package auth

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/ken5scal/go_todo_app/clock"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"time"
)

// ファイルを実行バイナリに埋め込む（多分相対パスで指定できるので実行環境を考慮したファイルパスを指定する必要がない）
// リアルなビルド環境では、ビルド時に鍵を配置する形になる。また秘密鍵はファイル形式でバイナリに埋め込んではいけない（リバエンの危険性）

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

//go:generate go run github.com/matryer/moq -out moq_test.go . Storage
type Storage interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Storage               Storage
	Clocker               clock.Clocker
}

func NewJWTer(s Storage, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Storage: s}
	privKey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("faile parsing key: private key: %w", err)
	}
	pubKey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("faile parsing key: pub key: %w", err)
	}

	j.PrivateKey = privKey
	j.PublicKey = pubKey
	j.Clocker = c
	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	return jwk.ParseKey(rawKey, jwk.WithPEM(true))
}

const (
	RoleKey     = "role"
	UserNameKey = "user_name"
	// lestrrat-go/jwx/v2/jwtデフォ
	defaultJWTHeader = "Authorization"
	expirationInMin  = 30 * time.Minute // in minu
)

func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("github.com/ken5scal/go_todo_app").
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(expirationInMin)).
		Claim(RoleKey, u.Role).
		Claim(UserNameKey, u.Name).
		Build()
	if err != nil {
		return nil, fmt.Errorf("GetToken: failed to build token: %w", err)
	}
	if err := j.Storage.Save(ctx, tok.JwtID(), u.ID); err != nil {
		return nil, err
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return signed, nil
}

func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		r,
		//jwt.ParseRequestのデフォ検索先
		jwt.WithHeaderKey(defaultJWTHeader),
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false))
	if err != nil {
		return nil, err
	}
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token: %w", err)
	}

	if _, err := j.Storage.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired: %w", token.JwtID(), err)
	}
	return token, nil
}
