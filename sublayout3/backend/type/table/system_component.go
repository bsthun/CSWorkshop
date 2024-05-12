package table

import "time"

type SystemComponent struct {
	Id            *uint64      `gorm:"primaryKey"`
	Name          *string      `gorm:"type:VARCHAR(255); not null"`
	SideAGroupId  *uint64      `gorm:"not null"`
	SideAGroup    *SystemGroup `gorm:"foreignKey:SideAGroupId"`
	SideBGroupId  *uint64      `gorm:"not null"`
	SideBGroup    *SystemGroup `gorm:"foreignKey:SideBGroupId"`
	QueryKey      *string      `gorm:"type:VARCHAR(255); null"`
	QueryValStart *int64       `gorm:"null"`
	QueryValEnd   *int64       `gorm:"null"`
	CreatedAt     *time.Time   `gorm:"not null"` // Embedded field
	UpdatedAt     *time.Time   `gorm:"not null"` // Embedded field
}
