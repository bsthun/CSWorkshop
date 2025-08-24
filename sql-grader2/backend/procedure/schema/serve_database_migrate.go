package schemaProcedure

import (
	"backend/util/orm"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/bsthun/gut"
)

func (r *Service) ServeDatabaseMigrate(ctx context.Context, mysqlDsn string, databaseName string, schemaContent string) *gut.ErrorInstance {
	// * create database
	tx := r.gorm.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", databaseName))
	if tx.Error != nil {
		return gut.Err(false, "failed to create database", tx.Error)
	}

	// * remove create database statements
	content := regexp.MustCompile(`(?im)^\s*CREATE\s+DATABASE\s+.*$`).ReplaceAllString(schemaContent, "")

	// * replace existing use statements or add new one
	if regexp.MustCompile(`(?im)^\s*USE\s+`).MatchString(content) {
		content = regexp.MustCompile(`(?im)^\s*USE\s+.*$`).ReplaceAllString(content, fmt.Sprintf("USE `%s`;", databaseName))
	} else {
		content = fmt.Sprintf("USE `%s`;\n\n%s", databaseName, content)
	}

	// * parse mysql dsn
	re := regexp.MustCompile(`^([^:]+):([^@]+)@tcp\(([^:]+):(\d+)\)/`)
	matches := re.FindStringSubmatch(mysqlDsn)
	if len(matches) != 5 {
		return gut.Err(false, "invalid mysql dsn format", nil)
	}

	// * execute schema file using mysql command
	cmd := exec.Command("mysql", "--protocol=TCP", "-h", matches[3], "-P", matches[4], "-u", matches[1], "-p"+matches[2], databaseName)

	cmd.Stdin = bytes.NewBuffer([]byte(content))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return gut.Err(false, "failed to execute schema", err)
	}

	// * test connection to migrated database
	conn, er := orm.Connect(mysqlDsn, databaseName)
	if er != nil {
		return er
	}
	defer conn.Close()

	return nil
}