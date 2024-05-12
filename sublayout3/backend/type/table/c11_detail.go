package table

import "time"

type C12Detail struct {
	Id        *uint64    `gorm:"primaryKey"`
	Date      *time.Time `gorm:"not null"`
	Title     *string    `gorm:"type:VARCHAR(255); not null"`
	Venue     *string    `gorm:"type:VARCHAR(255); not null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
