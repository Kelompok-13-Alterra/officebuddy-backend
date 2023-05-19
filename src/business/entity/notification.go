package entity

import (
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	office_id	uint
	user_id		uint
	step		string
	status		bool
}



