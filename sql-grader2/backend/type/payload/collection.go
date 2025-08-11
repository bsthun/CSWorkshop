package payload

import (
	"backend/type/common"
	"encoding/json"
	"time"
)

type CollectionListRequest struct {
	common.Paginate
	Name *string `json:"name"`
}

type Collection struct {
	Id        *uint64         `json:"id"`
	Name      *string         `json:"name"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt *time.Time      `json:"createdAt"`
	UpdatedAt *time.Time      `json:"updatedAt"`
}

type CollectionListResponse struct {
	Count       *uint64       `json:"count"`
	Collections []*Collection `json:"collections"`
}

type SemesterListRequest struct {
	common.Paginate
	Name *string `json:"name"`
	Sort *string `json:"sort" validate:"omitempty,oneof=name createdAt"`
}

type SemesterClass struct {
	Id           *uint64    `json:"id"`
	SemesterId   *uint64    `json:"semesterId"`
	Name         *string    `json:"name"`
	RegisterCode *string    `json:"registerCode"`
	JoineeCount  *uint64    `json:"joineeCount"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type Semester struct {
	Id        *uint64          `json:"id"`
	Name      *string          `json:"name"`
	Classes   []*SemesterClass `json:"classes"`
	CreatedAt *time.Time       `json:"createdAt"`
	UpdatedAt *time.Time       `json:"updatedAt"`
}

type SemesterListResponse struct {
	Count     *uint64     `json:"count"`
	Semesters []*Semester `json:"semesters"`
}

type SemesterCreateRequest struct {
	Name *string `json:"name" validate:"required"`
}

type SemesterEditRequest struct {
	Id   *uint64 `json:"id" validate:"required"`
	Name *string `json:"name" validate:"required"`
}

type ClassCreateRequest struct {
	SemesterId *uint64 `json:"semesterId" validate:"required"`
	Code       *string `json:"code" validate:"required"`
	Name       *string `json:"name" validate:"required"`
}

type CollectionCreateRequest struct {
	Name *string `json:"name" validate:"required"`
}

type CollectionSchemaUploadRequest struct {
	CollectionId string `json:"collectionId" validate:"required"`
}

type CollectionTableStructure struct {
	TableName *string `json:"tableName"`
	RowCount  *uint64 `json:"rowCount"`
}

type CollectionSchemaMetadata struct {
	SchemaFilename *string                     `json:"schemaFilename"`
	Structure      []*CollectionTableStructure `json:"structure"`
}

type CollectionQuestionCreateRequest struct {
	CollectionId *uint64 `json:"collectionId" validate:"required"`
	Description  *string `json:"description"`
}

type CollectionQuestionEditRequest struct {
	Id          *uint64 `json:"id" validate:"required"`
	Title       *string `json:"title" validate:"required"`
	Description *string `json:"description"`
	CheckQuery  *string `json:"checkQuery" validate:"required"`
	CheckPrompt *string `json:"checkPrompt" validate:"required"`
}

type CollectionQuestionListItem struct {
	Id       *uint64 `json:"id"`
	OrderNum *int32  `json:"orderNum"`
	Title    *string `json:"title"`
}

type CollectionQuestionListResponse struct {
	Questions []*CollectionQuestionListItem `json:"questions"`
}

type CollectionQuestionDetail struct {
	Id           *uint64    `json:"id"`
	CollectionId *uint64    `json:"collectionId"`
	OrderNum     *int32     `json:"orderNum"`
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	CheckQuery   *string    `json:"checkQuery"`
	CheckPrompt  *string    `json:"checkPrompt"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type CollectionQuestionDeleteRequest struct {
	Id *uint64 `json:"id" validate:"required"`
}
