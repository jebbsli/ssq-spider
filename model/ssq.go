package model

type SysConfig struct {
	TypeId string `gorm:"column:type_id"`
	Name   string `gorm:"column:name"`
	Value  string `gorm:"column:value"`
}

type SSQOpenNumber struct {
	OpenNo   string `gorm:"column:open_no"`
	RedNum   string `gorm:"column:red_num"`
	BlueNum  string `gorm:"column:blue_num"`
	BallSort string `gorm:"column:ball_sort"`
}
