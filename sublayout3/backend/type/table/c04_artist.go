package table

import "time"

type C04Artist struct {
	Id                   *uint64    `gorm:"primaryKey"`
	CoverUrl             *string    `gorm:"type:TEXT; null"`
	Name                 *string    `gorm:"type:VARCHAR(255); not null"`
	Bio                  *string    `gorm:"type:TEXT; null"`
	FollowerCount        *uint64    `gorm:"not null"`
	MonthlyListenerCount *uint64    `gorm:"not null"`
	IsVerified           *bool      `gorm:"not null"`
	CreatedAt            *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt            *time.Time `gorm:"not null"` // Embedded field
}
