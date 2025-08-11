package adminEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"backend/util/orm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleCollectionSchemaUpload(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get user from database
	user, err := r.database.P().UserGetById(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get user", err)
	}

	// * check if user is admin
	if user.IsAdmin == nil || !*user.IsAdmin {
		return gut.Err(false, "access denied", nil)
	}

	// * get collectionId from form
	collectionIdStr := c.FormValue("collectionId")
	if collectionIdStr == "" {
		return gut.Err(false, "collectionId is required", nil)
	}

	collectionId, er := gut.Decode(collectionIdStr)
	if er != nil {
		return er
	}

	// * get file from multipart form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return gut.Err(false, "failed to get file", err)
	}

	// * create directory
	dirPath := fmt.Sprintf(".local/collection/%s", gut.EncodeId(collectionId))
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return gut.Err(false, "failed to create directory", err)
	}

	// * save file
	filePath := filepath.Join(dirPath, "schema.sql")
	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return gut.Err(false, "failed to save file", err)
	}

	// * create temporary database
	tempDbName := fmt.Sprintf("tmp_%s", *gut.Random(gut.RandomSet.MixedAlphaNum, 8))

	// * create database
	tx := r.gorm.Exec(fmt.Sprintf("CREATE DATABASE `%s`", tempDbName))
	if tx.Error != nil {
		return gut.Err(false, "failed to create temporary database", tx.Error)
	}

	// * cleanup function
	cleanup := func() {
		r.gorm.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", tempDbName))
	}
	defer cleanup()

	// * connect to temporary database
	tempConn, er := orm.Connect(*r.config.MysqlDsn, tempDbName)
	if er != nil {
		return er
	}
	defer tempConn.Close()

	// * read schema file
	schemaContent, err := os.ReadFile(filePath)
	if err != nil {
		return gut.Err(false, "failed to read schema file", err)
	}

	// * execute schema in temporary database
	if _, err := tempConn.Exec(string(schemaContent)); err != nil {
		return gut.Err(false, "failed to execute schema", err)
	}

	// * get table list and row counts
	rows, err := tempConn.Query("SHOW TABLES")
	if err != nil {
		return gut.Err(false, "failed to get table list", err)
	}
	defer rows.Close()

	var structure []*payload.CollectionTableStructure
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return gut.Err(false, "failed to scan table name", err)
		}

		// * get row count
		var rowCount uint64
		if err := tempConn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s`", tableName)).Scan(&rowCount); err != nil {
			return gut.Err(false, "failed to get row count", err)
		}

		structure = append(structure, &payload.CollectionTableStructure{
			TableName: &tableName,
			RowCount:  &rowCount,
		})
	}

	// * marshal metadata
	metadata := payload.CollectionSchemaMetadata{
		Structure: structure,
	}

	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return gut.Err(false, "failed to marshal metadata", err)
	}

	// * update collection metadata
	_, err = r.database.P().CollectionUpdateMetadata(c.Context(), &psql.CollectionUpdateMetadataParams{
		Id:       &collectionId,
		Metadata: metadataJson,
	})
	if err != nil {
		return gut.Err(false, "failed to update collection metadata", err)
	}

	// * response
	return c.JSON(response.Success(c, "database structure analyzed successfully"))
}
