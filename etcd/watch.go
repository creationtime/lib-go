package etcd

import (
	"github.com/creationtime/lib-go/logger"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/source/etcd"
)

type ConfWatch struct {
	conf config.Config
}

func New(addresses []string, prefix string) (*ConfWatch, error) {
	source := etcd.NewSource(
		etcd.WithAddress(addresses...),
		etcd.WithPrefix(prefix),
		etcd.StripPrefix(true),
	)
	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	if err := conf.Load(source); err != nil {
		return nil, err
	}
	return &ConfWatch{conf: conf}, nil
}

func (c *ConfWatch) Watch(val interface{}, path ...string) error {
	if err := c.Get(path...).Scan(&val); err != nil {
		return err
	}
	w, err := c.conf.Watch(path...)
	if err != nil {
		return err
	}
	go func() {
		for {
			v, err := w.Next()
			if err != nil {
				logger.Errorf("watch [%s] next error:%s", path, err.Error())
				continue
			}
			if err := v.Scan(&val); err != nil {
				logger.Errorf("scan [%s] next error:%s", path, err.Error())
				continue
			}
			logger.Debugf("scan [%s] new value:%v", path, val)
		}
	}()
	return nil
}

func (c *ConfWatch) WatchCallback(cb func(b []byte) error, path ...string) error {
	b := c.Get(path...).Bytes()
	if err := cb(b); err != nil {
		return err
	}
	w, err := c.conf.Watch(path...)
	if err != nil {
		return err
	}
	go func() {
		for {
			v, err := w.Next()
			if err != nil {
				logger.Errorf("watch [%s] next error:%s", path, err.Error())
				continue
			}
			if err := cb(v.Bytes()); err != nil {
				logger.Errorf("scan [%s] next error:%s", path, err.Error())
				continue
			}
			logger.Debugf("scan [%s] new value:%s", path, string(v.Bytes()))
		}
	}()
	return nil
}

func (c *ConfWatch) Get(path ...string) reader.Value {
	return c.conf.Get(path...)
}

func (c *ConfWatch) Stop() error {
	return c.conf.Close()
}
