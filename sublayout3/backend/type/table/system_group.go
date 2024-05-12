package table

import "time"

type SystemGroup struct {
	Id        *uint64    `gorm:"primaryKey"`
	No        *uint64    `gorm:"not null"`
	Address   *string    `gorm:"type:VARCHAR(255); not null"`
	Side      *string    `gorm:"type:ENUM('a', 'b'); not null"`
	CreatedAt *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time `gorm:"not null"` // Embedded field
}
