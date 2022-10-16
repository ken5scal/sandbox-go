package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"testing"
)

func Test_run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		//テスト対象を別ゴルーチンで実行
		// errGroupを使う理由は、サーバーをgoroutineで立ち上げており
		// Shutdown時にerrがないか確認したいため
		return run(ctx)
	})

	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
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
