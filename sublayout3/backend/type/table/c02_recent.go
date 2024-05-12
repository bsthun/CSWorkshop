package table

import "time"

type C02Recent struct {
	Id        *uint64    `gorm:"primaryKey"`
	Name      *string    `gorm:"type:VARCHAR(255); not null"`
	ImageUrl  *string    `gorm:"type:TEXT; null"`
	IsPlaying *bool      `gorm:"not null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
