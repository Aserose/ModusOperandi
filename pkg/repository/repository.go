package repository

import (
	"github.com/Aserose/ModusOperandi/pkg/repository/boltDB/data"
)

type DB struct {
	BoltDB *data.BoltData
}

func NewDB(BoltDB *data.BoltData) *DB {
	return &DB{
		BoltDB: BoltDB,
	}
}
