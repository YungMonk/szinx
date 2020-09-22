package main

import (
	"fmt"

	"github.com/YungMonk/zinx/ziface"
	"github.com/YungMonk/zinx/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle 处理 Connection 主业务的钩子方法 Hook
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router handle ...")

	// 先读取客户端的数据， 再处理 ping... 数据
	fmt.Println("recv from clinet: msgID = ",
		request.GetMsgID(),
		", Data = ",
		string(request.GetData()),
	)

	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...")); err != nil {
		fmt.Println("send msg error ", err)
	}
}

// HelloZinxRouter 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Handle 处理 Connection 主业务的钩子方法 Hook
func (pr *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter handle ...")

	// 先读取客户端的数据， 再处理 ping... 数据
	fmt.Println("recv from clinet: msgID = ",
		request.GetMsgID(),
		", Data = ",
		string(request.GetData()),
	)

	if err := request.GetConnection().SendMsg(201, []byte("hello...hello...")); err != nil {
		fmt.Println("send msg error ", err)
	}
}

// DoConnectionBegin 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionBegin is Called.")
	if err := conn.SendMsg(202, []byte("Do Connnection is Begin")); err != nil {
		fmt.Println(err)
	}

	// 给当前的链接设置一些属性
	fmt.Println("Set conn Name, hee ...")
	conn.SetProperty("Name", "YungMonk")
	conn.SetProperty("Home", "https://github.com/YungMonk/zinx")
}

// DoConnncetionLost 断开连接之前执行的钩子函数
func DoConnncetionLost(conn ziface.IConnection) {
	fmt.Println("===> DoConnncetionLost is Called.")
	fmt.Println("conn ID=", conn.GetConnID(), "is lost...")

	if value, err := conn.GetProperty("Name"); err == nil {
		fmt.Printf("Name=%s\n", value)
	}

	if value, err := conn.GetProperty("Home"); err == nil {
		fmt.Printf("Home=%s\n", value)
	}
}

func main() {
	// 1.创建Server句柄，使用zinx的api
	s := znet.NewServer("[zinx.v0.5]")

	// 2.注册连接 Hook 钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnncetionLost)

	// 3.给服务注册路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloZinxRouter{})

	// 4.启动Server
	s.Serve()
}
