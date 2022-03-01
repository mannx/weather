package models

import "gorm.io/gorm"

// ConfigEntry contains a single city entry to monitor
type ConfigEntry struct {
	gorm.Model

	CityID int    // id of the given city
	Name   string // name of the city
	Active bool   // are we activily gathering data from this entry
}
