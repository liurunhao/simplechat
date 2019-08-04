package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

//func main() {
//	addr,_:=net.ResolveTCPAddr("tcp4","localhost:8089")     //创建服务器端地址
//	for i :=0;i<5 ;i++  {
//	conn,_:=net.DialTCP("tcp4",nil,addr)                //创建链接
//	conn.Write([]byte("客户数据"+strconv.Itoa(i)))                      //发送数据
//	b:=make([]byte,1024)
//	n,_:=conn.Read(b)                                   //接收服务器端返回的数据
//	fmt.Println("接收到服务器端的返回值为：",string(b[:n]))
//	conn.Close()                        //关闭连接
//	}
//}

type User struct {
	Username  string
	OtherUn   string
	Msg       string
	ServerMsg string
}

var (
	user=new(User)
	wg sync.WaitGroup
)

func main() {
	wg.Add(1)
	fmt.Println("请输入您的账号：")
	fmt.Scanln(&user.Username)
	fmt.Println("请输入您要发送消息的用户：")
	fmt.Scanln(&user.OtherUn)
	addr,_:=net.ResolveTCPAddr("tcp4","localhost:8089")
	conn,_:=net.DialTCP("tcp4",nil,addr)
	go func() {                                                 //发消息的协程
		fmt.Println("请输入您要发送的消息：(只提示一次)")
		for {
			fmt.Scanln(&user.Msg)
			if user.Msg=="exit" {
				conn.Close()
				wg.Done()
				os.Exit(0)
			}
			conn.Write([]byte(fmt.Sprintf("%s-%s-%s-%s",user.Username,user.OtherUn,user.Msg,user.ServerMsg)))
		}
	}()
	go func() {                 //接收消息的协程
		for {
			b := make([]byte, 1024)
			n, _ := conn.Read(b)
			arrey := strings.Split(string(b[:n]), "-")
			user2 := new(User)
			user2.Username = arrey[0]
			user2.OtherUn = arrey[1]
			user2.Msg = arrey[2]
			user2.ServerMsg = arrey[3]
			if user2.ServerMsg != "" {
				fmt.Println("\t\t来自服务器的消息：", user2.ServerMsg)
			} else {
				fmt.Println("\t\t", user2.Username, ":", user2.Msg)
			}
		}
	}()


	wg.Wait()   //保证程序一直在运行状态
}









