package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/arnicfil/go_blog_aggregator/internal/config"
	"github.com/arnicfil/go_blog_aggregator/internal/database"
	"github.com/arnicfil/go_blog_aggregator/internal/rss"
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
	cmds.Cmds["agg"] = handlerAggr
	cmds.Cmds["addFeed"] = middlewareLoggedIn(handlerAddFeed)
	cmds.Cmds["feeds"] = handlerListFeeds
	cmds.Cmds["follow"] = middlewareLoggedIn(handlerFollow)
	cmds.Cmds["following"] = middlewareLoggedIn(handlerFollowing)
	cmds.Cmds["unfollow"] = middlewareLoggedIn(handlerUnfollow)
	cmds.Cmds["browse"] = middlewareLoggedIn(handlerBrowse)

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

func handlerAggr(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Command requires 1 argumen \"timeN_between_regs\"")
	}

	time_between_regs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("Error while parsing duration argument: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_regs)

	ticker := time.NewTicker(time_between_regs)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		return errors.New("Command requires 2 arguments")
	}

	nameFeed := cmd.Arguments[0]
	urlFeed := cmd.Arguments[1]

	ctx := context.Background()

	feed, err := s.DbQ.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      nameFeed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Url:       urlFeed,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while creating new feed in database: %w", err)
	}

	_, err = s.DbQ.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	fmt.Print(feed)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("Command doesn't require arguments")
	}
	ctx := context.Background()
	feeds, err := s.DbQ.ListFeeds(ctx)
	if err != nil {
		return fmt.Errorf("Error while requesting feeds from database in handlerListFeeds: %w", err)
	}

	for i, feed := range feeds {
		feedUser, err := s.DbQ.RetrieveFeedUser(ctx, feed.Name)
		if err != nil {
			fmt.Printf("Error while retrieving user of feed %s: %v", feed.Name, err)
			continue
		}
		fmt.Printf("Feed number %d\n", i+1)
		fmt.Println("--------------------------")
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed url: %s\n", feed.Url)
		fmt.Printf("Feed user: %s\n\n", feedUser)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Command requires 1 argument")
	}

	ctx := context.Background()

	url := cmd.Arguments[0]
	feed, err := s.DbQ.FeedFromUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("Error while requesting feed with url in handlerFollow: %w", err)
	}

	_, err = s.DbQ.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while creating feed_follow relation in handlerFollow: %w", err)
	}

	fmt.Printf("User %s successfully follow feed %s\n", user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("Command doesnt require any arguments")
	}
	ctx := context.Background()

	user_feeds, err := s.DbQ.GetFeedFollowsForUser(ctx, user.Name)
	if err != nil {
		return fmt.Errorf("Error while requesting feeds for user %s: %w", s.Cfg.Name, err)
	}

	fmt.Printf("User %s follows feeds:\n", s.Cfg.Name)
	for _, feed := range user_feeds {
		fmt.Printf("  - %s\n", feed.Name_2)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("Command requires any 1 argument \"url\"")
	}

	ctx := context.Background()
	url := cmd.Arguments[0]

	feed, err := s.DbQ.FeedFromUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("Feed with this url does not exist: %w", err)
	}

	err = s.DbQ.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		ID:   user.ID,
		ID_2: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error while deleting feed_follow: %w", err)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	argLen := len(cmd.Arguments)
	if !(argLen == 1 || argLen == 2) {
		return errors.New("Command requires any 1 or 2 arguments feed, limit")
	}

	var limit int32
	if argLen != 2 {
		limit = 2
	} else {
		lim, err := strconv.Atoi(cmd.Arguments[1])
		if err != nil {
			fmt.Printf("Error while parsing limit argument, switching to default (2): %v", err)
			limit = 2
		} else {
			limit = int32(lim)
		}
	}

	ctx := context.Background()
	posts, err := s.DbQ.GetPostsForUser(ctx, database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: limit,
	})
	if err != nil {
		return fmt.Errorf("Error while retrieving posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("**%s**\n", post.Name)
		var pub string
		if post.PublishedAt.Valid {
			pub = post.PublishedAt.Time.String()
		} else {
			pub = "unknown"
		}

		var desc string
		if post.Description.Valid {
			desc = post.Description.String
		} else {
			desc = "none"
		}

		fmt.Printf("published at %s\n", pub)
		fmt.Printf("Description:\n%s\n\n", desc)
	}

	return nil
}
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.DbQ.GetUser(context.Background(), s.Cfg.Name)
		if err != nil {
			return fmt.Errorf("Error in middlewareLoggedIn: %w", err)
		}

		return handler(s, cmd, user)
	}
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

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeed, err := s.DbQ.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching next feed from database: %w", err)
	}

	err = s.DbQ.MarkFeedFetched(ctx,
		database.MarkFeedFetchedParams{
			ID:        nextFeed.ID,
			UpdatedAt: time.Now(),
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})

	feed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error while fetching feed from internet: %w", err)
	}

	for _, item := range feed.Channel.Item {
		/*
			parsedTime, err := time.Parse(,item.PubDate)
			val := true
			if err != nil {
				fmt.Printf("Error while parsing pubDate: %v", err)
				val = false
			}
		*/
		_, err := s.DbQ.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			FeedID: nextFeed.ID,
		})

		if err != nil {
			fmt.Printf("Error while inserting post %s into database: %v", item.Title, err)
		}
	}

	return nil
}
