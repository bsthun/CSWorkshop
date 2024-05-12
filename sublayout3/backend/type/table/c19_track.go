package table

import "time"

type C19Track struct {
	Id        *uint64    `gorm:"primaryKey"`
	Title     *string    `gorm:"type:VARCHAR(255); not null"`
	Artist    *string    `gorm:"type:VARCHAR(255); not null"`
	PlayCount *uint64    `gorm:"not null"`
	Length    *time.Time `gorm:"not null"`
	AlbumId   *uint64    `gorm:"not null"`
	Album     *C19Album  `gorm:"foreignKey:AlbumId"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
