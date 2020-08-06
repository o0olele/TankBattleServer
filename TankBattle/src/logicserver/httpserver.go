package main

import (
	"base/env"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
)

var pool *redis.Pool
var listkey = "ListKey"

// 连接池
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     64,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func RecordHandler(w http.ResponseWriter, r *http.Request) {
	values := r.FormValue("text")
	if values == "" {
		fmt.Fprintf(w, "try to use \"/echo?text=xxx!\"")
		return
	}

	fmt.Fprintf(w, values)

	temp_values := fmt.Sprintf("%s:%s:%s", r.RemoteAddr, GetDateFormat(), values)

	if pool != nil {
		conn := pool.Get()
		defer conn.Close()

		conn.Do("LPUSH", listkey, temp_values)
	}
}

func EchoAllHandler(w http.ResponseWriter, r *http.Request) {
	if pool != nil {
		conn := pool.Get()
		defer conn.Close()

		listLen, err := redis.Int(conn.Do("LLEN", listkey))
		if err != nil {
			fmt.Println("Redis LLEN error!")
			return
		}

		values, err := redis.MultiBulk(conn.Do("LRANGE", listkey, 0, listLen))
		items, err := redis.Strings(values, nil)
		for _, s := range items {
			fmt.Fprintf(w, s+"\n")
		}

		LogicGrpcClient_GetMe().Send("123456789")

	}

}

func ClearAllHandler(w http.ResponseWriter, r *http.Request) {
	if pool != nil {
		conn := pool.Get()
		defer conn.Close()

		_, err := conn.Do("DEL", listkey)
		if err != nil {
			fmt.Println("Redis DEL error!")
			return
		} else {
			fmt.Fprintf(w, "clear success!")
		}
	}
}

// 时间戳转年月日 时分秒
func GetDateFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func StartHttpServer() bool {
	pool = newPool("localhost:6379", "")

	http.HandleFunc("/add", RecordHandler)
	http.HandleFunc("/all", EchoAllHandler)
	http.HandleFunc("/clean", ClearAllHandler)

	addr := env.Get("logic", "listen")
	ln, err := net.Listen("tcp", addr)
	if nil != err {
		glog.Error("[Start] Bind Port Error, Port=", addr, ",", err)
		return false
	}

	go http.Serve(ln, nil)
	glog.Info("[Start] Bind Port Success, Port=", addr)

	return true
}
