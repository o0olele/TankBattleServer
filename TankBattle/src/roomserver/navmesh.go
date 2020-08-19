package main

import (
	"base/env"
	"fmt"
	"io/ioutil"
	"unsafe"

	detour "github.com/fananchong/recastnavigation-go/Detour"
)

const (
	NAVMESHSET_MAGIC     int32 = int32('M')<<24 | int32('S')<<16 | int32('A')<<8 | int32('T')
	NAVMESHSET_VERSION   int32 = 1
	TILECACHESET_MAGIC   int32 = int32('T')<<24 | int32('S')<<16 | int32('A')<<8 | int32('T')
	TILECACHESET_VERSION int32 = 1
)

type NavMeshSetHeader struct {
	magic    int32
	version  int32
	numTiles int32
	params   detour.DtNavMeshParams
}

type NavMeshTileHeader struct {
	tileRef  detour.DtTileRef
	dataSize int32
}

func LoadStaticMesh(path string) *detour.DtNavMesh {
	meshData, err := ioutil.ReadFile(path)
	detour.DtAssert(err == nil)

	header := (*NavMeshSetHeader)(unsafe.Pointer(&(meshData[0])))
	detour.DtAssert(header.magic == NAVMESHSET_MAGIC)
	detour.DtAssert(header.version == NAVMESHSET_VERSION)

	navMesh := detour.DtAllocNavMesh()
	state := navMesh.Init(&header.params)
	detour.DtAssert(detour.DtStatusSucceed(state))

	d := int32(unsafe.Sizeof(*header))
	fmt.Println("header tilesize: ", d, header.numTiles)
	for i := 0; i < int(header.numTiles); i++ {
		tileHeader := (*NavMeshTileHeader)(unsafe.Pointer(&(meshData[d])))
		fmt.Println("header tilesize: ", tileHeader.tileRef, tileHeader.dataSize)
		if tileHeader.tileRef == 0 || tileHeader.dataSize == 0 {
			break
		}
		d += int32(unsafe.Sizeof(*tileHeader))

		var t detour.DtTileRef
		data := meshData[d : d+tileHeader.dataSize]
		state = navMesh.AddTile(data, int(tileHeader.dataSize), detour.DT_TILE_FREE_DATA, tileHeader.tileRef, &t)
		detour.DtAssert(detour.DtStatusSucceed(state))
		d += tileHeader.dataSize
		fmt.Println("tile:")
	}

	return navMesh
}

func CreateQuery(mesh *detour.DtNavMesh, maxNode int) *detour.DtNavMeshQuery {
	query := detour.DtAllocNavMeshQuery()
	detour.DtAssert(query != nil)
	status := query.Init(mesh, maxNode)
	detour.DtAssert(detour.DtStatusSucceed(status))
	return query
}

type NavMesh struct {
	mesh   *detour.DtNavMesh
	query  *detour.DtNavMeshQuery
	filter *detour.DtQueryFilter
}

var mNavMesh *NavMesh

func NavMesh_GetMe() *NavMesh {
	if nil == mNavMesh {
		mNavMesh = &NavMesh{}
		mNavMesh.mesh = LoadStaticMesh(env.Get("room", "navmesh"))
		mNavMesh.query = CreateQuery(mNavMesh.mesh, 1024)
		mNavMesh.filter = detour.DtAllocDtQueryFilter()
	}
	return mNavMesh
}
