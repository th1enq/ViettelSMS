package entity

import "github.com/lib/pq"

type AuthUser struct {
	ID       uint           `gorm:"primaryKey"`
	Username string         `gorm:"unique;not null"`
	Password string         `gorm:"not null"`
	Blocked  bool           `gorm:"not null;default:false"`
	Scopes   pq.StringArray `gorm:"type:text[];not null;default:'{}'"`
}
