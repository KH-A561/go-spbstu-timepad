package inmem

import (
	"log"
	"universityTimepad/files"
	"universityTimepad/model"
)

type DataInitFunc[E model.Entity] struct {
	initFunc func() *[]E
}

func (d DataInitFunc[E]) GetInitFunc() func() *[]E {
	return d.initFunc
}

type Config struct {
	FacInitFunc   DataInitFunc[model.Faculty]
	GroupInitFunc DataInitFunc[model.Group]
}

func NewDefault() *Config {
	result := Config{}
	faculties := make([]model.Faculty, 0)
	groups := make([]model.Group, 0)
	facultiesWithGroups, err := files.ReadFacultiesWithGroups()

	if err != nil {
		log.Fatalf("Error reading faculties: %v", err)
	}

	for k, v := range *facultiesWithGroups {
		faculties = append(faculties, *k)
		groups = append(groups, *v...)
	}

	FacInitFunc := DataInitFunc[model.Faculty]{}
	FacInitFunc.initFunc = func() *[]model.Faculty {
		return &faculties
	}
	result.FacInitFunc = FacInitFunc

	GroupInitFunc := DataInitFunc[model.Group]{}
	GroupInitFunc.initFunc = func() *[]model.Group {
		return &groups
	}
	result.GroupInitFunc = GroupInitFunc

	return &result
}
