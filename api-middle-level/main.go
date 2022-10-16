package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ken5scal/api-go-mid-level/api"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	dbUser     = "docker" //os.Getenv("ROOTUSER")
	dbDatabase = "sampledb"
)

func main() {
	// goroutineによる並行処理を安全に処理するためのmux, waitgroup
	go cutIngredient()
	boilWater()

	// goroutineによる並行処理を安全に処理するためのchannel
	ch1, ch2 := make(chan int), make(chan string)
	defer close(ch1)
	defer close(ch2)

	go doubleInt(1, ch1)
	go doubleString("hoge", ch2)

	for i := 0; i < 2; i++ {
		select {
		case result := <-ch1:
			fmt.Println(result)
		case result := <-ch2:
			fmt.Println(result)
		}
	}

	dbPwd := os.Getenv("ROOTPASS")
	if dbPwd == "" {
		log.Fatal("ENV ROOTPASS is empty")
	}
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPwd, dbDatabase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatalf("fail to connect DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := api.NewRouter(db)

	log.Println("server start as port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func doubleString(src string, strCh chan<- string) {
	strCh <- strings.Repeat(src, 2)
}

func doubleInt(src int, intCh chan<- int) {
	result := src * 2
	//time.Sleep(5 * time.Second)
	intCh <- result
}

func cutIngredient() {
	fmt.Println("start cutIngredient")
	time.Sleep(1 * time.Second)
	fmt.Println("finish cutIngredient")
}

func boilWater() {
	fmt.Println("start boilWater")
	time.Sleep(2 * time.Second)
	fmt.Println("finish boilWater")
}
