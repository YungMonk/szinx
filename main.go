package main

import (
	"fmt"
	"szinx/apis"
	"szinx/core"

	"github.com/YungMonk/zinx/ziface"
	"github.com/YungMonk/zinx/zlog"
	"github.com/YungMonk/zinx/znet"
)

// OnConnectionAdd 当前客户端创建连接之后执行的 Hook 函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)

	// 给客户端发送MsgID=1的消息，同步当前的playerID给客户端
	player.SyncPid()

	// 给客户端发送MsgID=200的消息，同步当前player的位置给客户端
	player.BroadCastStartPosition()

	// 将新上线的玩家添加到世界管理模块中
	core.WorldMgrObj.AddPlayer(player)

	// 将当前连接绑定到一个Pid玩家ID的属性
	conn.SetProperty("pid", player.Pid)

	// 在当前玩家上线之后，触发同步当前玩家位置信息（告知周围玩家当前玩家已经上线）
	player.SyncSurrounding()

	fmt.Println("\n====> Player pid=", player.Pid, " is arrived ====")
}

// OnConnectionLost 当前客户端断开连接之前执行的 Hook 函数
func OnConnectionLost(conn ziface.IConnection) {
	// 获取当前连接绑定的玩家 ID
	pid, _ := conn.GetProperty("pid")

	// 获取当前连接周边的玩家信息
	player := core.WorldMgrObj.Players[pid.(int32)]

	// 触发玩家下线的业务
	player.Offline()

	fmt.Printf("====> Player pid=%d will offline <====", pid)
}

func main() {
	zlog.SetLevel(zlog.LogDebug)

	// 1.创建Server句柄，使用zinx的api
	s := znet.NewServer("[zinx.v0.5]")

	// 2.注册连接 Hook 钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 3.给服务注册路由
	s.AddRouter(2, &apis.WorldChatAPI{})
	s.AddRouter(3, &apis.MoveAPI{})

	// 4.启动Server
	s.Serve()
}
