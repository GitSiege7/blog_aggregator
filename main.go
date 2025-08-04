package main

import (
	"database/sql"
	"fmt"
	"os"

	database "github.com/GitSiege7/blog_aggregator/internal/database"

	_ "github.com/lib/pq"

	config "github.com/GitSiege7/blog_aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read: %v", err))
		return
	}

	db, err := sql.Open("postgres", c.Db_url)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to open db: %v", err))
	}

	dbQueries := database.New(db)

	s := state{
		dbQueries,
		&c,
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Printf("error: insufficient args (%v)\n", len(args))
		os.Exit(1)
	}

	err = cmds.run(&s, command{args[0], args[1:]})
	if err != nil {
		fmt.Println(fmt.Errorf("failed to run: %v", err))
		return
	}
}
