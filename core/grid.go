package core

import (
	"fmt"
	"sync"
)

// Grid 一个AOI地图中格子的类型
type Grid struct {
	// 格子ID
	GID int
	// 格子左边边界坐标
	MinX int
	// 格子右边边界坐标
	MaxX int
	// 格子下边边界坐标
	MinY int
	// 格子上边边界坐标
	MaxY int
	// 当前格子中玩家/物品的ID集合
	playerIDs map[int]bool
	// 保护当前集合的锁
	pIDLock sync.RWMutex
}

// NewGrid 初始化当前格子的方法
func NewGrid(gID, minx, maxx, miny, maxy int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minx,
		MaxX:      maxx,
		MinY:      miny,
		MaxY:      maxy,
		playerIDs: make(map[int]bool),
	}
}

// Add 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true

}

// Remove 从格子删除一个玩家
func (g *Grid) Remove(playerID int) error {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	if _, ok := g.playerIDs[playerID]; !ok {
		return fmt.Errorf("playerID=%d is not found", playerID)
	}

	delete(g.playerIDs, playerID)

	return nil
}

// GetPlayerIDs 获取格子中的所有玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	for pid := range g.playerIDs {
		playerIDs = append(playerIDs, pid)
	}

	return
}

// 调试打印出格子中的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf(
		"Grid id:%d, minX:%d, maxX:%d, miny:%d, maxy:%d, playerids:%+v\n",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs,
	)
}
