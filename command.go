package main

import (
	"fmt"

	config "github.com/GitSiege7/blog_aggregator/internal/config"
)

type state struct {
	c *config.Config
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

	err := s.c.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to set user at login: %v", err)
	}

	fmt.Println("User logged in!")

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
