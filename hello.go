package main

import (
	"syrenity/server/middleware"
	"syrenity/server/routes"
	"syrenity/server/socket"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {
	connString := "host=localhost port=5432 user=postgres password=postgres dbname=syrenity"

	db := sqlx.MustConnect("postgres", connString)
	defer db.Close()

	websocketServer := socket.NewServer()
	go websocketServer.Run()

	router := gin.Default()

	router.GET("/ws", func(c *gin.Context) {
		socket.Serve(c, websocketServer, db)
	})

	api := router.Group("/api")
	api.Use(middleware.RequireToken(db))
	{
		routes.RegisterUserRoutes(api, db)
		routes.RegisterServerRoutes(api, db)
		routes.RegisterChannelRoutes(api, db, websocketServer)
	}

	router.Use((middleware.ValidateParams(db)))

	router.Run(":3000")
}
