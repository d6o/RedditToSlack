package main

import (
	"flag"
	"fmt"
	"github.com/disiqueira/RedditToSlack/pkg/slack"
	"github.com/disiqueira/RedditToSlack/pkg/slack/rtm"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
	"log"
	"os"
)

const slackText = "/r/%s: <%s> %s (http://redd.it/%s) [%s - %s]"

type slackBot struct {
	bot     reddit.Bot
	slack   *slack.Agent
	channel string
	user    string
}

func (r *slackBot) Post(p *reddit.Post) error {
	text := fmt.Sprintf(slackText, p.Subreddit, p.Author, p.Title, p.ID, p.Domain, p.URL)

	m := slack.Message{
		Type:    "message",
		Channel: r.channel,
		User:    r.user,
		Text:    text,
	}

	fmt.Println(text)
	return r.slack.SendMessage(m)
}

func main() {
	subreddit := flag.String("subreddit", "", "the subreddit you want to watch")
	slackToken := flag.String("slackToken", "", "token to post on Slack")
	slackUser := flag.String("slackUser", "", "user to post on Slack")
	slackChannel := flag.String("slackChannel", "", "channel to post on Slack")
	flag.Parse()
	if *subreddit == "" {
		log.Println("You need to inform a subreddit.")
		os.Exit(1)
	}
	if *slackUser == "" {
		log.Println("You need to inform a Slack User.")
		os.Exit(1)
	}
	if *slackChannel == "" {
		log.Println("You need to inform a Slack Channel.")
		os.Exit(1)
	}
	if *slackToken == "" {
		log.Println("You need to inform a Slack Token.")
		os.Exit(1)
	}

	agentPath := fmt.Sprintf("./%s.agent", *subreddit)
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		log.Printf("File not found (%s) \n", agentPath)
		os.Exit(1)
	}

	bot, err := reddit.NewBotFromAgentFile(agentPath, 0)
	if err != nil {
		panic(err)
	}

	slackAgent, err := startSlack(*slackToken)
	if err != nil {
		panic(err)
	}

	cfg := graw.Config{
		Subreddits: []string{
			*subreddit,
		},
	}
	handler := &slackBot{
		bot:     bot,
		slack:   slackAgent,
		channel: *slackChannel,
		user:    *slackUser,
	}
	_, wait, err := graw.Run(handler, bot, cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("graw run failed: ", wait())
}

func startSlack(token string) (slackAgent *slack.Agent, err error) {
	realTime, err := rtm.New(token)
	if err != nil {
		return nil, err
	}
	return slack.New(realTime)
}
