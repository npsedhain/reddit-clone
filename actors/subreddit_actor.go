package actors

import (
	"reddit/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type SubredditActor struct {
    subreddits map[string]*Subreddit
}

type Subreddit struct {
    Name        string
    Description string
    CreatorId   string
    Members     map[string]bool // UserId -> membership status
}

func NewSubredditActor() *SubredditActor {
    return &SubredditActor{
        subreddits: make(map[string]*Subreddit),
    }
}

func (state *SubredditActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
        case *messages.CreateSubreddit:
            response := &messages.CreateSubredditResponse{}

            if _, exists := state.subreddits[msg.Name]; exists {
                response.Success = false
                response.Error = "Subreddit already exists"
            } else {
                // Create new subreddit
                subreddit := &Subreddit{
                    Name:        msg.Name,
                    Description: msg.Description,
                    CreatorId:   msg.CreatorId,
                    Members:     make(map[string]bool),
                }

                // Add creator as first member
                subreddit.Members[msg.CreatorId] = true

                // Store subreddit
                state.subreddits[msg.Name] = subreddit

                response.Success = true
                response.SubId = msg.Name
            }

            context.Respond(response)

        case *messages.JoinSubreddit:
            response := &messages.JoinSubredditResponse{}

            if subreddit, exists := state.subreddits[msg.SubredditName]; exists {
                if _, isMember := subreddit.Members[msg.UserId]; isMember {
                    response.Success = false
                    response.Error = "User is already a member"
                } else {
                    subreddit.Members[msg.UserId] = true
                    response.Success = true
                }
            } else {
                response.Success = false
                response.Error = "Subreddit not found"
            }

            context.Respond(response)

        case *messages.GetSubredditMembers:
            response := &messages.GetSubredditMembersResponse{}

            if subreddit, exists := state.subreddits[msg.SubredditName]; exists {
                members := make([]string, 0, len(subreddit.Members))
                for memberId := range subreddit.Members {
                    members = append(members, memberId)
                }
                response.Members = members
                response.Success = true
            } else {
                response.Success = false
                response.Error = "Subreddit not found"
            }

            context.Respond(response)

        case *messages.LeaveSubreddit:
            response := &messages.LeaveSubredditResponse{}

            if subreddit, exists := state.subreddits[msg.SubredditName]; exists {
                if msg.UserId == subreddit.CreatorId {
                    response.Success = false
                    response.Error = "Creator cannot leave their subreddit"
                } else if _, isMember := subreddit.Members[msg.UserId]; isMember {
                    delete(subreddit.Members, msg.UserId)
                    response.Success = true
                } else {
                    response.Success = false
                    response.Error = "User is not a member of this subreddit"
                }
            } else {
                response.Success = false
                response.Error = "Subreddit not found"
            }

            context.Respond(response)
    }
}

