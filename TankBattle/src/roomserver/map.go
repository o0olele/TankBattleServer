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
		X: float64(getrand(uint32(common.MapHeight / 2))),
		Y: float64(getrand(uint32(common.MapWidth / 2))),
	}
	return &common.Obstacle{
		Pos:    pos,
		Height: 1,
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
		obstacle[uint32(i)] = o
	}
	return &obstacle
}
