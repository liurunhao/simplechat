package main

import (
	"fmt"
	"net"
	"strings"
)

//func main() {
//	addr,_:=net.ResolveTCPAddr("tcp4","localhost:8089")  //创建服务器地址
//	lis,_:=net.ListenTCP("tcp4",addr)                             //创建监听器
//	fmt.Println("服务器已经启动")
//	for {                       //加死循环让服务器端一直接收
//		conn, _ := lis.Accept() //通过监听器接收客户端发来的数据，阻塞式的，接收不到就阻塞
//		go func() {             //在服务器端加协程，即可解决并发过程中效率的问题
//		b := make([]byte, 1024) //因为连接读取的返回值是byte类型的字符串，所以定义一个
//		n, _ := conn.Read(b)    //转换数据,n是长度
//		fmt.Println("获取到的数据为：", string(b [:n]))
//		conn.Write(append([]byte("server:"), b[:n]...)) //向客户端返回数据
//		conn.Close()                                    //关闭连接
//		}()
//	}
//}

type User struct {
	Username  string
	OtherUn   string
	Msg       string
	ServerMsg string
}

var (
	userMap=make(map[string] net.Conn)          //map中放conn对象
	user = new(User)
)

func main() {
	addr,_:=net.ResolveTCPAddr("tcp4","localhost:8089")
	lis,_:=net.ListenTCP("tcp4",addr)
	fmt.Println("服务器已经启动")
	for {
		conn, _ := lis.Accept()
		go func() {
			for{                          //建立死循环让他们可以一直发消息
				b := make([]byte, 1024)
				n, _ := conn.Read(b)
				arrey:=strings.Split(string(b[:n]),"-")   //将接收到的数据转化成字符串然后拆分，用arrey接收
				user.Username = arrey[0]   //把客户端传过来的数据解析成user对象
				user.OtherUn  = arrey[1]
				user.Msg      = arrey[2]
				user.ServerMsg= arrey[3]
				userMap[user.Username]=conn

				if v ,ok:=userMap[user.OtherUn];ok&&v!=nil {
					n,err:=v.Write([]byte(fmt.Sprintf("%s-%s-%s-%s",user.Username,user.OtherUn,user.Msg,user.ServerMsg)))
					if n<=0||err!=nil {
						delete(userMap,user.OtherUn)    //如果有错就把map里的其他用户删掉，然后关闭本次连接
						conn.Close()
						v.Close()
						break
					}
				}else {
					user.ServerMsg="对方不在线"
					conn.Write([]byte(fmt.Sprintf("%s-%s-%s-%s",user.Username,user.OtherUn,user.Msg,user.ServerMsg)))
				}
			}
		}()
	}
}