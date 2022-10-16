package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net"
	"net/http"
	"testing"
)

func Test_run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		//テスト対象を別ゴルーチンで実行
		// errGroupを使う理由は、サーバーをgoroutineで立ち上げており
		// Shutdown時にerrがないか確認したいため
		return run(ctx, l)
	})

	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get : %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := fmt.Sprintf("hello, %s", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	cancel()
	if err := eg.Wait(); err != nil {
		fmt.Println(err.Error())
		t.Fatal(err)
	}
}
