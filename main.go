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

    // Create actors
    userProps := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewUserActor()
    })
    userPID := system.Root.Spawn(userProps)

    subredditProps := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewSubredditActor()
    })
    subredditPID := system.Root.Spawn(subredditProps)

    postProps := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewPostActor()
    })
    postPID := system.Root.Spawn(postProps)

    // Create comment actor
    commentProps := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewCommentActor()
    })
    commentPID := system.Root.Spawn(commentProps)

    dmProps := actor.PropsFromProducer(func() actor.Actor {
        return actors.NewDirectMessageActor()
    })
    dmPID := system.Root.Spawn(dmProps)

    // Register users in parallel
    users := []string{"user1", "user2"}
    var userFutures []*actor.Future
    for _, username := range users {
        registerMsg := &messages.RegisterUser{
            Username: username,
            Password: "password123",
        }
        future := system.Root.RequestFuture(userPID, registerMsg, 5*time.Second)
        userFutures = append(userFutures, future)
    }

    // Wait for all user registrations
    for i, future := range userFutures {
        result, _ := future.Result()
        if response, ok := result.(*messages.RegisterUserResponse); ok && response.Success {
            fmt.Printf("User registered: %s\n", users[i])
        }
    }

    // Create subreddits in parallel
    createSubMsg1 := &messages.CreateSubreddit{
        Name:        "programming",
        Description: "Programming discussions",
        CreatorId:   "user1",
    }
    createSubMsg2 := &messages.CreateSubreddit{
        Name:        "storytelling",
        Description: "Amazing stories",
        CreatorId:   "user2",
    }
    subredditMsgs := []*messages.CreateSubreddit{createSubMsg1, createSubMsg2}
    var subredditFutures []*actor.Future
    for _, msg := range subredditMsgs {
        future := system.Root.RequestFuture(subredditPID, msg, 5*time.Second)
        subredditFutures = append(subredditFutures, future)
    }

    // Wait for all subreddit creations
    for i, future := range subredditFutures {
        result, _ := future.Result()
        if response, ok := result.(*messages.CreateSubredditResponse); ok && response.Success {
            fmt.Printf("Subreddit created: %s\n", subredditMsgs[i].Name)
        }
    }

    // Create posts in parallel
    posts := []messages.Post{
        {
            Title:         "Post",
            Content:      "Hello, this is my first post!",
            AuthorId:     "user1",
            SubredditName: createSubMsg1.Name,
        },
        {
            Title:         "Another Post",
            Content:      "Learning Go programming",
            AuthorId:     "user2",
            SubredditName: createSubMsg2.Name,
        },
    }

    var postFutures []*actor.Future
    for _, post := range posts {
        createPostMsg := &messages.Post{
            Title:         post.Title,
            Content:      post.Content,
            AuthorId:     post.AuthorId,
            SubredditName: post.SubredditName,
        }
        future := system.Root.RequestFuture(postPID, createPostMsg, 5*time.Second)
        postFutures = append(postFutures, future)
    }

    var postIds []string

    // Wait for all post creations
    for i, future := range postFutures {
        result, _ := future.Result()
        if response, ok := result.(*messages.CreatePostResponse); ok && response.Success {
            fmt.Printf("Post created: %s by %s\n", posts[i].Title, posts[i].AuthorId)
            postIds = append(postIds, response.PostId)
        }
    }

    if len(postIds) > 0 {
        comments := []messages.CreateComment{
            {
                Content:  "Great first post!",
                AuthorID: "user2",
                PostID:   postIds[0],
                ParentID: "", // Empty for top-level comments
            },
            {
                Content:  "Thanks for sharing this.",
                AuthorID: "user1",
                PostID:   postIds[0],
                ParentID: "", // Empty for top-level comments
            },
        }

        var commentFutures []*actor.Future
        for _, comment := range comments {
            future := system.Root.RequestFuture(commentPID, &comment, 5*time.Second)
            commentFutures = append(commentFutures, future)
        }

        // Wait for comment creation responses
        for i, future := range commentFutures {
            result, err := future.Result()
            if err != nil {
                fmt.Printf("Error creating comment: %v", err)
                continue
            }
            if response, ok := result.(*messages.CreateCommentResponse); ok {
                fmt.Printf("Comment created by %s, Success: %v, ID: %s",
                    comments[i].AuthorID, response.Success, response.CommentID)
            }
        }
    }

    // Send some direct messages
    var sentMsgID string
    sendDM := &messages.SendDirectMessage{
        FromUserID: "user1",
        ToUserID:   "user2",
        Content:    "Hey, how are you?",
    }
    dmFuture := system.Root.RequestFuture(dmPID, sendDM, 5*time.Second)
    if result, err := dmFuture.Result(); err == nil {
        if response, ok := result.(*messages.SendDirectMessageResponse); ok && response.Success {
            fmt.Printf("Direct message sent successfully, ID: %s\n", response.MessageID)
            sentMsgID = response.MessageID
        }
    }

    // reply to the direct message
    replyDM := &messages.SendDirectMessage{
        ParentID:   sentMsgID,
        Content:    "I'm good, thanks!",
        FromUserID: "user2",
        ToUserID:   "user1",
    }
    replyFuture := system.Root.RequestFuture(dmPID, replyDM, 5*time.Second)
    if result, err := replyFuture.Result(); err == nil {
        if response, ok := result.(*messages.SendDirectMessageResponse); ok && response.Success {
            fmt.Printf("Direct message replied successfully, ID: %s\n", response.MessageID)
        }
    }

    // Get messages for user2
    getMessages := &messages.GetUserMessages{
        UserID: "user2",
    }

    getMsgsFuture := system.Root.RequestFuture(dmPID, getMessages, 5*time.Second)
    if result, err := getMsgsFuture.Result(); err == nil {
        if response, ok := result.(*messages.GetUserMessagesResponse); ok && response.Success {
            fmt.Println("\nMessages for user2:")
            for _, dm := range response.Messages {
                fmt.Printf("From: %s\nContent: %s\n\n", dm.FromUserID, dm.Content)
            }
        }
    }

    // Get all posts from subreddit
    getPostsMsg := []messages.GetSubredditPosts{
        {
            SubredditName: "programming",
        },
        {
            SubredditName: "storytelling",
        },
    }

    var subPostFutures []*actor.Future
    for _, msg := range getPostsMsg {
        future := system.Root.RequestFuture(postPID, &msg, 5*time.Second)
        subPostFutures = append(subPostFutures, future)
    }

    // Wait for all post in the subreddits
    for i, future := range subPostFutures {
        result, _ := future.Result()
        if response, ok := result.(*messages.GetSubredditPostsResponse); ok && response.Success {
            fmt.Printf("\nPosts in %s:\n", getPostsMsg[i].SubredditName)
        for _, post := range response.Posts {
            fmt.Printf("Title: %s\nAuthor: %s\nContent: %s\n\n",
                post.Title, post.AuthorId, post.Content)
        }
        }
    }

    // Keep the program running
    time.Sleep(1 * time.Second)
}
