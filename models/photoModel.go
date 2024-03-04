package models

type Photo struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Title    string `gorm:"not null;varchar(255)" json:"title"`
	Caption  string `gorm:"not null;varchar(255)" json:"caption"`
	PhotoUrl string `gorm:"not null;varchar(255)" json:"photo_url"`
	UserID   int    `gorm:"not null;varchar(255)" json:"user_id"`
}
