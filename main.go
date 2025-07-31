package main

import (
	"fmt"

	config "github.com/GitSiege7/blog_aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read: %v", err))
	}

	c.SetUser("CJ")

	new_c, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read: %v", err))
	}

	fmt.Println(new_c)
}
