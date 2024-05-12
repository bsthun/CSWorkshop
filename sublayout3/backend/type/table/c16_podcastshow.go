package table

import "time"

type C16PodcastShow struct {
	Id        *uint64    `gorm:"primaryKey"`
	ImageUrl  *string    `gorm:"TEXT; null"`
	Title     *string    `gorm:"type:VARCHAR(255); not null"`
	Author    *string    `gorm:"type:VARCHAR(255); not null"`
	About     *string    `gorm:"TEXT; null"`
	Rating    *float64   `gorm:"not null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
