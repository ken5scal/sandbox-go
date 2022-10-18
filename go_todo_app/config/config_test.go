package config

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			wantErr: false,
			name:    "デフォルト値",
			want: &Config{
				Env:        "dev",
				Port:       80,
				DBHost:     "127.0.0.1",
				DBPort:     33306,
				DBUser:     "todo",
				DBPassword: "todo",
				DBName:     "todo",
				RedisHost:  "127.0.0.1",
				RedisPort:  36379,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
