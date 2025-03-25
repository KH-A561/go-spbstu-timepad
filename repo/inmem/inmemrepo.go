package inmem

import (
	"context"
	"errors"
	"slices"
	"universityTimepad/model"
)

type MemoryRepository[E model.Entity] struct {
	Elements *[]E
}

func NewMemoryRepository[E model.Entity](dataProvider func() *[]E) *MemoryRepository[E] {
	result := new(MemoryRepository[E])
	result.Elements = dataProvider()
	return result
}

func (r *MemoryRepository[E]) GetByID(ctx context.Context, id int) (*E, error) {
	index := slices.IndexFunc(*r.Elements, func(e E) bool { return e.GetId() == id })
	elements := *r.Elements
	if index != -1 {
		return &(elements[index]), nil
	} else {
		return nil, errors.New("element not found")
	}
}

func (r *MemoryRepository[E]) GetAll(ctx context.Context) (*[]E, error) {
	return r.Elements, nil
}
