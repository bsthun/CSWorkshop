package table

import "time"

type C06Search struct {
	Id           *uint64    `gorm:"primaryKey"`
	TrackTitle   *string    `gorm:"type:VARCHAR(255); not null"`
	TrackArtwork *string    `gorm:"type:TEXT; null"`
	ArtistName   *string    `gorm:"type:VARCHAR(255); not null"`
	Duration     *string    `gorm:"not null"`
	IsLiked      *bool      `gorm:"not null"`
	CreatedAt    *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt    *time.Time `gorm:"not null"` // Embedded field
}
