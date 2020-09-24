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
