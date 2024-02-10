package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ereminiu/voting/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var mode string
var cmd string
var ver int

func init() {
	flag.StringVar(&mode, "mode", "test", "config mode")
	flag.StringVar(&cmd, "cmd", "up", "up/down/force")
	flag.IntVar(&ver, "v", 0, "force version")
	flag.Parse()
}

func main() {
	// load configs
	cfg, err := config.LoadConfigs(mode)
	if err != nil {
		log.Fatalln(err)
	}
	// databaseURL := "postgres://ys-user:qwerty@localhost:5432/ys-db?sslmode=disable"

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	// TODO: fix it
	mpath := "file://internal/migrate/migrations"
	if envPath := os.Getenv("MIGRATIONS_PATH"); envPath != "" {
		mpath = envPath
	}
	m, err := migrate.New(mpath, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// mode, cmd = "test", "up"
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	// switch cmd {
	// case "up":
	// 	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	// 		log.Fatal(err)
	// 	}
	// 	break
	// case "down":
	// 	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	// 		log.Fatal(err)
	// 	}
	// 	break
	// case "force":
	// 	if err := m.Force(ver); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	// 		log.Fatal(err)
	// 	}
	// 	break
	// }

	// m.Down() - to discard changes
	// m.Force() - to fix dirty version of migrations

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Applied migration: %d, Dirty: %t\n", version, dirty)
}
