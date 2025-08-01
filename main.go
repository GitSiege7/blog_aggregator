package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"

	config "github.com/GitSiege7/blog_aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read: %v", err))
		return
	}

	s := state{
		&c,
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Printf("error: insufficient args (%v)\n", len(args))
		os.Exit(1)
	}

	err = cmds.run(&s, command{args[0], args[1:]})
	if err != nil {
		fmt.Println(fmt.Errorf("failed to run: %v", err))
		return
	}
}
