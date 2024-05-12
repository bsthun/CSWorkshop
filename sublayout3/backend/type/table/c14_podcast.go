package table

import "time"

type C14PodcastGenre struct {
	Id            *uint64            `gorm:"primaryKey"`
	Title         *string            `gorm:"type:VARCHAR(255); not null"`
	SectionId     *uint64            `gorm:"not null"`
	Section       *C14PodcastSection `gorm:"foreignKey:SectionId; references:Id"`
	CoverImageUrl *string            `gorm:"TEXT; null"`
	CreatedAt     *time.Time         `gorm:"not null"` // Embedded field
	UpdatedAt     *time.Time         `gorm:"not null"` // Embedded field
}
