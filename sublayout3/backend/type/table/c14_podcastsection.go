package table

import "time"

type C14PodcastSection struct {
	Id        *uint64    `gorm:"primaryKey"`
	Title     *string    `gorm:"type:VARCHAR(255); not null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
