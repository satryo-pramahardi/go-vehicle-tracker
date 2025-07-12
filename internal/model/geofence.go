package model

type Geofence struct {
	ID        int64   `gorm:"primaryKey"`
	Name      string  `gorm:"not null"`
	CenterLat float64 `gorm:"not null"`
	CenterLng float64 `gorm:"not null"`
	Radius    float64 `gorm:"not null"` // in meters
	Active    bool    `gorm:"default:true"`
}

func (Geofence) TableName() string {
	return "geofences"
}
