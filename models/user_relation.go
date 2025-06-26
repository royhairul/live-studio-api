package models

import (
	"gorm.io/gorm"
)

type UserRelation struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	ID         uint   `gorm:"primaryKey"`
	ParentID   uint   // user yang membuat
	ChildID    uint   // user yang dibuat
	Relation   string // optional, misalnya: "created", "managed", dst
}
