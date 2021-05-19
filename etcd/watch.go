package etcd

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/util/log"
)

type EtcdWatchConf struct {
	conf config.Config
}

func NewWatch(addresses []string, prefix string) (*EtcdWatchConf, error) {
	source := etcd.NewSource(
		etcd.WithAddress(addresses...),
		etcd.WithPrefix(prefix),
		etcd.StripPrefix(true),
	)
	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	if err = conf.Load(source); err != nil {
		return nil, err
	}
	return &EtcdWatchConf{conf: conf}, nil
}

func (c *EtcdWatchConf) Watch(val interface{}, path ...string) error {
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
				log.Errorf("watch[%s] next error:%s", path, err.Error())
				continue
			}
			if err = v.Scan(&val); err != nil {
				log.Errorf("scan[%s] next error:%s", path, err.Error())
				continue
			}
			log.Debugf("scan[%s] new value:%v", val)
		}
	}()
	return nil
}

func (c *EtcdWatchConf) WatchCallback(cb func(b []byte) error, path ...string) error {
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
				log.Errorf("watch[%s] next error:%s", path, err.Error())
				continue
			}
			if err = cb(v.Bytes()); err != nil {
				log.Errorf("scan[%s] next error:%s", path, err.Error())
				continue
			}
			log.Debugf("scan[%s] new value:%v", string(v.Bytes()))
		}
	}()
	return nil
}

func (c *EtcdWatchConf) Get(path ...string) reader.Value {
	return c.conf.Get(path...)
}

func (c *EtcdWatchConf) Stop() error {
	return c.conf.Close()
}
