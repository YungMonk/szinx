package apis

import (
	"fmt"
	"szinx/core"
	"szinx/pb"

	"github.com/YungMonk/zinx/ziface"
	"github.com/YungMonk/zinx/znet"
	"github.com/golang/protobuf/proto"
)

// WorldChatAPI 世界聊天的路由业务
type WorldChatAPI struct {
	znet.BaseRouter
}

// Handle 处理 Connection 主业务的钩子方法 Hook
func (wc *WorldChatAPI) Handle(request ziface.IRequest) {
	// 1.解析客户端传递的proto协议
	protoMsg := &pb.Talk{}
	if err := proto.Unmarshal(request.GetData(), protoMsg); err != nil {
		fmt.Println("Talk proto unmarshal err:", err)
		return
	}

	// 2.当前的聊天数据是那个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Printf("pid not found")
		return
	}

	// 3.根据pid得到player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 4.将这个消息广播给其它的玩家
	player.Talk(protoMsg.Content)
}
