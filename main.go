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

	addr := net.JoinHostPort("localhost", port)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var conn net.Conn
	conn, err := DialWithContext("tcp", addr, ctx)
	if err != nil {
		defer cancel()
		return err
	}

	go func() {
		<-time.After(1 * time.Second)
		fmt.Println("cancelled")
		cancel()
	}()

	go func() {
		io.Copy(os.Stdout, conn)
	}()

	io.Copy(conn, os.Stdin)

	return nil
}

func DialWithContext(network, address string, ctx context.Context) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	// ctxで親がキャンセルを呼ばれたら発火する
	context.AfterFunc(ctx, func() {
		fmt.Println("closing")
		conn.Close()
		// conn.SetDeadline(time.Now()) // こっちはコネクションを閉じないけどI/Oが止まる
	})
	return conn, nil
}
