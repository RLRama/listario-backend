package models

type UserProject struct {
	UserID    uint
	ProjectID uint
	Role      string `gorm:"default:'member'"` // owner, admin, member
}
