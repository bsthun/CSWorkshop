package table

import "time"

type SystemGroup struct {
	Id        *uint64    `gorm:"primaryKey"`
	No        *uint64    `gorm:"not null"`
	Address   *string    `gorm:"type:VARCHAR(255); not null"`
	Port      *string    `gorm:"type:VARCHAR(255); not null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
