package db

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
)

type Conf struct {
	InstanceName string
	DriverName   string
	DataSource   string
	IsSlave      bool
}

type DB struct {
	opts          Options
	dataSet       map[string]*goqu.Database
	masterDataSet map[string]*goqu.Database
	slaveDataSet  map[string]*goqu.Database
}

func New(cs []*Conf, opts ...Option) (*DB, error) {
	options := newOptions(opts...)
	db := &DB{
		opts:          options,
		dataSet:       make(map[string]*goqu.Database),
		masterDataSet: make(map[string]*goqu.Database),
		slaveDataSet:  make(map[string]*goqu.Database),
	}
	for _, c := range cs {
		if _, ok := db.dataSet[c.InstanceName]; !ok {
			sess, err := db.open(c.DriverName, c.DataSource)
			if err != nil {
				return nil, err
			}
			db.dataSet[c.InstanceName] = sess

			if c.IsSlave {
				db.slaveDataSet[c.InstanceName] = db.dataSet[c.InstanceName]
			} else {
				db.masterDataSet[c.InstanceName] = db.dataSet[c.InstanceName]
			}
		}
	}

	if len(db.masterDataSet) == 0 {
		return nil, fmt.Errorf("Master instance is empty ")
	}

	return db, nil
}

func (d *DB) open(driverName, dataSourceName string) (*goqu.Database, error) {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(d.opts.maxOpen)
	conn.SetMaxIdleConns(d.opts.maxIdle)
	conn.SetConnMaxLifetime(d.opts.maxLifetime)
	db := goqu.New(driverName, conn)
	l := zerolog.New(logger.DefaultLogger.Options().Out)
	db.Logger(&l)
	return db, nil
}

func (d *DB) GetInstance(instanceName string, isSlaver ...bool) *goqu.Database {
	slaver := len(isSlaver) > 0 && isSlaver[0] == true
	switch slaver {
	case true:
		return d.SlaveSession(instanceName)
	case false:
		return d.MasterSession(instanceName)
	}

	return d.MasterSession(instanceName)
}

// MasterSession 主库
func (d *DB) MasterSession(instanceName string) *goqu.Database {
	if len(instanceName) == 0 {
		for _, db := range d.masterDataSet {
			return db
		}
	}

	if i, ok := d.masterDataSet[instanceName]; ok {
		return i
	}
	return nil
}

// SlaveSession 从库
func (d *DB) SlaveSession(instanceName string) *goqu.Database {
	if len(d.slaveDataSet) == 0 {
		return d.MasterSession(instanceName)
	}

	if len(instanceName) == 0 {
		for _, db := range d.slaveDataSet {
			return db
		}
	}

	if db, ok := d.slaveDataSet[instanceName]; ok {
		return db
	}

	return nil
}
