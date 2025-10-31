package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"

	"github.com/arnicfil/go_blog_aggregator/internal/config"
	"github.com/arnicfil/go_blog_aggregator/internal/database"
)

type commands struct {
	Cmds map[string]func(*state, command) error
}

func returnCommands() commands {
	cmds := commands{Cmds: make(map[string]func(*state, command) error)}

	cmds.Cmds["login"] = handlerLogin
	cmds.Cmds["register"] = handlerRegister
	cmds.Cmds["reset"] = handlerReset
	cmds.Cmds["users"] = handlerListUsers

	return cmds
}

type command struct {
	Name      string
	Arguments []string
}

type state struct {
	Cfg *config.Config
	DbQ *database.Queries
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Command requires 1 argument")
	}
	ctx := context.Background()

	_, err := s.DbQ.GetUser(ctx, cmd.Arguments[0])
	if err != nil {
		return errors.New("User doesn't exists")
	}

	s.Cfg.Name = cmd.Arguments[0]
	err = s.Cfg.Write()
	if err != nil {
		return fmt.Errorf("Erorr setting user name in handlerLogin: %w", err)
	}

	fmt.Println("User name has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Command requires 1 argument")
	}

	ctx := context.Background()

	_, err := s.DbQ.GetUser(ctx, cmd.Arguments[0])
	if err == nil {
		return errors.New("User already exists")
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
	}
	user, err := s.DbQ.CreateUser(ctx, newUser)
	s.Cfg.Name = user.Name

	fmt.Printf("User %s was created\n", user.Name)
	fmt.Printf("CreatedAt %s\n", user.CreatedAt)
	fmt.Printf("UpdatedAt %s\n", user.UpdatedAt)
	fmt.Printf("Id%s\n", user.ID)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("Command doesn't require arguments")
	}
	ctx := context.Background()

	err := s.DbQ.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("Error while deleting entries: %w", err)
	}

	return nil
}

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("Command doesn't require arguments")
	}

	ctx := context.Background()
	currentUser, err := s.DbQ.GetUser(ctx, s.Cfg.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving current user in handlerListUsers: %w", err)
	}

	users, err := s.DbQ.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("Error retrieving all users in handlerListUsers: %w", err)
	}
	for _, user := range users {
		stringToPrint := " * %s"
		if user == currentUser.Name {
			stringToPrint = stringToPrint + " (current)"
		}
		fmt.Printf(stringToPrint+"\n", user)
	}

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("Error command %s does not exist in commands", cmd.Name)
	}

	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("Erorr executing command %s: %w", cmd.Name, err)
	}
	return nil
}
