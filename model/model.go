package model

import (
	"time"
)

type Entity interface {
	Faculty | Group
	EntityName() string
	GetId() int
}

type Faculty struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}

func (f Faculty) GetId() int {
	return f.Id
}

func (Faculty) EntityName() string {
	return "Faculty"
}

type Group struct {
	FacultyId int    `json:"faculty_id"`
	Name      string `json:"name"`
	Id        int    `json:"id"`
	Level     int    `json:"level"`
	Type      string `json:"type"`
	Kind      int    `json:"kind"`
	Year      int    `json:"year"`
}

func (g Group) GetId() int {
	return g.Id
}

func (Group) EntityName() string {
	return "Group"
}

type Lesson struct {
	FacultyId int    `json:"faculty_id"`
	GroupIds  []int  `json:"group_ids"`
	Subject   string `json:"subject"`
	Type      string `json:"type"`
	TeacherId string `json:"teacher_id"`
	PlaceId   int    `json:"place_id"`
}

type Teacher struct {
	Id   int    `json:"teacher_id"`
	Name string `json:"name"`
}

type Place struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
}

type Day struct {
	Date    time.Time `json:"date"`
	Lessons []Lesson  `json:"lessons"`
}

type Week struct {
	IsOdd bool  `json:"is_odd"`
	Days  []Day `json:"days"`
}
