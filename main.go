package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	snapshotsapi "github.com/containerd/containerd/api/services/snapshots/v1"
	"github.com/containerd/containerd/v2/contrib/snapshotservice"
	"github.com/containerd/containerd/v2/plugins/snapshots/native"
)

func main() {
	// 调用参数为：<unix-sock-path> <snapshot path>
	if len(os.Args) < 3 {
		fmt.Printf("invalid args: usage: %s <unix addr> <root>\n", os.Args[0])
		os.Exit(1)
	}

	// 创建 grpc 服务
	rpc := grpc.NewServer()

	// 直接初始化 native 用于演示
	sn, err := native.NewSnapshotter(os.Args[2])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// FromSnapshotter 函数是 contrib/snapshotservice/service.go 提供的包装函数。
	// 该函数 snapshots.Snapshotter 接口的方法转换成对应的 rpc 调用
	service := snapshotservice.FromSnapshotter(sn)

	// 注册 grpc 服务
	snapshotsapi.RegisterSnapshotsServer(rpc, service)

	// Listen and serve
	l, err := net.Listen("unix", os.Args[1])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	if err := rpc.Serve(l); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
