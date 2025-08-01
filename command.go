package main

import (
	"context"
	"fmt"
	"os"
	"time"

	config "github.com/GitSiege7/blog_aggregator/internal/config"
	"github.com/GitSiege7/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db *database.Queries
	c  *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	callback map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no arguments")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		fmt.Printf("no user found: %v\n", cmd.args[0])
		os.Exit(1)
	}

	err := s.c.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Println("User logged in!")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no arguments")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		fmt.Println("User already exists")
		os.Exit(1)
	}

	uuid := uuid.New()
	time := time.Now()

	s.db.CreateUser(context.Background(), database.CreateUserParams{uuid, time, time, cmd.args[0]})

	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("User created:\nUUID: %v\nCreated: %v\nUpdated: %v\nName: %v\n", uuid, time, time, cmd.args[0])

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Resets(context.Background()); err != nil {
		return fmt.Errorf("failed to reset: %v", err)
	}

	fmt.Println("Reset successful")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	names, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	for _, name := range names {
		fmt.Printf("* %v", name)
		if name == s.c.Current_user_name {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	err := c.callback[cmd.name](s, cmd)
	if err != nil {
		return fmt.Errorf("failed to run: %v", err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.callback[name]; ok {
		return
	}

	c.callback[name] = f
}
