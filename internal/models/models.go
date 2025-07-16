package models

type Employee struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Sex    string  `json:"sex"`
	Age    int     `json:"age"`
	Salary float64 `json:"salary"`
}
