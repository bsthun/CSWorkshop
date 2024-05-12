package table

import "time"

type C11Detail struct {
	Id            *uint64    `gorm:"primaryKey"`
	Title         *string    `gorm:"type:VARCHAR(255); not null"`
	Venue         *string    `gorm:"type:VARCHAR(255); not null"`
	CoverImageUrl *string    `gorm:"type:TEXT; null"`
	CreatedAt     *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt     *time.Time `gorm:"not null"` // Embedded field
}
