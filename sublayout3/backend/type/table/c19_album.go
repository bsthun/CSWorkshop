package table

import "time"

type C19Album struct {
	Id             *uint64    `gorm:"primaryKey"`
	ImageUrl       *string    `gorm:"TEXT; null"`
	Type           *string    `gorm:"type:ENUM('album', 'single', 'compilation', 'ep'); not null"`
	Title          *string    `gorm:"type:VARCHAR(255); not null"`
	ArtistImageUrl *string    `gorm:"TEXT; null"`
	ArtistName     *string    `gorm:"type:VARCHAR(255); not null"`
	ReleaseYear    *uint64    `gorm:"not null"`
	CreatedAt      *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt      *time.Time `gorm:"not null"` // Embedded field
}
