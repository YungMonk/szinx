package main

import (
	"fmt"
	"szinx/apis"
	"szinx/core"

	"github.com/YungMonk/zinx/ziface"
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

	fmt.Println("====> Player pid=", player.Pid, " is arrived ====")
}

func main() {
	// 1.创建Server句柄，使用zinx的api
	s := znet.NewServer("[zinx.v0.5]")

	// 2.注册连接 Hook 钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	// s.SetOnConnStop(DoConnncetionLost)

	// 3.给服务注册路由
	s.AddRouter(2, &apis.WorldChatAPI{})
	s.AddRouter(3, &apis.MoveAPI{})

	// 4.启动Server
	s.Serve()
}
