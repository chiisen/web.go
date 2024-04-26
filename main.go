package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		"daydb-svc.h1-db-dev", "mobile_api", "a:oY%~^E+VU0", 1433, "HKNetGame_HJ")

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, state_desc FROM sys.databases WHERE name = 'HKNetGame_HJ';")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var name, state_desc string // 這裡的變數類型和數量應該與你的資料表欄位匹配
		err := rows.Scan(&name, &state_desc)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("name: %s, state_desc: %s\n", name, state_desc)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	router.POST("/api/Member/login", func(c *gin.Context) {
		c.String(http.StatusOK, "login")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// 服務連線
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中斷信號以優雅地關閉服務器（設置 5 秒的超時時間）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
