package core

import "sync"

// WorldManager 当前世界总管理模块
type WorldManager struct {
	// AOIManager 当前世界地图的 AOI 管理模块
	AoiManager *AOIManager

	// 当前全部在线的 Players 集合
	Players map[int32]*Player

	// 保护 Players 的锁
	pLock sync.RWMutex
}

// WorldMgrObj 提供一个对外的世界管理模块句柄（全局）
var WorldMgrObj *WorldManager

// 初始化世界管理模块
func init() {
	WorldMgrObj = &WorldManager{
		// 创建世界
		AoiManager: NewAOIManager(
			AOIMINX,
			AOIMAXX,
			AOICNTX,
			AOIMINY,
			AOIMAXY,
			AOICNTY,
		),
		// 初始化 Players 集合
		Players: make(map[int32]*Player),
	}
}

// AddPlayer 添加一个 Player
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	// 将 Player 添加到 AOIManager 中
	wm.AoiManager.AddPidToGridByPos(int(player.Pid), player.X, player.Z)
}

// RemovePlayerByPid 删除一个 Player
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	// 取得当前玩家
	player := wm.Players[pid]

	// 将 Player 从 AOIManger 中移除
	wm.AoiManager.RemovePidFromGridByPos(int(pid), player.X, player.Z)

	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()

}

// GetPlayerByPid 通过玩家ID查询player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) (player *Player) {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

// GetAllPlayers 获取全部在线玩家
func (wm *WorldManager) GetAllPlayers() (players []*Player) {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players = make([]*Player, 0)
	for _, player := range wm.Players {
		players = append(players, player)
	}

	return players
}
