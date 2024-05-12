package table

import "time"

type C17PodcastEpisode struct {
	Id          *uint64         `gorm:"primaryKey"`
	ImageUrl    *string         `gorm:"TEXT; null"`
	Title       *string         `gorm:"type:VARCHAR(255); not null"`
	ShowId      *uint64         `gorm:"not null"`
	Show        *C16PodcastShow `gorm:"foreignKey:ShowId; references:Id"`
	Description *string         `gorm:"TEXT; null"`
	Date        *time.Time      `gorm:"not null"`
	Length      *uint64         `gorm:"not null"`
	CreatedAt   *time.Time      `gorm:"not null"` // Embedded field
	UpdatedAt   *time.Time      `gorm:"not null"` // Embedded field
}
