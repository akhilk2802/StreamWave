package models

type User struct {
	UserID  uint   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name"`
	EmailID string `json:"emailId" gorm:"unique"`
}
