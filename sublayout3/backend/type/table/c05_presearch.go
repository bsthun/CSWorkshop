package table

import "time"

type C05Presearch struct {
	Id        *uint64    `gorm:"primaryKey"`
	Name      *string    `gorm:"VARCHAR(255); not null"`
	CoverUrl  *string    `gorm:"TEXT; null"`
	Color     *string    `gorm:"VARCHAR(255); null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
