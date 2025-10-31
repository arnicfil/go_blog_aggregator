package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/arnicfil/go_blog_aggregator/internal/config"
	"github.com/arnicfil/go_blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func run() error {
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("Error while reading config: %w\n", err)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		return fmt.Errorf("Error while opening database: %w\n", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	st := state{Cfg: &cfg, DbQ: dbQueries}
	cmds := returnCommands()

	if len(os.Args) < 2 {
		return errors.New("Error need at least 1 argument\n")
	}

	input_command := command{
		Name:      os.Args[1],
		Arguments: os.Args[2:],
	}

	err = cmds.run(&st, input_command)
	if err != nil {
		return fmt.Errorf("Error while calling function: %w\n", err)
	}

	err = cfg.Write()
	if err != nil {
		return fmt.Errorf("Error while writing config: %w", err)
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}
