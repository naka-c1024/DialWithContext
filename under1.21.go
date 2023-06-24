package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	// "host:port "という形式のネットワークアドレスに結合
	addr := net.JoinHostPort("localhost", port)

	ctx := context.Background()
	// キャンセルするためのコンテキストを作成
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, cancelFunc, err := DialWithContext("tcp", addr, ctx) // 実装処理
	if err != nil {
		defer cancel()
		return err
	}

	// n秒後にキャンセルする
	go func() {
		<-time.After(5 * time.Second)
		fmt.Printf("\ncancelled\n")
		cancelFunc() // DialWithContextのgoroutineを終了させる
		conn.Close()
	}()

	go func() {
		io.Copy(os.Stdout, conn)
	}()

	// EOFに達するか、エラーが発生するまで動き続ける、コネクションが切れているかどうかは送ってみないとわからないので、最後にEnterは必要
	io.Copy(conn, os.Stdin)

	return nil
}

func DialWithContext(network, address string, ctx context.Context) (net.Conn, func(), error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, func() {}, err
	}

	parent := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-parent:
			return
		}
	}()

	return conn, func() {
		close(parent)
	}, nil
}
