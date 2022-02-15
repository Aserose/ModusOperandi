package boltDB

import (
	"github.com/Aserose/ModusOperandi/pkg/config"
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"github.com/boltdb/bolt"
	"time"
)

func ConnectBoltDB(log logger.Logger, cfg config.CfgBolt) *bolt.DB {
	db, err := bolt.Open(cfg.RootName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Panicf("%s: %s", log.CallInfoStr(), err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(cfg.BucketInstructionName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Panicf("%s: %s", log.CallInfoStr(), err.Error())
	}

	return db
}
