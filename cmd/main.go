package main

import (
	"log"
	"time"

	"github.com/Mark-Grigorev/auth/internal/config"
	"github.com/Mark-Grigorev/auth/internal/controller"
	"github.com/Mark-Grigorev/auth/internal/db"
	jwtmanager "github.com/Mark-Grigorev/auth/internal/jwt_manager"
	"github.com/Mark-Grigorev/auth/internal/logic"
)

func main() {
	cfg := config.Read()

	db, err := db.New(cfg.DBConfig.DBConnection)
	if err != nil {
		log.Fatalf("error init db - %s", err.Error())
	}
	jwt := jwtmanager.New(cfg.JWTConfig.SecretKey, time.Duration(cfg.JWTConfig.TokenDuration)*time.Second)

	controller.New(cfg.AppConfig, logic.New(cfg, db, jwt)).Start()
}
