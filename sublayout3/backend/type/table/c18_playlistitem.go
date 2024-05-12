package table

import "time"

type C18PlaylistItem struct {
	Id          *uint64             `gorm:"primaryKey"`
	ArtworkUrl  *string             `gorm:"TEXT; null"`
	Title       *string             `gorm:"type:VARCHAR(255); not null"`
	Description *string             `gorm:"TEXT; null"`
	SectionId   *uint64             `gorm:"not null"`
	Section     *C18PlaylistSection `gorm:"foreignKey:SectionId"`
	CreatedAt   *time.Time          `gorm:"not null"` // Embedded field
	UpdatedAt   *time.Time          `gorm:"not null"` // Embedded field
}
