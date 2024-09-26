package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string // 每个用户都有一个chan
	conn net.Conn    // 当前用户和对端用户通信的连接

	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	// 从连接中获取远程地址
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessge()

	return user
}

// 用户上线业务
func (this *User) Online() {
	// 1.当前用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock() // map是引用类型，加锁
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 2.当前用户上线，广播消息
	this.server.BroadCast(this, "已上线")
}

// 用户下线业务
func (this *User) Offline() {
	// 1.当前用户下线，将用户到从onlineMap中删除
	this.server.mapLock.Lock() // map是引用类型，加锁
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 2.当前用户上线，广播消息
	this.server.BroadCast(this, "下线")

}

// 给当前User对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户都有哪些
		this.server.mapLock.Lock()
		for _, u := range this.server.OnlineMap {
			onlineMsg := "[" + u.Addr + "]" + u.Name + ":" + "在线... \n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式：rename|张三	修改用户名为张三
		newName := strings.Split(msg, "|")[1]

		// 判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名被使用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)

			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("您已经更新用户名:" + this.Name + "\n")

		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式： to|张三|消息内容
		// 获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用 \"to|张三|你好啊\"格式。\n")
			return
		}
		// 根据用户名 得到对方的User对象
		remeteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}

		// 获取消息内容，通过对方的User对象将消息内容发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("无消息内容，请重发\n")
			return
		}

		remeteUser.SendMsg(this.Name + "对您说:" + content + "\n")
	} else {
		this.server.BroadCast(this, msg)
	}

}

// 监听当前User channel的方法，一旦有消息，就发送给对端客户端
func (this *User) ListenMessge() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
