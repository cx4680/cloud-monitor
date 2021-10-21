package dao

import "github.com/jinzhu/gorm"

type InstanceDao struct {
	db *gorm.DB
}

func NewInstanceDao(db *gorm.DB) *InstanceDao {
	return &InstanceDao{db: db}
}
