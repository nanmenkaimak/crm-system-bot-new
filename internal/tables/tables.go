package tables

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Room struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Bed struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	Name      string    `json:"name"`
	Cost      int       `json:"cost"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Customer struct {
	ID           int       `json:"id"`
	BedID        int       `json:"bed_id"`
	FullName     string    `json:"full_name"`
	Photo        string    `json:"photo"`
	Phone        string    `json:"phone"`
	Info         string    `json:"info"`
	Money        int       `json:"money"`
	IsHere       bool      `json:"is_here"`
	ArrivalDay   time.Time `json:"arrival_day"`
	DepartureDay time.Time `json:"departure_day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Expenses struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Money     int       `json:"money"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BedWithCustomer struct {
	Bed      Bed      `json:"bed"`
	Customer Customer `json:"customer"`
}
