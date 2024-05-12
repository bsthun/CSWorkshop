package table

import "time"

type C12Playlist struct {
	Id        *uint64    `gorm:"primaryKey"`
	ImageUrl  *string    `gorm:"TEXT; null"`
	Type      *string    `gorm:"type:ENUM('private', 'public'); not null"`
	Title     *string    `gorm:"type:VARCHAR(255); not null"`
	Creator   *string    `gorm:"type:VARCHAR(255); not null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
