package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"

	"github.com/ereminiu/voting/internal/config"
	"github.com/ereminiu/voting/internal/pkg/repository"
	"github.com/ereminiu/voting/internal/pkg/service"
	"github.com/gin-gonic/gin"
)

var mode string

func init() {
	flag.StringVar(&mode, "mode", "test", "config mode")
	flag.Parse()
}

func main() {
	// load configs
	cfg, err := config.LoadConfigs(mode)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := repository.NewDB(cfg)
	if err != nil {
		slog.Error("error during db connection: ", err)
		return
	}

	pollService, err := service.NewPollService(db)
	if err != nil {
		slog.Error("error during poll service creation: ", err)
		return
	}

	router := gin.Default()

	router.POST("/add-hero", func(ctx *gin.Context) {
		var input struct {
			HeroName string `json:"hero_name"`
		}
		err := ctx.BindJSON(&input)
		if err != nil {
			log.Fatalln(err)
			return
		}

		id, err := pollService.CreateHero(input.HeroName)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadGateway, gin.H{
				"message": "user not found",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"id": id,
		})
	})

	router.Run(":8000")

	err = db.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
