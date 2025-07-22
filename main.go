package main 

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	
)

const listenAdrr = "0.0.0.0"


func main(){
	dest := os.Args[1]
	destaddr ,err := net.ResolveIPAddr("ip4",dest)
	if err != nil {
		fmt.Printf("error resolving destination")
		os.Exit(1)
	}
	conn,err := icmp.ListenPacket("ip4:icmp",listenAdrr)
	if err != nil {
		fmt.Println("error while connecting to icmp packet")
	}

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body :&icmp.Echo{
			ID : os.Getpid() & 0xffff ,
			Seq : 1,
			Data : []byte("hello"),
		},

	}
	msgdata ,err := msg.Marshal(nil)
	if err != nil {
		fmt.Printf("erro will marshaling the icmp %v\n",err)
		os.Exit(1)
	}
	startTime := time.Now()
	if _,err := conn.WriteTo(msgdata,destaddr);err != nil {
		fmt.Printf("error while sending the data %v\n",err)
		os.Exit(1)
	}
	fmt.Printf("Pinging  %s (%s) with %d bytes of data\n",dest,destaddr.String(),len(msgdata))

	replydata := make([]byte,1500)
	err = conn.SetReadDeadline(time.Now().Add(5*time.Second))
	if err != nil {
		fmt.Printf("error setting in readline deadline :%v\n",err)
		os.Exit(1)
	}

	n,peer,err := conn.ReadFrom(replydata)
	if err != nil {
		fmt.Printf("error while getting the result")
		os.Exit(1)
	}
	rtt := time.Since(startTime)

	replymsg,err := icmp.ParseMessage(1,replydata[:n])
	if err != nil {
		fmt.Printf("error while parsing the icmp message repy %v\n",err)
		os.Exit(1)
	}

	switch replymsg.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("Reply from %s: bytes=%d time=%s\n", peer, n, rtt)
	default:
		fmt.Printf("Got unexpected ICMP message: type=%v, code=%v\n", replymsg.Type, replymsg.Code)
	}
}

