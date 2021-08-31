package main

import (
	"context"
	"dm-go/proto"
	"flag"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
	"time"
)

type SlaveInfo struct {
	hostname string
	addr     string
}

var slaves map[string]string

var (
	bind = flag.String("bind", "0.0.0.0:8972", "server address")
)

type DMService int

func (t *DMService) Register(ctx context.Context, args *pb.DMServiceInfo, reply *pb.Void) error {
	_, ok := slaves[args.Addr]
	if ok {
		log.Infof("Slave [%s]'s information updated\n", args.Hostname)
	} else {
		log.Infof("New slave [%s] added into the pool\n", args.Hostname)
	}

	slaves[args.Addr] = args.Hostname
	return nil
}

func main() {
	flag.Parse()
	slaves = make(map[string]string)

	go CreateMasterServer()
	go CheckAvailability()
	time.Sleep(time.Minute * 10) // 持续 10 分钟的服务器
}

func CreateMasterServer() {
	s := server.NewServer()
	s.RegisterName("DM-Master", new(DMService), "")
	s.Serve("tcp", *bind)

	//err := s.UnregisterAll()
	//if err != nil {
	//	panic(err)
	//}
}

func CheckAvailability() {
	for {
		time.Sleep(time.Second * 3)
		log.Infof("Slave pool have %d slaves.\n", len(slaves))
		for addr, _ := range slaves { // because we use "addr" as the key
			if !CheckOnline(addr) {
				delete(slaves, addr)
			}
		}
	}
}

// CheckOnline 尝试调用 slave 的 AreYouAlive 函数
func CheckOnline(addr string) bool {
	option := client.DefaultOption
	option.ConnectTimeout = time.Second * 5 // 连接超时 5 秒

	d, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	xclient := client.NewXClient("DM-Slave", client.Failfast, client.RoundRobin, d, option)
	defer xclient.Close()
	args := &pb.Void{}
	resp := &pb.Response{}
	err := xclient.Call(context.Background(), "AreYouAlive", args, resp)
	if err != nil {
		log.Errorf("Slave dead: %v\n", err)
		return false
	}
	return true
}
