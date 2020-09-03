package main

import (
	"common"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/golang/glog"
)

func getrand(limit uint32) uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Uint32() % limit
}

func NewRandObstacle() *common.Obstacle {
	pos := common.Pos{
		X: float64(getrand(uint32(common.MapWidth))) - float64(common.MapWidth/2),
		Y: float64(getrand(uint32(common.MapHeight))) - float64(common.MapHeight/2),
	}
	return &common.Obstacle{
		Pos:    pos,
		Height: 1,
		Length: 1,
		Width:  3,
	}
}

func NewObstacle(x, y, length, width, height uint32) *common.Obstacle {
	return &common.Obstacle{
		Pos: common.Pos{
			X: float64(x),
			Y: float64(y),
		},
		Length: length,
		Width:  width,
		Height: height,
	}
}

func GenerateRandMap() *map[uint32]*common.Obstacle {
	obstacle := make(map[uint32]*common.Obstacle)
	for i := 0; i < 3; i++ {
		o := NewRandObstacle()
		o.Id = uint32(i)
		obstacle[uint32(i)] = o
	}
	return &obstacle
}

type obj struct {
	Pos []common.Pos
}

func GenerateMap() *map[uint32]*common.Obstacle {
	file, err := ioutil.ReadFile("../../config/map.json")
	if err != nil {
		glog.Error("[config] Read Map file error")
		return nil
	}
	var pos []common.Pos
	err = json.Unmarshal(file, &pos)
	if err != nil {
		glog.Error("[config] UnMarshal Map file error")
	}
	ret := make(map[uint32]*common.Obstacle)
	for i, p := range pos {
		ret[uint32(i)] = &common.Obstacle{
			Id:     uint32(i),
			Pos:    p,
			Height: 1,
			Width:  1,
			Length: 1,
		}
	}
	return &ret
}
