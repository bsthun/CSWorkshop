package adminEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/response"
	"backend/type/tuple"
	"backend/util/orm"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

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

	collectionId, err := gut.Decode(collectionIdStr)
	if err != nil {
		return gut.Err(false, "invalid collection id", err)
	}

	// * get file from multipart form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return gut.Err(false, "failed to get file", err)
	}

	// * create directory
	dirPath := fmt.Sprintf(".local/collection/%s", gut.Base62(collectionId))
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return gut.Err(false, "failed to create directory", err)
	}

	// * save file
	filePath := filepath.Join(dirPath, "schema.sql")
	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return gut.Err(false, "failed to save schema file", err)
	}

	// * read original content
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return gut.Err(false, "failed to read schema file", err)
	}

	// * create temporary database
	tempDatabaseName := fmt.Sprintf("tmp_%s", *gut.Random(gut.RandomSet.MixedAlphaNum, 8))

	// * remove create database statements
	content := regexp.MustCompile(`(?im)^\s*CREATE\s+DATABASE\s+.*$`).ReplaceAllString(string(raw), "")

	// * replace existing use statements or add new one
	if regexp.MustCompile(`(?im)^\s*USE\s+`).MatchString(content) {
		content = regexp.MustCompile(`(?im)^\s*USE\s+.*$`).ReplaceAllString(content, fmt.Sprintf("USE `%s`;", tempDatabaseName))
	} else {
		content = fmt.Sprintf("USE `%s`;\n\n%s", tempDatabaseName, content)
	}

	// * create database
	tx := r.gorm.Exec(fmt.Sprintf("CREATE DATABASE `%s`", tempDatabaseName))
	if tx.Error != nil {
		return gut.Err(false, "failed to create temporary database", tx.Error)
	}

	// * cleanup function
	cleanup := func() {
		r.gorm.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", tempDatabaseName))
	}
	defer cleanup()

	re := regexp.MustCompile(`^([^:]+):([^@]+)@tcp\(([^:]+):(\d+)\)/`)
	matches := re.FindStringSubmatch(*r.config.MysqlDsn)
	if len(matches) != 5 {
		return gut.Err(false, "invalid mysql dsn format", nil)
	}

	// * execute schema file using mysql command
	cmd := exec.Command("mysql", "--protocol=TCP", "--ssl-verify-server-cert=false", "-h", matches[3], "-P", matches[4], "-u", matches[1], "-p"+matches[2], tempDatabaseName)

	cmd.Stdin = bytes.NewBuffer([]byte(content))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return gut.Err(false, "failed to execute schema", err)
	}

	// * connect to temporary database to get table information
	tempConn, er := orm.Connect(*r.config.MysqlDsn, tempDatabaseName)
	if er != nil {
		return er
	}
	defer tempConn.Close()

	// * get table list and row counts
	rows, err := tempConn.Query("SHOW TABLES")
	if err != nil {
		return gut.Err(false, "failed to get table list", err)
	}
	defer rows.Close()

	var structure []*tuple.CollectionTableStructure
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

		structure = append(structure, &tuple.CollectionTableStructure{
			TableName: &tableName,
			RowCount:  &rowCount,
		})
	}

	// * update collection metadata
	_, err = r.database.P().CollectionUpdateMetadata(c.Context(), &psql.CollectionUpdateMetadataParams{
		Id: &collectionId,
		Metadata: &tuple.CollectionSchemaMetadata{
			SchemaFilename: &fileHeader.Filename,
			Structure:      structure,
		},
	})
	if err != nil {
		return gut.Err(false, "failed to update collection metadata", err)
	}

	// * response
	return c.JSON(response.Success(c, "database structure analyzed successfully"))
}
