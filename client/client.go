package main

import (
	"context"
	"dm-go/proto"
	"flag"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"log"
	"os"
	"time"
)

var (
	serveraddr  = flag.String("serveraddr", "localhost:8972", "server address")
	serviceaddr = flag.String("serviceaddr", "localhost:9967", "server address")
)

type DMHeartBeat int

func (t *DMHeartBeat) AreYouAlive(ctx context.Context, args *pb.Void, reply *pb.Response) error {
	reply.Code = 0
	reply.Msg = ""
	return nil
}

func main() {
	flag.Parse()

	go RegisterSlave()
	go HeartBeatService()

	time.Sleep(time.Minute * 1)
}

func RegisterSlave() {
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+*serveraddr, "")

	for {
		xclient := client.NewXClient("DM-Master", client.Failtry, client.RoundRobin, d, client.DefaultOption)

		hostname, err := os.Hostname() // hostname may change, so we put it in the loop
		args := &pb.DMServiceInfo{
			Name:     "Core",
			Desc:     "DM core service",
			Addr:     *serviceaddr,
			Hostname: hostname,

			Major: 1,
			Minor: 0,
			Patch: 0,
		}

		reply := &pb.Void{}
		err = xclient.Call(context.Background(), "Register", args, reply)
		if err != nil {
			log.Printf("failed to call: %v\n", err)
		}

		time.Sleep(time.Second * 10)
	}

}

func HeartBeatService() {
	s := server.NewServer()
	_ = s.RegisterName("DM-Slave", new(DMHeartBeat), "")
	s.Serve("tcp", *serviceaddr)
	time.Sleep(time.Minute)
	err := s.UnregisterAll()
	if err != nil {
		panic(err)
	}
}
