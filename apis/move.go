package apis

import (
	"fmt"
	"szinx/core"
	"szinx/pb"

	"github.com/YungMonk/zinx/ziface"
	"github.com/YungMonk/zinx/znet"
	"google.golang.org/protobuf/proto"
)

// MoveAPI 玩家移动的路由业务
type MoveAPI struct {
	znet.BaseRouter
}

// Handle 处理 Connection 主业务的钩子方法 Hook
func (m *MoveAPI) Handle(request ziface.IRequest) {
	// 1.解析客户端传递的proto协议
	positionProtoMsg := &pb.Position{}

	if err := proto.Unmarshal(request.GetData(), positionProtoMsg); err != nil {
		fmt.Println("Move proto unmarshal err:", err)
		return
	}

	// 2.获取当前发送位置信息的是哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Printf("pid not found")
		return
	}
	fmt.Printf(
		"Player pid=%d，move(%f,%f,%f,%f)",
		pid,
		positionProtoMsg.X,
		positionProtoMsg.Y,
		positionProtoMsg.Z,
		positionProtoMsg.V,
	)
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 3.更新当前玩家的坐标，并广播给周边的玩家（九宫格内的玩家）
	player.UpdatePos(
		positionProtoMsg.X,
		positionProtoMsg.Y,
		positionProtoMsg.Z,
		positionProtoMsg.V,
	)
}
