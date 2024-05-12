package table

import "time"

type C03Track struct {
	Id                  *uint64    `gorm:"primaryKey"`
	Title               *string    `gorm:"type:VARCHAR(255); not null"`
	ArtworkUrl          *string    `gorm:"type:TEXT; null"`
	IsLiked             *bool      `gorm:"not null"`
	ArtistName          *string    `gorm:"type:VARCHAR(255); not null"`
	ArtistImageUrl      *string    `gorm:"type:TEXT; null"`
	ArtistListenerCount *uint64    `gorm:"not null"`
	ArtistBio           *string    `gorm:"type:TEXT; null"`
	IsArtistVerified    *bool      `gorm:"not null"`
	ISArtistFollowed    *bool      `gorm:"not null"`
	CreatedAt           *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt           *time.Time `gorm:"not null"` // Embedded field
}
