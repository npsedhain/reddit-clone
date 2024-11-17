package actors

import (
	"reddit/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

type PostActor struct {
    posts map[string]*messages.Post                 // PostId -> Post
    subredditPosts map[string][]string             // SubredditName -> []PostId
    votes map[string]map[string]bool              // PostId -> UserId -> IsUpvote
}

func NewPostActor() *PostActor {
    return &PostActor{
        posts: make(map[string]*messages.Post),
        subredditPosts: make(map[string][]string),
        votes: make(map[string]map[string]bool),
    }
}

func (state *PostActor) handleVote(msg *messages.Vote) *messages.VoteResponse {
    if _, exists := state.posts[msg.TargetID]; !exists {
        return &messages.VoteResponse{Success: false, Error: "Post not found"}
    }

    if state.votes[msg.TargetID] == nil {
        state.votes[msg.TargetID] = make(map[string]bool)
    }

    // Check if user has already voted
    previousVote, hasVoted := state.votes[msg.TargetID][msg.UserID]
    if hasVoted {
        // If voting the same way, remove vote (toggle)
        if previousVote == msg.IsUpvote {
            delete(state.votes[msg.TargetID], msg.UserID)
        } else {
            // Change vote direction
            state.votes[msg.TargetID][msg.UserID] = msg.IsUpvote
        }
    } else {
        // New vote
        state.votes[msg.TargetID][msg.UserID] = msg.IsUpvote
    }

    return &messages.VoteResponse{Success: true}
}

func (state *PostActor) calculateKarmaChange(isUpvote bool) int {
    if isUpvote {
        return 1
    }
    return -1
}

func (state *PostActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.Post:
        response := &messages.CreatePostResponse{}

        postId := uuid.New().String()

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

    case *messages.Vote:
        response := state.handleVote(msg)

        // If vote was successful, notify user actor to update karma
        if response.Success {
            post := state.posts[msg.TargetID]
            karmaUpdate := &messages.UpdateKarma{
                UserID: post.AuthorId,
                Change: state.calculateKarmaChange(msg.IsUpvote),
            }
            context.Send(context.Parent(), karmaUpdate)
        }

        context.Respond(response)
    }
}
