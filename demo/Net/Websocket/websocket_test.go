// go_bind_cross_test.go
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

// 获取本机第一个非环回 IPv4
func localIPv4() (string, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifs {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					return ip4.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("找不到非环回 IPv4")
}

// 启动一个简单 HTTP 服务，绑定 addr，返回 text
func startService(addr, text string, wg *sync.WaitGroup) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, text)
	})
	server := &http.Server{Handler: mux}

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve(ln)
	}()
	return ln, nil
}

// 测试连接四种组合
func testConn(name, url string) {
	start := time.Now()
	resp, err := http.Get(url)
	dur := time.Since(start)
	if err != nil {
		log.Printf("%s → FAIL: %v (%.3fs)\n", name, err, dur.Seconds())
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Printf("%s → OK: %q (%.3fs)\n", name, string(body), dur.Seconds())
}

func Test(t *testing.T) {
	ipv4, err := localIPv4()
	if err != nil {
		log.Fatal(err)
	}
	port := "8080"
	addrLoop := "127.0.0.1:" + port
	addrExt := ipv4 + ":" + port

	var wg sync.WaitGroup
	lnA, err := startService(addrLoop, "Service A (loopback)", &wg)
	if err != nil {
		log.Fatalf("启动 Service A 失败: %v\n", err)
	}
	lnB, err := startService(addrExt, "Service B (external)", &wg)
	if err != nil {
		log.Fatalf("启动 Service B 失败: %v\n", err)
	}

	// 等待服务就绪
	time.Sleep(200 * time.Millisecond)

	log.Println("---- 测试开始 ----")
	testConn("A from 127.0.0.1", "http://127.0.0.1:"+port+"/")
	testConn("A from "+ipv4, "http://"+ipv4+":"+port+"/")
	testConn("B from 127.0.0.1", "http://127.0.0.1:"+port+"/")
	testConn("B from "+ipv4, "http://"+ipv4+":"+port+"/")

	lnA.Close()
	lnB.Close()
	wg.Wait()
	log.Println("---- 测试结束 ----")
}
