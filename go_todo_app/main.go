package main

import (
	"context"
	"fmt"
	"github.com/ken5scal/go_todo_app/config"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
)

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するので、
		// Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello, %s", (r.URL.Path[1:]))
		}),
	}

	// errが含まれるゴルーチン間の並行処理の実装が簡単
	// 別ゴルーチン上で実行する関数からerrorを戻り値として受け取れる（sync.WaitGroupはできない）
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// 別ゴルーチン。func() errorというシグネチャの関数を起動できる
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// 別ゴルーチンの終了を待って、その中のerrorを返す
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}
