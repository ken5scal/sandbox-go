package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		// 引数で受け取ったnet.Listenerを利用するので、 Addrフィールドは指定しない
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// errが含まれるゴルーチン間の並行処理の実装が簡単
	// 別ゴルーチン上で実行する関数からerrorを戻り値として受け取れる（sync.WaitGroupはできない）
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// 別ゴルーチン。func() errorというシグネチャの関数を起動できる
		if err := s.srv.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// 別ゴルーチンの終了を待って、その中のerrorを返す
	return eg.Wait()
}
