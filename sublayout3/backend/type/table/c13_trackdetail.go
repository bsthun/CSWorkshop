package table

import "time"

type C13TrackDetail struct {
	Id             *uint64    `gorm:"primaryKey"`
	TrackTitle     *string    `gorm:"type:VARCHAR(255); not null"`
	ArtworkUrl     *string    `gorm:"type:TEXT; null"`
	ArtistName     *string    `gorm:"type:VARCHAR(255); not null"`
	ArtistImageUrl *string    `gorm:"type:TEXT; null"`
	ReleasedYear   *uint64    `gorm:"not null"`
	Length         *time.Time `gorm:"not null"`
	IsLiked        *bool      `gorm:"not null"`
	Lyrics         *string    `gorm:"type:TEXT; null"`
	CreatedAt      *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt      *time.Time `gorm:"not null"` // Embedded field
}
