package core

import (
	"fmt"
	"math/rand"
	"sync"

	"szinx/pb"

	"github.com/YungMonk/zinx/ziface"
	"github.com/golang/protobuf/proto"
)

// Player 玩家对象
type Player struct {
	Pid  int32              // 玩家 id
	Conn ziface.IConnection // 当前玩家的连接（用于和客户端的连接）
	X    float32            // 平面的 x 坐标
	Y    float32            // 高度
	Z    float32            // 平面的 y 坐标
	V    float32            // 玩家的旋转的角度（0-360）
}

// PIDGen PlayerID 生成器
var PIDGen int32 = 1 // 用来生成玩家 id 的计数器
// IDLock 保护 PIDGen 的 Mutex
var IDLock sync.Mutex

// NewPlayer 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家 ID
	IDLock.Lock()
	id := PIDGen
	PIDGen++
	IDLock.Unlock()

	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机在160坐标点，基于平面x轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), // 随机在140坐标点，基于平面y轴若干偏移
		V:    0,                            // 角度为0
	}
}

// SendMsg 提供一个发送给客户端消息的方法
// 主要是将pb的protobuf数据序列化后，再调用zinx的SendMsg方法
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	// 将proto Message结构体数据序列化，转化为二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Printf("marshal data error:%s\n", err)
		return
	}

	// 将转化后的二进制文件通过zinx框架的SendMsg方法发送给客户端
	if p.Conn == nil {
		fmt.Printf("connection in player is nil\n")
		return
	}
	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Printf("player send msg error:%s", err)
		return
	}

	return
}

// SyncPid 告知客户端玩家Pid，同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	// 组建 MsgID:1 的 proto 数据
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 将消息发送给客户端
	p.SendMsg(1, protoMsg)
}

// BroadCastStartPosition 广播玩家的上线地点
func (p *Player) BroadCastStartPosition() {
	// 组建 MsgID:200 的 proto 数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.Z,
			},
		},
	}

	// 将消息发送给客户端
	p.SendMsg(200, protoMsg)
}

// Talk 玩家广播聊天消息到世界
func (p *Player) Talk(content string) {
	// 组建 MsgID:200 的 proto 数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 得到当前所有在线的玩家
	players := WorldMgrObj.GetAllPlayers()

	// 向所有玩家（包括自己）发送 MsgID:200 消息
	for _, player := range players {
		// player 分别给对应的客户端发送消息
		player.SendMsg(200, protoMsg)
	}
}

// SyncSurrounding 在当前玩家上线之后，触发同步当前玩家位置信息（告知周围玩家当前玩家已经上线）
func (p *Player) SyncSurrounding() {
	// 1.获取当前玩家周围的玩家有哪些（九宫格）
	players := p.GetSurroundingPlayers()

	// 2.将当前的玩家信息通过 MsgID:200 发给周围的玩家（让其它玩家看到当前玩家）
	// 2.1 组建 MsgID:200 的 proto 数据
	broadCastProtoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.Z,
			},
		},
	}
	// 2.2 分别给周围玩家的客户端发送消息为200的信息 broadCastProtoMsg
	for _, player := range players {
		player.SendMsg(200, broadCastProtoMsg)
	}

	// 3.将周围的玩家位置信息发送给当前玩家 MsgID:202（让当前玩家看到周围的玩家）
	// 3.1 组建 MsgID:202 的 proto 数据
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		pbPlayer := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersProtoMsg = append(playersProtoMsg, pbPlayer)
	}
	syncProtoMsg := &pb.SyncPlayer{
		Ps: playersProtoMsg[:],
	}

	// 3.2 将组建好的数据发送给当前玩家的客户端
	p.SendMsg(202, syncProtoMsg)
}

// UpdatePos 更新当前玩家的坐标（广播玩家当前位置的移动信息）
func (p *Player) UpdatePos(x, y, z, v float32) {
	IDLock.Lock()
	p.X, p.Y, p.Z, p.V = x, y, z, v
	IDLock.Unlock()

	// 给其它玩家广播当前玩家位置变动信息
	broadcastProtoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 获取当前玩家周边的玩家
	players := p.GetSurroundingPlayers()

	// 给周围的玩家发送位置变动信息
	for _, pler := range players {
		pler.SendMsg(200, broadcastProtoMsg)
	}
}

// GetSurroundingPlayers 获取当前玩家周围（九宫格内）的玩家信息
func (p *Player) GetSurroundingPlayers() []*Player {
	pids := WorldMgrObj.AoiManager.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

// Offline 玩家下线
func (p *Player) Offline() {
	// 获取当前玩家周边九宫格内的玩家信息
	players := p.GetSurroundingPlayers()

	// 给周边的玩家发送 MsgID:201 信息
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, protoMsg)
	}

	// 将当前玩家从AOI管理器删除
	WorldMgrObj.AoiManager.RemovePidFromGridByPos(int(p.Pid), p.X, p.Z)

	// 将当前玩家从世界管理器删除
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
