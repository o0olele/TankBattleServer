package main

import (
	"common"
	"math/rand"
	"time"
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
