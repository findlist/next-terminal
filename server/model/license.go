package model

type License struct {
	ID            string `gorm:"primary_key,type:varchar(36)" json:"id"`
	Type          string `gorm:"type:varchar(20)" json:"type"`
	MachineID     string `gorm:"type:varchar(50)" json:"machineId"`
	Name          string `gorm:"type:varchar(100)" json:"name"`
	Phone         string `gorm:"type:varchar(20)" json:"phone"`
	Address       string `gorm:"type:varchar(500)" json:"address"`
	Asset         int    `gorm:"column:asset" json:"asset"`
	SurplusAssets int64  `json:"surplusAssets"`
	Concurrent    int    `gorm:"column:concurrent" json:"concurrent"`
	User          int    `gorm:"column:user" json:"user"`
	Expired       string `gorm:"column:expired" json:"expired"`
}

func (l *License) TableName() string {
	return "licenses"
}
