package actors

import (
	"fmt"
	"reddit/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type PostActor struct {
    posts map[string]*messages.Post                 // PostId -> Post
    subredditPosts map[string][]string             // SubredditName -> []PostId
}

func NewPostActor() *PostActor {
    return &PostActor{
        posts: make(map[string]*messages.Post),
        subredditPosts: make(map[string][]string),
    }
}

func (state *PostActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.Post:
        response := &messages.CreatePostResponse{}

        // Generate post ID (simple implementation)
        postId := fmt.Sprintf("post_%d", time.Now().UnixNano())

        // Create post
        post := &messages.Post{
            PostId:        postId,
            Title:         msg.Title,
            Content:       msg.Content,
            AuthorId:      msg.AuthorId,
            SubredditName: msg.SubredditName,
            Timestamp:     time.Now().Unix(),
        }

        // Store post
        state.posts[postId] = post

        // Add to subreddit posts
        if _, exists := state.subredditPosts[msg.SubredditName]; !exists {
            state.subredditPosts[msg.SubredditName] = make([]string, 0)
        }
        state.subredditPosts[msg.SubredditName] = append(state.subredditPosts[msg.SubredditName], postId)

        response.Success = true
        response.PostId = postId

        context.Respond(response)

    case *messages.GetSubredditPosts:
        response := &messages.GetSubredditPostsResponse{}

        if postIds, exists := state.subredditPosts[msg.SubredditName]; exists {
            posts := make([]messages.Post, 0)
            for _, postId := range postIds {
                if post, ok := state.posts[postId]; ok {
                    posts = append(posts, *post)
                }
            }
            response.Success = true
            response.Posts = posts
        } else {
            response.Success = false
            response.Error = "No posts found for this subreddit"
        }

        context.Respond(response)
    }
}
