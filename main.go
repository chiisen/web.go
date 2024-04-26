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

	"github.com/go-redis/redis/v8"
)

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	router.POST("/api/Member/login", func(c *gin.Context) {
		c.String(http.StatusOK, "login")
	})

	router.GET("/healthcheck", func(c *gin.Context) {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "redis-cluster.h1-redis-dev:6379",
			Password: "h1devredis1688", // no password set
			DB:       0,                // use default DB
		})

		ctx := context.Background()

		val, err := rdb.Ping(ctx).Result()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(val)
		}

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

		var name, state_desc string
		for rows.Next() {
			err := rows.Scan(&name, &state_desc)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("name: %s, state_desc: %s\n", name, state_desc)
		}
		str := fmt.Sprintf("val: %v, name: %s, state_desc: %s", val, name, state_desc)
		c.String(http.StatusOK, str)
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
