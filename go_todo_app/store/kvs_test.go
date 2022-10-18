package store

import (
	"context"
	"errors"
	"github.com/ken5scal/go_todo_app/config"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/testutil"
	"reflect"
	"testing"
	"time"
)

func TestKVS_Load(t *testing.T) {
	t.Parallel()
	c := testutil.OpenRedisForTest(t)
	sut := &KVS{c}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		key := "TestKVS_Load_ok"
		uid := entity.UserID(123)
		ctx := context.Background()
		c.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() { c.Del(ctx, key) })

		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("want no error, but got %v", err)
		}
		if got != uid {
			t.Errorf("want %d, but got %v", uid, got)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		t.Parallel()
		key := "TestKVS_Save_notFound"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v(value = %d)", ErrNotFound, err, got)
		}
	})
}

func TestKVS_Save(t *testing.T) {
	t.Parallel()
	c := testutil.OpenRedisForTest(t)
	sut := &KVS{c}
	key := "TestKVS_Save"
	uid := entity.UserID(123)
	ctx := context.Background()
	t.Cleanup(func() { c.Del(ctx, key) })
	if err := sut.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func TestNewKVS(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *KVS
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKVS(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKVS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKVS() got = %v, want %v", got, tt.want)
			}
		})
	}
}
