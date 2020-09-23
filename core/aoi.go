package core

import (
	"fmt"
)

// AOIManager AOI 区域管理模块
type AOIManager struct {
	// 区域的左边边界坐标
	MinX int
	// 区域的右边边界坐标
	MaxX int
	// 区域的下边边界坐标
	MinY int
	// 区域的上边边界坐标
	MaxY int
	// X 方向格子数量
	CntsX int
	// Y 方向格子数量
	CntsY int
	// 当前区域中有那些格子 map-key=格子的ID，value=格子对象
	Grids map[int]*Grid
}

// NewAOIManager 初始化 AOI 区域管理模块的方法
func NewAOIManager(minx, maxx, cntsx, miny, maxy, cntsy int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minx,
		MaxX:  maxx,
		MinY:  miny,
		MaxY:  maxy,
		CntsX: cntsx,
		CntsY: cntsy,
		Grids: make(map[int]*Grid),
	}

	// 给 AOI 初始化区域的格子的所有的格子进行编号和初始化
	for y := 0; y < cntsy; y++ {
		for x := 0; x < cntsx; x++ {
			// 计算格子的ID，根据x，y编号
			// 格子编号 = y*cntsx + x
			gid := y*cntsx + x
			aoiMgr.Grids[gid] = NewGrid(
				gid,
				aoiMgr.MinX+x*aoiMgr.gridWith(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWith(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength(),
			)
		}
	}

	return aoiMgr
}

// 得到每个格子在x方向的宽度
func (am *AOIManager) gridWith() int {
	return (am.MaxX - am.MinX) / am.CntsX
}

// 得到每个格子在y方向的长度
func (am *AOIManager) gridLength() int {
	return (am.MaxY - am.MinY) / am.CntsY
}

// GetSurroundGridsByGid 根据当前格子的 gid 获取其周边格子
func (am *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	// 当前格子是否在AOIManager中
	if _, ok := am.Grids[gid]; !ok {
		return nil
	}

	// 初始化 grids 返回值
	grids = append(grids, am.Grids[gid])

	// 1.通过当前的 gid 计算其 x 的编号，x=gid%CntX
	x := gid % am.CntsX

	// 2.判断 gid 的左边是否有格子，如果有放入 grids 中
	if x > 0 {
		grids = append(grids, am.Grids[gid-1])
	}

	// 3.判断 gid 的右边是否有格子，如果有放入 grids 中
	if x < am.CntsX-1 {
		grids = append(grids, am.Grids[gid+1])
	}

	// 4.将 x 轴中的格子全部取出，放入 gidsX中
	gidsX := make([]int, 0, len(grids))
	for _, grid := range grids {
		gidsX = append(gidsX, grid.GID)
	}

	//  5.然后遍历 gidsX 中的格子 gid'
	for _, gidd := range gidsX {
		// 当前格子gid'的y轴的编号，y=gid'/CntX
		y := gidd / am.CntsX
		// 判断 gid' 的上边是否有格子，如果有放入 grids
		if y > 0 {
			grids = append(grids, am.Grids[gidd-am.CntsX])
		}

		// 判断 gid' 的下边是否有格子，如果有放入 grids
		if y < am.CntsY-1 {
			grids = append(grids, am.Grids[gidd+am.CntsX])
		}
	}

	return grids
}

// GetGidByPos 通过 x，y来获取格子的gid
func (am *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - am.MinX) / am.CntsX
	idy := (int(y) - am.MinY) / am.CntsY

	return idy*am.CntsX + idx
}

// GetPidsByPos 通过横纵坐标获取周边九宫格内的所有 playerIDs
func (am *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	// 得到当前坐标的gid
	gid := am.GetGidByPos(x, y)

	// 通过gid得到gids
	gids := am.GetSurroundGridsByGid(gid)

	// 通完gids得到所有的playerIDs
	for _, Grid := range gids {
		playerIDs = append(playerIDs, Grid.GetPlayerIDs()...)
		fmt.Printf("===> [Grid] ID:%d, Pids:%+v", Grid.GID, Grid.GetPlayerIDs())
	}

	return playerIDs
}

// AddPidToGrid 给格子中添加一个 playerid
func (am *AOIManager) AddPidToGrid(pid, gid int) {
	am.Grids[gid].Add(pid)
}

// RemovePidFromGrid 从格子中移除一个 playerid
func (am *AOIManager) RemovePidFromGrid(pid, gid int) {
	am.Grids[gid].Remove(pid)
}

// GetPidsByGid 从格子中取出所有 playerids
func (am *AOIManager) GetPidsByGid(gid int) (playerids []int) {
	return am.Grids[gid].GetPlayerIDs()
}

// AddPidToGridByPos 通过坐标添加一个 playerid 到格子中
func (am *AOIManager) AddPidToGridByPos(pid int, x, y float32) {
	gid := am.GetGidByPos(x, y)
	am.Grids[gid].Add(pid)
}

// RemovePidFromGridByPos 通过坐标移除格子中的一个 playerid
func (am *AOIManager) RemovePidFromGridByPos(pid int, x, y float32) {
	gid := am.GetGidByPos(x, y)
	am.Grids[gid].Remove(pid)
}

// 打印格子信息
func (am *AOIManager) String() string {
	// 打印 AOIManager 信息
	str := fmt.Sprintf(
		"AOIManager:\nMinX:%d, MaxX:%d, MinY:%d, MaxY:%d, CntsX:%d, CntsY:%d\n Grids:\n",
		am.MinX, am.MaxX, am.MinY, am.MaxY, am.CntsX, am.CntsY,
	)

	// 打印所有的 Grid 信息
	for _, gval := range am.Grids {
		str += fmt.Sprintln(gval)
	}

	return str
}
