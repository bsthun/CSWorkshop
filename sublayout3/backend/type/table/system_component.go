package table

import "time"

type SystemComponent struct {
	Id        *uint64      `gorm:"primaryKey"`
	Name      *string      `gorm:"type:VARCHAR(255); not null"`
	GroupId   *uint64      `gorm:"not null"`
	Group     *SystemGroup `gorm:"foreignKey:GroupId"`
	CreatedAt *time.Time   `gorm:"not null"` // Embedded field
	UpdatedAt *time.Time   `gorm:"not null"` // Embedded field
}
