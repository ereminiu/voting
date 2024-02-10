package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/ereminiu/voting/internal/config"
	"github.com/ereminiu/voting/internal/events"
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
	// set up logger
	log := setupLogger()

	// load configs
	cfg, err := config.LoadConfigs(mode)
	if err != nil {
		log.Error(
			"error occured",
			err,
		)
		return
	}
	log.Info(
		"server is started",
		slog.String("port", cfg.Port),
	)

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

	router.POST("/create-poll", func(ctx *gin.Context) {
		var input events.PollEvent
		err := ctx.BindJSON(&input)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "wrong input format",
			})
			return
		}
		pollId, choiceIds, err := pollService.CreatePoll(input)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "something went wrong",
				"error":   err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message":    "poll is created",
			"poll_id":    pollId,
			"choice_ids": choiceIds,
		})
	})

	router.POST("/add-hero", func(ctx *gin.Context) {
		var input struct {
			HeroName string `json:"hero_name"`
		}
		err := ctx.BindJSON(&input)
		if err != nil {
			log.Error(
				"error",
				err,
			)
			return
		}

		id, err := pollService.CreateHero(input.HeroName)
		if err != nil {
			log.Error(
				"error",
				err,
			)
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
		log.Error(
			"error",
			err,
		)
	}
}

func setupLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
}
