package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid,
		CreatedAt: time,
		UpdatedAt: time,
		Name:      cmd.args[0]})

	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("User Created:\nUUID: %v\nCreated: %v\nUpdated: %v\nName: %v\n", uuid, time, time, cmd.args[0])

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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: agg {interval} ('30s', '10m')")
	}

	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed parseduration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Println("usage: addfeed {name} {url}")
		os.Exit(1)
	}

	feed_UUID := uuid.New()
	time := time.Now()

	err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        feed_UUID,
		CreatedAt: time,
		UpdatedAt: time,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID})
	if err != nil {
		return err
	}

	follow_UUID := uuid.New()

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        follow_UUID,
		CreatedAt: time,
		UpdatedAt: time,
		UserID:    user.ID,
		FeedID:    feed_UUID,
	})

	fmt.Printf("Feed Created:\nUUID: %v\nCreated: %v\nUpdated: %v\nName: %v\nURL: %v\nUser: %v\n", feed_UUID, time, time, cmd.args[0], cmd.args[1], user.Name)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed getfeeds: %v", err)
	}

	for _, feed := range feeds {
		username, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed getuserbyid: %v", err)
		}

		fmt.Printf("Name: %v\nURL: %v\nUser: %v\n", feed.Name, feed.Url, username)
		fmt.Printf("\n")
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Printf("usage: follow {url}")
		os.Exit(1)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed getfeedbyurl: %v", err)
	}

	uuid := uuid.New()
	time := time.Now()

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid,
		CreatedAt: time,
		UpdatedAt: time,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed createfeedfollow: %v", err)
	}

	fmt.Println("Created Follow:")
	fmt.Printf("Name: %v\n", feed_follow.FeedName)
	fmt.Printf("Current User: %v\n", feed_follow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed getfollows: %v", err)
	}

	fmt.Printf("%v's feeds:\n", s.c.Current_user_name)
	for _, follow := range follows {
		fmt.Printf(" - %v\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: unfollow {url}")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed getfeedbyurl: %v", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed deletefeedfollow: %v", err)
	}

	fmt.Println("Successfully unfollowed!")

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.c.Current_user_name)
		if err != nil {
			return fmt.Errorf("failed getuser: %v", err)
		}

		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed getnextfeedtofetch: %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed markfeedfetched: %v", err)
	}

	fmt.Println("Fetching feed...")

	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed fetchfeed: %v", err)
	}

	Decode_escaped(rss)

	for _, item := range rss.Channel.Item {
		now := time.Now()
		UUID := uuid.New()

		date, err := time.Parse("Mon, 02 Jan 2006 15:04:05 Z0700", item.PubDate)
		if err != nil {
			return fmt.Errorf("failed timeparse on %v: %v", item.PubDate, err)
		}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        UUID,
			CreatedAt: now,
			UpdatedAt: now,
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Url: rss.Channel.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: date,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}

			fmt.Println(err)
		}
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) != 0 {
		parsed, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("failed parseint: %v", err)
		}

		limit = int32(parsed)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("failed getpostsforuser: %v", err)
	}

	for _, post := range posts {
		fmt.Println(post.Title.String)
		fmt.Println(post.Description.String)
		fmt.Println(post.PublishedAt)
		fmt.Println()
	}

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	return c.callback[cmd.name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.callback[name]; ok {
		return
	}

	c.callback[name] = f
}
