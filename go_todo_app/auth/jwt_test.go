package auth

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ken5scal/go_todo_app/clock"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/store"
	"github.com/ken5scal/go_todo_app/testutil/fixture"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestEmbed(t *testing.T) {
	want := []byte("BEGIN PUBLIC KEY")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}

	want = []byte("BEGIN PRIVATE KEY")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPrivKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	c := clock.FixedClocker{}
	testUser := fixture.User(&entity.User{ID: entity.UserID(123)})
	mockStorage := &StorageMock{
		SaveFunc: func(ctx context.Context, key string, userID entity.UserID) error {
			if userID != testUser.ID {
				t.Errorf("want %d, but got %d", testUser.ID, userID)
			}
			return nil
		}}
	type fields struct {
		Storage Storage
		Clocker clock.Clocker
	}
	type args struct {
		user *entity.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{Storage: mockStorage, Clocker: c},
			args:    args{testUser},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			sut, err := NewJWTer(tt.fields.Storage, tt.fields.Clocker)
			if err != nil {
				t.Fatal(err)
			}

			got, err := sut.GenerateToken(ctx, *tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("token is empty")
			}
		})
	}
}

type FixedFutureClocker struct{}

func (c FixedFutureClocker) Now() time.Time {
	return clock.FixedClocker{}.Now().Add(100 * time.Hour)
}

func TestJWTer_GetToken(t *testing.T) {
	fixedClocker := clock.FixedClocker{}
	pubKey, err := parse(rawPubKey)
	if err != nil {
		t.Fatal(err)
	}
	prvKey, err := parse(rawPrivKey)
	if err != nil {
		t.Fatal(err)
	}
	testUser := fixture.User(&entity.User{ID: entity.UserID(123)})
	mockStorage := &StorageMock{
		LoadFunc: func(ctx context.Context, key string) (entity.UserID, error) {
			return testUser.ID, nil
		},
	}
	mockStorageWithNoToken := &StorageMock{
		LoadFunc: func(ctx context.Context, key string) (entity.UserID, error) {
			return entity.UserID(0), store.ErrNotFound
		},
	}

	type fields struct {
		PrivateKey jwk.Key
		PublicKey  jwk.Key
		Storage    Storage
		Clocker    clock.Clocker
	}
	type args struct {
		user entity.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    jwt.Token
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				PrivateKey: prvKey,
				PublicKey:  pubKey,
				Storage:    mockStorage,
				Clocker:    fixedClocker,
			},
			args:    args{user: *testUser},
			want:    nil,
			wantErr: false,
		},
		{
			name: "fail(expired token)",
			fields: fields{
				PrivateKey: prvKey,
				PublicKey:  pubKey,
				Storage:    mockStorage,
				Clocker:    FixedFutureClocker{},
			},
			args: args{
				user: *testUser,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail(not in KVS)",
			fields: fields{
				PrivateKey: prvKey,
				PublicKey:  pubKey,
				Storage:    mockStorageWithNoToken,
				Clocker:    fixedClocker,
			},
			args: args{
				user: *testUser,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := &JWTer{
				PrivateKey: tt.fields.PrivateKey,
				PublicKey:  tt.fields.PublicKey,
				Storage:    tt.fields.Storage,
				Clocker:    tt.fields.Clocker,
			}

			wantJwt, _ := jwt.NewBuilder().
				IssuedAt(fixedClocker.Now()).
				Expiration(fixedClocker.Now().Add(expirationInMin)).
				Claim(RoleKey, tt.args.user.Role).
				Claim(UserNameKey, tt.args.user.Name).
				Build()
			b, _ := jwt.Sign(wantJwt, jwt.WithKey(jwa.RS256, sut.PrivateKey))

			r := httptest.NewRequest(http.MethodGet, "https://hogehoge.com", nil)
			r.Header.Set(defaultJWTHeader, fmt.Sprintf(`Bearer %s`, string(b)))
			got, err := sut.GetToken(context.Background(), r)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.want == nil {
				return
			}

			if !reflect.DeepEqual(got, wantJwt) {
				t.Errorf("GetToken() got = %v, want %v", got, wantJwt)
			}
		})
	}
}
