package main

import (
	"runtime"
	"os/exec"
	"log"
	"flag"
	"net"
	"time"
	"strconv"
)

// UpdateDate can update system date.
//
// You should provide args like 2017/04/30, is year/month/day
func UpdateDate(date string) error {
	cmd := exec.Command("cmd", "/c", "date", date)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// UpdateTime can update system time.
//
// You should provide args like 13:45:01.23, is hour:mine:second
func UpdateTime(time string) error {
	cmd := exec.Command("cmd", "/c", "time", time)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func server(serverAddr string) {
	ServerAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Println(err)
	}
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error: ", err)
		}
		resp := string(buf[0:n])
		log.Println("Received ", resp, " from ", addr)
		if resp != "TimeSync" {
			continue
		}
		now := time.Now().Unix()
		ServerConn.WriteToUDP([]byte(strconv.FormatInt(now, 10)), addr)
	}
}

func main() {
	if runtime.GOOS != "windows" {
		panic("Not Windows")
	}
	taip := flag.String("type", "client", "[server|client]")
	serverAddr := flag.String("ServerAddr", "192.168.123.155:2345", "server listen addr, or server addr")
	clientAddr := flag.String("ClientAddr", "192.168.123.155:0", "client addr")
	//netType:= flag.String("netType", "udp", "[tcp|udp]")
	flag.Parse()
	log.Println(*taip)
	switch *taip {
	case "server":
		server(*serverAddr)
	case "client":
		ServerAddr, err := net.ResolveUDPAddr("udp", *serverAddr)
		if err != nil {
			log.Println(err)
		}
		LocalAddr, err := net.ResolveUDPAddr("udp", *clientAddr)
		if err != nil {
			log.Println(err)
		}
		Conn, err := net.ListenUDP("udp", LocalAddr)
		defer Conn.Close()
		if err != nil {
			log.Println(err)
		}

		var addr net.Addr
		var n int
		buf := make([]byte, 1024)
		for {
			_, err = Conn.WriteToUDP([]byte("TimeSync"), ServerAddr)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 1)
			Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			n, addr, err = Conn.ReadFrom(buf)
			if err != nil {
				log.Println(err)
			} else {
				break
			}
		}
		resp := string(buf[0:n])
		log.Println("Received ", resp, " from ", addr)
		ServerUnixTime, err := strconv.ParseInt(resp, 10, 64)
		if err != nil {
			log.Println(err)
		}
		serverDate := time.Unix(ServerUnixTime, 0).Format("2006-01-02")
		serverTime := time.Unix(ServerUnixTime, 0).Format("15:04:05")
		log.Println(serverDate)
		log.Println(serverTime)
		UpdateDate(serverDate)
		UpdateTime(serverTime)
	default:
		log.Fatalln("Please tell me the type.")
		return
	}
}
