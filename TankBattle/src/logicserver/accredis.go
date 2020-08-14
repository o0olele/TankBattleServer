package main

import (
	"base/env"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
)

const global_id string = "GLOBALID"

type AccRedis struct {
	pool *redis.Pool
}

var mAccRedis *AccRedis

func AccRedis_GetMe() *AccRedis {
	if nil == mAccRedis {
		mAccRedis = &AccRedis{}
		mAccRedis.pool = &redis.Pool{
			MaxIdle:     64,
			MaxActive:   1000,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", env.Get("logic", "redis"))
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

	return mAccRedis
}

func (this *AccRedis) GetIncID() (uint32, error) {
	conn := this.pool.Get()
	defer conn.Close()

	r, err := conn.Do("INCR", global_id)
	if nil != err {
		glog.Error("[Redis] Global Id Inc Fail")
		return 0, err
	}

	value, err := redis.Int(r, nil)
	if nil != err {
		glog.Error("[Redis] Global Id Parse Fail")
		return 0, err
	}

	return uint32(value), nil
}
