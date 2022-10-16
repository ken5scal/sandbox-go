package services_test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ken5scal/api-go-mid-level/services"
	"log"
	"os"
	"testing"
)

var aSer *services.MyAppService

func TestMain(m *testing.M) {
	dbPwd := os.Getenv("ROOTPASS")
	if dbPwd == "" {
		log.Fatal("ENV ROOTPASS is empty")
	}
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", "docker", dbPwd, "sampledb")
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatal(err)
	}

	aSer = services.NewMyAppService(db)
	m.Run()
}

func BenchmarkMyAppService_GetArticleService(b *testing.B) {
	articleID := 1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := aSer.GetArticleService(articleID); err != nil {
			b.Error(err)
			break
		}

	}
}
