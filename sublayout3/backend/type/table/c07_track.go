package table

import "time"

type C07Track struct {
	Id          *uint64    `gorm:"primaryKey"`
	Title       *string    `gorm:"type:VARCHAR(255); not null"`
	Artwork     *string    `gorm:"type:TEXT; null"`
	ListenCount *uint64    `gorm:"not null"`
	Duration    *time.Time `gorm:"not null"`
	ArtworkUrl  *string    `gorm:"type:TEXT; null"`
	CreatedAt   *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt   *time.Time `gorm:"not null"` // Embedded field
}
