package redis

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	addr := "10.4.61.59:30379"
	cli, err := NewClient(addr, 1, "")
	if err != nil {
		assert.NoError(t, err)
	} else {
		assert.NotNil(t, cli)
	}
}

func BenchmarkNewClient(b *testing.B) {
	addr := "10.4.61.59:30379"
	cli, err := NewClient(addr, 0, "")
	assert.NoError(b, err)
	assert.NotNil(b, cli)
	ctx := context.TODO()
	err = cli.Set(ctx, "a", 1, time.Minute).Err()
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if v, err := cli.Get(ctx, "a").Result(); err != nil {
			assert.NoError(b, err)
			assert.EqualValues(b, 1, v)
		}
	}
}

func TestNewCluster(t *testing.T) {
	addrs := []string{"10.4.61.59:30379"}
	cli, err := NewCluster(addrs)
	assert.NoError(t, err)
	assert.NotNil(t, cli)
}

func Test_ttL(T *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	if _time, err := TTl(ctx, "lock:seckill"); err != nil {
		T.Fatal(err)
	} else {
		T.Log(_time)
	}
}

func Test_bit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	t.Run("getBit：", func(t *testing.T) {
		key := "bit"
		offset := int64(5)
		fmt.Println(GetBit(ctx, key, offset))
	})

	t.Run("getBitAll：", func(t *testing.T) {
		key := "bit1"
		if res, err := GetBitAll(ctx, key); err != nil {
			t.Error(err)
			return
		} else {
			for k, v := range res {
				fmt.Printf("bit_k: %d,bit_v:%d\n", k, v)
			}
		}
	})

	t.Run("用户签到：", func(t *testing.T) {
		//用户uid
		uid := 1
		//记录有uid的key
		cacheKey := fmt.Sprintf("sign:%d", uid)
		//开始有签到功能的日期
		startDate := "2021-05-01"
		//今天的日期
		todayDate := "2021-05-25"

		//计算offset
		startTime := now.MustParse(startDate).Unix()
		todayTime := now.MustParse(todayDate).Unix()
		offset := int64(math.Ceil(float64(todayTime-startTime) / 86400))

		t.Logf("今天是第 %d 天", offset)
		//签到
		if err := SetBit(ctx, cacheKey, offset, 1); err != nil {
			t.Errorf("sign err: %s", err.Error())
			return
		}
		//查询签到情况
		bitStatus := GetBit(ctx, cacheKey, offset)
		var signRes string
		switch bitStatus {
		case 1:
			signRes = "今天已经签到啦!"
		default:
			signRes = "还没有签到呢!"
		}
		t.Log(signRes)
		//计算总签到次数
		signList, err := GetBitAll(ctx, cacheKey)
		if err != nil {
			t.Errorf("getBitAll failed, key: %s, err: %s", cacheKey, err.Error())
			return
		}
		signCou := BitCount(ctx, cacheKey, 0, int64(len(signList)))
		t.Logf("总计签到 %d 天", signCou)
	})

	t.Run("统计活跃用户：", func(t *testing.T) {
		//日期对应的活跃用户
		data := map[string][]int32{
			"2021-05-10": {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			"2021-05-11": {1, 2, 3, 4, 5, 6, 7, 8},
			"2021-05-12": {1, 2, 3, 4, 5, 6},
			"2021-05-13": {1, 2, 3, 4},
			"2021-05-14": {1, 2},
		}
		key := "activeStat"
		hashTag := "{activeStat}"
		//批量设置活跃状态
		for date, uIds := range data {
			cacheKey := fmt.Sprintf("%s:%s:%s", key, hashTag, date)
			for _, uid := range uIds {
				if err := SetBit(ctx, cacheKey, int64(uid), 1); err != nil {
					t.Errorf("setbit failed, err: %s", err.Error())
					return
				}
			}
		}
		keys := []string{
			"activeStat:{activeStat}:2021-05-13", "activeStat:{activeStat}:2021-05-14", "activeStat:{activeStat}:2021-05-12",
		}
		if _, err := BitOpAnd(ctx, key, keys...); err != nil {
			t.Errorf("BitOpAnd err: %s", err.Error())
			return
		}
		t.Logf("总活跃用户数：%d", BitCount(ctx, key, 0, int64(len(keys))))
	})

	t.Run("用户在线状态：", func(t *testing.T) {
		cacheKey := "userOnline"
		//随机批量设置在线状态
		uIdsMax := 50
		var onlineStatus int32
		for i := 0; i < uIdsMax; i++ {
			onlineStatus = int32(rand.Intn(2))
			if err := SetBit(ctx, cacheKey, int64(i), onlineStatus); err != nil {
				t.Errorf("set online err: %s", err.Error())
				return
			}
		}
		//获取在线情况
		onlineList, err := GetBitAll(ctx, cacheKey)
		if err != nil {
			t.Errorf("getBitAll failed, key: %s, err: %s", cacheKey, err.Error())
			return
		}
		for uid, online := range onlineList {
			t.Logf("uid: %d, onlineStatus: %d", uid, online)
		}
	})

	t.Run("布隆过滤器：", func(t *testing.T) {
		//https://juejin.cn/post/6855839313859461133
	})
}

//------------连接池实现------------
var gClient *redis.Client

func handler(w http.ResponseWriter, r *http.Request) {
	gClient.Ping(context.TODO()).Result()
	printRedisPool(gClient.PoolStats())
	fmt.Fprintf(w, "Hello")
}

func printRedisPool(stats *redis.PoolStats) {
	fmt.Printf("Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d\n",
		stats.Hits, stats.Misses, stats.Timeouts, stats.TotalConns, stats.IdleConns, stats.StaleConns)
}

func printRedisOption(opt *redis.Options) {
	fmt.Printf("Network=%v\n", opt.Network)
	fmt.Printf("Addr=%v\n", opt.Addr)
	fmt.Printf("Password=%v\n", opt.Password)
	fmt.Printf("DB=%v\n", opt.DB)
	fmt.Printf("MaxRetries=%v\n", opt.MaxRetries)
	fmt.Printf("MinRetryBackoff=%v\n", opt.MinRetryBackoff)
	fmt.Printf("MaxRetryBackoff=%v\n", opt.MaxRetryBackoff)
	fmt.Printf("DialTimeout=%v\n", opt.DialTimeout)
	fmt.Printf("ReadTimeout=%v\n", opt.ReadTimeout)
	fmt.Printf("WriteTimeout=%v\n", opt.WriteTimeout)
	fmt.Printf("PoolSize=%v\n", opt.PoolSize)
	fmt.Printf("MinIdleConns=%v\n", opt.MinIdleConns)
	fmt.Printf("MaxConnAge=%v\n", opt.MaxConnAge)
	fmt.Printf("PoolTimeout=%v\n", opt.PoolTimeout)
	fmt.Printf("IdleTimeout=%v\n", opt.IdleTimeout)
	fmt.Printf("IdleCheckFrequency=%v\n", opt.IdleCheckFrequency)
	fmt.Printf("TLSConfig=%v\n", opt.TLSConfig)

}

func Test_pool(t *testing.T) {
	ctx := context.TODO()
	network := "tcp"
	addr := "127.0.0.1:6379"
	gClient = redis.NewClient(&redis.Options{
		//连接信息
		Network:  network, //网络类型，tcp or unix，默认tcp
		Addr:     addr,    //主机名+冒号+端口，默认localhost:6379
		Password: "",      //密码
		DB:       0,       // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//可自定义连接函数
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			netDialer := net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.DialContext(ctx, network, addr)
		},

		//钩子函数
		OnConnect: func(ctx context.Context, cn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
			fmt.Printf("conn=%v\n", cn)
			return nil
		},
	})
	_, _ = gClient.Set(ctx, "poolas", 123, -1).Result()
	get := gClient.Get(ctx, "poolas")
	defer gClient.Close()
	t.Logf("redis get: %+v", get)
	cmd := gClient.Do(ctx, "ping")
	t.Logf("redis cmd: %+v", cmd)
	con := gClient.Conn(ctx)
	t.Logf("redis con: %+v", con)

	printRedisOption(gClient.Options())
	printRedisPool(gClient.PoolStats())

	http.HandleFunc("/", handler)

	_ = http.ListenAndServe(":8080", nil)
}
