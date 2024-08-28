package models

import "gorm.io/gorm"

type Stream struct {
	gorm.Model
	Title       string
	Description string
	URL         string
}
