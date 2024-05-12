package table

import "time"

type C08Track struct {
	Id        *uint64        `gorm:"primaryKey"`
	AlbumId   *uint64        `gorm:"not null"`
	Album     *C08Album      `gorm:"foreignKey:AlbumId"`
	Title     *string        `gorm:"type:VARCHAR(255); not null"`
	Artist    *string        `gorm:"type:VARCHAR(255); not null"`
	PlayCount *uint64        `gorm:"not null"`
	Duration  *time.Duration `gorm:"not null"`
	CreatedAt *time.Time     `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time     `gorm:"not null"` // Embedded field
}
