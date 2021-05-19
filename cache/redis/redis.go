package redis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewClient(addr string, db int, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          db,
		Password:    password,
		DialTimeout: time.Second * 30,
	})
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	if res := client.Ping(ctx); res.Err() != nil {
		return nil, res.Err()
	}
	return client, nil
}

func NewCluster(addrs []string) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       addrs,
		DialTimeout: time.Second * 30,
	})
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	if res := client.Ping(ctx); res.Err() != nil {
		return nil, res.Err()
	}
	return client, nil
}

func GetByKey(ctx context.Context, key string) (string, error) {
	code, err := cli.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return code, nil
}

func SaveByKey(ctx context.Context, key string, value interface{}) error {
	_, err := cli.Set(ctx, key, value, SaveExpire).Result()
	return err
}

func DelByKey(ctx context.Context, key string) error {
	_, err := cli.Del(ctx, key).Result()
	return err
}

func SetExpire(ctx context.Context, key string, hour int32) error {
	if err := cli.Expire(ctx, key, time.Duration(hour)*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func TTl(ctx context.Context, key string) (time.Duration, error) {
	return cli.TTL(ctx, key).Result()
}

func Setnx(ctx context.Context, key string, value int64) (bool, error) {
	return cli.SetNX(ctx, key, value, 5*time.Minute).Result()
}

//https://blog.csdn.net/lihao21/article/details/49104695
func Lock(ctx context.Context, key string, expire int64) (bool, error) {
	if expire == 0 {
		expire = SetnxExpire
	}
	currentTime := time.Now().Unix()
	isLock, err := Setnx(ctx, key, currentTime+expire)
	if err != nil {
		return false, fmt.Errorf("key: %s, Setnx err: %s\n", key, err.Error())
	}
	if !isLock {
		keyValue, err := GetByKey(ctx, key)
		if err != nil {
			return false, fmt.Errorf("key: %s, GetByKey err: %s\n", key, err.Error())
		}
		//此时不存在，说明已失效或被释放，需再次主动请求
		if len(keyValue) == 0 {
			return Lock(ctx, key, expire)
		}
		lockTime, err := strconv.Atoi(keyValue)
		if err != nil {
			return false, fmt.Errorf("key: %s, strconv.Atoi err: %s\n", key, err.Error())
		}
		//锁已过期，重置
		if int64(lockTime) < currentTime {
			if err := Unlock(ctx, key); err != nil {
				return false, fmt.Errorf("key: %s, Unlock err: %s\n", key, err.Error())
			}
			isLock, err = Setnx(ctx, key, currentTime+expire)
			if err != nil {
				return false, fmt.Errorf("key: %s, again setnx err: %s\n", key, err.Error())
			}
		}
	}

	if isLock {
		if err := SetExpire(ctx, key, int32(expire)); err != nil {
			return false, fmt.Errorf("key: %s, SetExpire err: %s\n", key, err.Error())
		}
		return true, nil
	}

	return false, nil
}

func Unlock(ctx context.Context, key string) error {
	return DelByKey(ctx, key)
}

//设置bit
func SetBit(ctx context.Context, key string, offset int64, value int32) error {
	if _, err := cli.SetBit(ctx, key, offset, int(value)).Result(); err != nil {
		return err
	}
	return nil
}

//获取bit中的偏移量的对应值
func GetBit(ctx context.Context, key string, offset int64) int64 {
	return cli.GetBit(ctx, key, offset).Val()
}

//获取Bitmaps 指定范围值为 1 的位个数
func BitCount(ctx context.Context, key string, startOffset, endOffset int64) int64 {
	return cli.BitCount(ctx, key, &redis.BitCount{
		Start: startOffset,
		End:   endOffset,
	}).Val()
}

//BITOP 命令支持 AND 、 OR 、 NOT 、 XOR 这四种操作中的任意一种参数
//对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destkey 上。
func BitOpAnd(ctx context.Context, key string, keys ...string) (int64, error) {
	return cli.BitOpAnd(ctx, key, keys...).Result()
}

func GetBitAll(ctx context.Context, key string) ([]int, error) {
	info := cli.Get(ctx, key).Val() //返回ascii的字符
	if len(info) == 0 {
		return nil, nil
	}
	//fmt.Printf("ascii字符, redis: %v\n", info) //\x04\x80

	//本质上，数据都是二进制组成的，能不能把 D 转换为 []byte形式，之后通过解析每个 byte所代表的二进制
	//将字符转换成[]byte形式
	origns, err := Bytes(info)
	if err != nil {
		return nil, err
	}
	size := len(origns) * 8
	// buffer := bytes.Buffer{}
	buffer := poolBuffer.Get().(*bytes.Buffer)

	obj := make([]int, size, size)
	//对每个byte转换成二进制
	//参照 ascii码表 一个byte最大的数值就是254，2^8 次方，之后把ascii的字符对应的十进制转换为二进制就可以了
	//fmt.Printf("byte[]: %v\n", origns) //[4 128]
	for _, orign := range origns {
		//fmt.Printf("ascii字符, byte: %v\n", orign) //4  128
		item := strconv.FormatUint(uint64(orign), 2)
		itemLen := len(item)
		//fmt.Printf("ascii字符, byte转二进制: %s\n", item) //100  10000000
		//向前补位
		for i := 0; i < 8-itemLen; i++ {
			buffer.WriteString("0")
		}
		buffer.WriteString(item)
		//fmt.Printf("ascii字符, byte转二进制: %s\n", buffer.String()) //00000100  0000010010000000
	}
	newBit := buffer.String()
	//fmt.Printf("二进制: %s\n", newBit) //0000010010000000
	buffer.Reset()
	poolBuffer.Put(buffer)
	for idx, l := range newBit {
		//处理字符=1的
		if l == 49 {
			//fmt.Printf("二进制:1, ascii字符, byte: %d, bit_key: %d\n", l, idx) //49,5  49,8
			obj[idx] = 1
		} else {
			//fmt.Printf("二进制:0, ascii字符, byte: %d\n", l) //48
		}
	}
	return obj, nil
}

func Bytes(reply interface{}) ([]byte, error) {
	switch _type := reply.(type) {
	case []byte:
		return _type, nil
	case string:
		return []byte(_type), nil
	case nil:
		return nil, errors.New("redigo: nil returned")
	case error:
		return nil, _type
	}
	return nil, fmt.Errorf("redigo: unexpected type for Bytes, got type %T", reply)
}
