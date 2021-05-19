package etcd

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pquerna/ffjson/ffjson"
)

var (
	ClientV3Config clientv3.Config
	ClientV3       *clientv3.Client
	err            error
)

func NewKV(addresses []string) error {
	ClientV3Config = clientv3.Config{
		Endpoints:   addresses,
		DialTimeout: 5 * time.Second,
	}
	if ClientV3, err = clientv3.New(ClientV3Config); err != nil {
		return err
	}
	return nil
}

// SetKV 新增或修改指定字段值，key不传则在path下直接保存val
func SetKV(path string, val interface{}, assignKey ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var saveVal string
	switch val.(type) {
	case string: //需为json字符串
		saveVal = val.(string)
	default:
		saveVal = genPutVal(ctx, path, val, assignKey...)
	}
	if len(saveVal) == 0 {
		return nil
	}

	if _, err = ClientV3.Put(ctx, path, saveVal); err != nil {
		return err
	}
	return nil
}

// GetKV 获取path下指定key的值，key不传则直接获取path下数据
func GetKV(path string, assignKey ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getResp, err := ClientV3.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return []byte(""), nil
	}

	if len(assignKey) > 0 && len(assignKey[0]) > 0 {
		var dataMap map[string]interface{}
		if err = json.Unmarshal(getResp.Kvs[0].Value, &dataMap); err != nil {
			return []byte(""), err
		}
		if v, ok := dataMap[assignKey[0]]; !ok {
			return []byte(""), nil
		} else {
			bv, _ := ffjson.Marshal(v)
			return bv, nil
		}
	}

	return getResp.Kvs[0].Value, nil
}

// DelKV 删除path下指定key，key不传则直接删除path
func DelKV(path string, assignKey ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(assignKey) == 0 {
		_, err = ClientV3.Delete(ctx, path)
		return err
	}

	getResp, err := ClientV3.Get(ctx, path)
	if err != nil {
		return err
	}
	if len(getResp.Kvs) == 0 {
		return nil
	}

	var dataMap map[string]interface{}
	if err = json.Unmarshal(getResp.Kvs[0].Value, &dataMap); err != nil {
		return err
	}
	delete(dataMap, assignKey[0])

	if len(dataMap) == 0 {
		_, err = ClientV3.Delete(ctx, path)
		return err
	}

	bd, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}
	if _, err = ClientV3.Put(ctx, path, string(bd)); err != nil {
		return err
	}

	return nil
}

func genPutVal(ctx context.Context, path string, val interface{}, assignKey ...string) string {
	if len(assignKey) == 0 || len(assignKey[0]) == 0 {
		b, _ := ffjson.Marshal(val)
		return string(b)
	}
	getResp, err := ClientV3.Get(ctx, path)
	if err != nil {
		return ""
	}

	if len(getResp.Kvs) == 0 {
		dataMap := map[string]interface{}{assignKey[0]: val}
		bd, _ := json.Marshal(dataMap)
		return string(bd)
	}

	var dataMap map[string]interface{}
	if err = json.Unmarshal(getResp.Kvs[0].Value, &dataMap); err != nil {
		return ""
	}
	// 存在即修改，不存在则新增
	dataMap[assignKey[0]] = val
	bd, _ := json.Marshal(dataMap)
	return string(bd)
}
