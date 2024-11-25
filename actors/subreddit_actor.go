package actors

import (
	"math/rand"
	"reddit/messages"
	"time"

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
                response.ActorPID = context.Self()
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
                    response.SubId = msg.SubredditName
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
                if _, isMember := subreddit.Members[msg.UserId]; isMember {
                    delete(subreddit.Members, msg.UserId)
                    response.Success = true
                    response.SubId = msg.SubredditName
                } else {
                    response.Success = false
                    response.Error = "User is not a member of this subreddit"
                }
            } else {
                response.Success = false
                response.Error = "Subreddit not found"
            }

            context.Respond(response)

        case *messages.GetSubreddits:
            response := &messages.GetSubredditsResponse{}
            response.Success = true


            // Convert map keys to slice for random selection
            allSubreddits := make([]string, 0, len(state.subreddits))
            for subredditName := range state.subreddits {
                allSubreddits = append(allSubreddits, subredditName)
            }

            // Initialize random number generator
            rand := rand.New(rand.NewSource(time.Now().UnixNano()))

            // Get minimum between 5 and total number of subreddits
            numToReturn := 5
            if len(allSubreddits) < 5 {
                numToReturn = len(allSubreddits)
            }

            // Initialize result slice
            response.Subreddits = make([]string, 0, numToReturn)

            // Randomly select subreddits
            for i := 0; i < numToReturn; i++ {
                // Get random index
                randomIndex := rand.Intn(len(allSubreddits))
                // Add randomly selected subreddit
                response.Subreddits = append(response.Subreddits, allSubreddits[randomIndex])
                // Remove selected item to avoid duplicates
                allSubreddits = append(allSubreddits[:randomIndex], allSubreddits[randomIndex+1:]...)
            }

            context.Respond(response)
    }
}

