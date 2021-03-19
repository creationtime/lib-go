package db

import (
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/micro/go-micro/v2/logger"
	"github.com/rs/zerolog"
)

type Conf struct {
	InstanceName string
	DriverName   string
	DataSource   string
}

type DB struct {
	opts    Options
	dataSet map[string]*goqu.Database
}

func New(cs []*Conf, opts ...Option) (*DB, error) {
	options := newOptions(opts...)
	db := &DB{
		opts:    options,
		dataSet: make(map[string]*goqu.Database),
	}
	for _, c := range cs {
		if _, ok := db.dataSet[c.InstanceName]; !ok {
			sess, err := db.open(c.DriverName, c.DataSource)
			if err != nil {
				return nil, err
			}
			db.dataSet[c.InstanceName] = sess
		}
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

func (d *DB) GetInstance(instanceName string) *goqu.Database {
	if sess, ok := d.dataSet[instanceName]; ok {
		return sess
	}
	for _, v := range d.dataSet {
		return v
	}
	return nil
}
