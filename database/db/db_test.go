package db

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	mysqlConf = &Conf{
		InstanceName: "mysql_rw",
		DriverName:   "mysql",
		DataSource:   "test:test123@tcp(10.4.61.88:30006)/d_test?charset=utf8",
	}
	pgConf = &Conf{
		InstanceName: "postgres_rw",
		DriverName:   "postgres",
		DataSource:   "postgres://test:test123@10.4.61.88:30032/d_test?sslmode=disable",
	}
	curInstance = "postgres_rw"
)

func TestNew(t *testing.T) {
	cs := []*Conf{mysqlConf, pgConf}
	db, err := New(cs)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestDB_GetInstance(t *testing.T) {
	cs := []*Conf{mysqlConf, pgConf}
	db, err := New(cs)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	d := db.GetInstance(curInstance)
	assert.NotNil(t, d)
}

type User struct {
	Name string
	Age  int
}

func TestNew2(t *testing.T) {
	cs := []*Conf{mysqlConf, pgConf}
	db, err := New(cs)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	d := db.GetInstance(curInstance)
	assert.NotNil(t, d)

	u := &User{
		Name: "test_002",
		Age:  14,
	}
	r, err := d.Insert("t_user").
		Prepared(true).
		Cols("name,age").
		Rows(u).
		OnConflict(goqu.DoUpdate("name", goqu.Record{"age": u.Age})).
		Executor().Exec()
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestNew3(t *testing.T) {
	cs := []*Conf{mysqlConf, pgConf}
	db, err := New(cs)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	d := db.GetInstance(curInstance)
	assert.NotNil(t, d)

	u := &User{
		Name: "test_002",
		Age:  14,
	}

	err = d.WithTx(func(tx *goqu.TxDatabase) error {
		_, err := tx.Insert("t_user").
			Prepared(true).
			Rows(u).
			OnConflict(goqu.DoUpdate("name", goqu.Record{"age": u.Age})).
			Executor().Exec()
		return err
	})
	assert.NoError(t, err)
}
