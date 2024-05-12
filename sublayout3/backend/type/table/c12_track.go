package table

import "time"

type C12Track struct {
	Id         *uint64    `gorm:"primaryKey"`
	ArtworkUrl *string    `gorm:"TEXT; null"`
	Title      *string    `gorm:"type:VARCHAR(255); not null"`
	ArtistName *string    `gorm:"type:VARCHAR(255); not null"`
	AlbumName  *string    `gorm:"type:VARCHAR(255); not null"`
	DateAdded  *time.Time `gorm:"not null"`
	Length     *time.Time `gorm:"not null"`
	CreatedAt  *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt  *time.Time `gorm:"not null"` // Embedded field
}
