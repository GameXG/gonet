package main

import (
	"net"
	"time"

	"fmt"

	"github.com/gamexg/gonet/tuntap"
)

// 尝试手工实现 tcp 协议，这样测试一下性能。

func main() {
	f := tuntap.NewTunTapF()
	lAddr := net.IPv4(192, 168, 111, 1)
	network := net.IPv4(192, 168, 111, 0)
	mask := net.IPv4(255, 255, 255, 0)

	tap, err := f.NewTun(lAddr, network, mask)
	if err != nil {
		panic(err)
	}
	defer tap.Close()

	go tap.Server()

	go func() {
		for buf := range tap.GetReadChan() {
			fmt.Println(buf)
		}
	}()

	func() {

		buf := make([]byte, 1024)
		sTime := time.Now()
		for i := 0; i < 10000; i++ {
			err := tap.WritePack(buf)
			if err != nil {
				panic(err)
			}
		}
		eTime := time.Now()
		fmt.Println("同步写入10000包耗时：", eTime.Sub(sTime))
		fmt.Println("平均耗时： ", eTime.Sub(sTime)/10000)
	}()

	time.Sleep(1 * time.Hour)
}

/*

tcp 协议要求使用4个元素确定一个连接，虽然一些系统一个端口只能建立一个对外连接连接，但是这里还是按照 tcp 标准来实现吧。

计划只模拟远端设备，并不管本地的ip等。

计划结构是这样的：

远端池->远端ip设备->远端端口(连接池)->单条连接

这里狗屎依靠tcp协议的4个元素进行的分路。

最后一步需要注意下，计划使用窗口大小而不是丢包来限制客户端的发包速度。


完成端口测试结果
同步写入10000包耗时： 13.7945879s
平均耗时：  1.379458ms

*/
