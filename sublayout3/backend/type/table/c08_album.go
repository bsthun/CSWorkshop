package table

import "time"

type C08Album struct {
	Id          *uint64    `gorm:"primaryKey"`
	ArtworkUrl  *string    `gorm:"type:TEXT; null"`
	Title       *string    `gorm:"type:VARCHAR(255); not null"`
	Type        *string    `gorm:"type:ENUM('album', 'single', 'compilation', 'ep'); not null"`
	ReleaseYear *uint64    `gorm:"not null"`
	CreatedAt   *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt   *time.Time `gorm:"not null"` // Embedded field
}
