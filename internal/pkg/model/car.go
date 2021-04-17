package model

import "time"

type Car struct {
	ID  int64
	Lat float64
	Lng float64
}

type NearestCar struct {
	ID             int64
	ArrivalMinutes int64
}

type CachedCar struct {
	ID             int64     `json:"id"`
	ArrivalMinutes int64     `json:"eta_minutes"`
	TimeCreated    time.Time `json:"time_created"`
}
