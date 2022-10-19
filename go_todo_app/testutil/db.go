package testutil

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// OpenDBForTest はポート番号（≒実行環境）に応じて接続所情報を切り替える
// ポート番号はほぼほぼ固定化されているため、環境変数から読み取らずともおｋとしている
func OpenDBForTest(t *testing.T) *sqlx.DB {
	port := 33306 // localの場合
	if _, defined := os.LookupEnv("CI"); defined {
		port = 3306 // CI (Github Action)
	}
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("todo:todo@tcp(127.0.0.1:%d)/todo?parseTime=true", port),
	)
	if err != nil {
		t.Fatalf("failed to connect db: %s", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return sqlx.NewDb(db, "mysql")
}
