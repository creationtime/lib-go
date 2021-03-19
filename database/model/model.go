package model

import (
	"time"
)

type Model struct {
	TableName  string    `jsonpb:"-" db:"-" goqu:"skipinsert,skipupdate"`
	Id         int64     `jsonpb:"id" db:"id" goqu:"skipinsert,skipupdate"`
	CreateTime time.Time `jsonpb:"-" db:"createTime" goqu:"skipinsert,skipupdate"`
	UpdateTime time.Time `jsonpb:"-" db:"updateTime" goqu:"skipinsert,skipupdate"`
}
