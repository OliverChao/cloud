package main

import (
	"cloud/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "cloud_" + defaultTableName
	}
}

func main() {
	//testConServer()
	testConDB()
	//testRedis()

}

func testRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "toor",
		DB:       0,
		PoolSize: 10,
	})
	defer client.Close()

	client.FlushAll()

	client.HMSet("user:1:info", map[string]interface{}{
		"name":  "oliver",
		"age":   21,
		"lover": "annabelle",
		"s":     "oliver loves annabelle",
	})

	all := client.HGetAll("user:1:info")
	fmt.Println(all.Result())
}

func testConDB() {
	//"user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
	mysqlSign := "root:toor@(127.0.0.1:3306)/cloud?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", mysqlSign)
	if err != nil {
		//gin.Logger()
		fmt.Println("err....db...")
		return
	}
	defer db.Close()

	db.DropTableIfExists(&model.Model{})

	db.AutoMigrate(&model.Model{})
	if err = db.AutoMigrate(&model.Model{}).Error; err != nil {
		fmt.Println("auto migrate tables failed:", err)
		return
	} else {
		fmt.Println("auto migrate tables successfully")
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(55)
	db.DB().SetConnMaxLifetime(5 * time.Minute)

}

func testConServer() {
	ret := gin.New()
	ret.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		//c.String(200,"hahaha")
	})

	//ret.Run()
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: ret,
	}
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		go func() {
			s := <-c
			fmt.Printf("got signal [%s], exiting pipe now", s)
			if err := server.Close(); nil != err {
			}
			os.Exit(0)
		}()
	}()

	if err := server.ListenAndServe(); nil != err {
		fmt.Printf("listen and serve failed: " + err.Error())
	}
}
