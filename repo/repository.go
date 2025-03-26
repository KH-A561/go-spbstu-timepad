package repo

import (
	"context"
	"errors"
	"log"
	"universityTimepad/model"
	"universityTimepad/repo/inmem"
	"universityTimepad/repo/postgres"
)

type StorageType string

const (
	StorageTypePostgres StorageType = "pg"
	StorageTypeInmemory StorageType = "mem"
)

type Config struct {
	StorageType StorageType
	Postgres    *postgres.Config
	Inmemory    *inmem.Config
}

type Repository[E model.Entity] interface {
	GetByID(ctx context.Context, id int) (*E, error)
	GetByName(ctx context.Context, name string) (*E, error)
	GetAll(ctx context.Context) (*[]E, error)
}

func New[E model.Entity](cfg *Config, dataProvider func() *[]E) (Repository[E], error) {
	switch cfg.StorageType {
	case StorageTypePostgres:
		log.Println("Postgres storage is not implemented yet")
		return nil, errors.New("postgres storage is not implemented yet")
	case StorageTypeInmemory:
		return NewMemoryRepository[E](dataProvider), nil
	default:
		return nil, errors.New("unknown storage type")
	}
}
