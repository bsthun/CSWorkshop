package table

import "time"

type C07Artist struct {
	Id                   *uint64    `gorm:"primaryKey"`
	Name                 *string    `gorm:"type:VARCHAR(255); not null"`
	IsVerified           *bool      `gorm:"not null"`
	CoverImageUrl        *string    `gorm:"type:TEXT; null"`
	Bio                  *string    `gorm:"type:TEXT; null"`
	IsFollowed           *bool      `gorm:"not null"`
	MonthlyListenerCount *uint64    `gorm:"not null"`
	CreatedAt            *time.Time `gorm:"not null"` // Embedded field
	UpdatedAt            *time.Time `gorm:"not null"` // Embedded field
}
