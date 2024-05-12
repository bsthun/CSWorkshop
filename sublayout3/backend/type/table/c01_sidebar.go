package table

import "time"

type C01Sidebar struct {
	Id         *uint64    `gorm:"primaryKey"`
	Type       *string    `gorm:"type:ENUM('likedSong', 'likedEpisode', 'album', 'playlist', 'podcast', 'artist'); not null"`
	Title      *string    `gorm:"type:VARCHAR(255); not null"`
	Creator    *string    `gorm:"type:VARCHAR(255); null"`
	ArtworkUrl *string    `gorm:"type:TEXT; null"`
	IsPlaying  *bool      `gorm:"not null"`
	AddedAt    *time.Time `gorm:"not null"`
	CreatedAt  *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt  *time.Time `gorm:"not null"` // Embedded field
}
