package tuple

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type CollectionTableStructure struct {
	TableName *string `json:"tableName"`
	RowCount  *uint64 `json:"rowCount"`
}

type CollectionSchemaMetadata struct {
	SchemaFilename *string                     `json:"schemaFilename"`
	Structure      []*CollectionTableStructure `json:"structure"`
}

func (r *CollectionSchemaMetadata) Scan(value any) error {
	if value == nil {
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, r)
}

func (r *CollectionSchemaMetadata) Value() (driver.Value, error) {
	return json.Marshal(r)
}
