package main

import "common"

type PlayerInfo struct {
	id     uint32
	name   string
	friend map[uint32]*PlayerInfo
	rank   uint32
	pos    common.Pos
	HP     uint32
	MovDir uint32
}
