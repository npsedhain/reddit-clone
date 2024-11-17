package main

import (
	"fmt"
	"reddit/actors"
	"reddit/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

func main() {
    // Create the actor system
    system := actor.NewActorSystem()

    // Create the user actor
    userProps := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewUserActor()
    })
    userPID := system.Root.Spawn(userProps)

    // Create the subreddit actor
    subredditProps := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewSubredditActor()
    })
    subredditPID := system.Root.Spawn(subredditProps)

    // Register users
    users := []string{"user1", "user2"}
    for _, username := range users {
        registerMsg := &messages.RegisterUser{
            Username: username,
            Password: "password123",
        }
        result, _ := system.Root.RequestFuture(userPID, registerMsg, 5*time.Second).Result()
        if response, ok := result.(*messages.RegisterUserResponse); ok && response.Success {
            fmt.Printf("User registered: %s\n", username)
        }
    }

    // Create a subreddit
    createSubMsg := &messages.CreateSubreddit{
        Name:        "programming",
        Description: "Programming discussions",
        CreatorId:   "user1",
    }
    result, _ := system.Root.RequestFuture(subredditPID, createSubMsg, 5*time.Second).Result()
    if response, ok := result.(*messages.CreateSubredditResponse); ok && response.Success {
        fmt.Printf("Subreddit created: %s\n", createSubMsg.Name)
    }

    // Join subreddit
    joinMsg := &messages.JoinSubreddit{
        UserId:        "user2",
        SubredditName: "programming",
    }
    result, _ = system.Root.RequestFuture(subredditPID, joinMsg, 5*time.Second).Result()
    if response, ok := result.(*messages.JoinSubredditResponse); ok && response.Success {
        fmt.Printf("User %s joined subreddit: %s\n", joinMsg.UserId, joinMsg.SubredditName)
    }

    // Leave subreddit
    leaveMsg := &messages.LeaveSubreddit{
        UserId:        "user2",
        SubredditName: "programming",
    }
    result, _ = system.Root.RequestFuture(subredditPID, leaveMsg, 5*time.Second).Result()
    if response, ok := result.(*messages.LeaveSubredditResponse); ok && response.Success {
        fmt.Printf("User %s left subreddit: %s\n", leaveMsg.UserId, leaveMsg.SubredditName)
    }

    // Get members to verify
    getMembersMsg := &messages.GetSubredditMembers{
        SubredditName: "programming",
    }
    result, _ = system.Root.RequestFuture(subredditPID, getMembersMsg, 5*time.Second).Result()
    if response, ok := result.(*messages.GetSubredditMembersResponse); ok && response.Success {
        fmt.Printf("Current members of %s: %v\n", getMembersMsg.SubredditName, response.Members)
    }

    // Keep the program running
    time.Sleep(1 * time.Second)
}
