package actors

import (
	"reddit/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/uuid"
)

type CommentActor struct {
    comments map[string]messages.Comment
    // Map of postID to comment IDs
    postComments map[string][]string
		votes map[string]map[string]bool
}

func NewCommentActor() *CommentActor {
    return &CommentActor{
        comments:     make(map[string]messages.Comment),
        postComments: make(map[string][]string),
				votes: make(map[string]map[string]bool),
    }
}


func (state *CommentActor) handleVote(msg *messages.Vote) *messages.VoteResponse {
	if _, exists := state.comments[msg.TargetID]; !exists {
			return &messages.VoteResponse{Success: false, Error: "Comment not found"}
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

func (state *CommentActor) calculateKarmaChange(isUpvote bool) int {
	if isUpvote {
			return 1
	}
	return -1
}

func (c *CommentActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *messages.CreateComment:
        commentID := uuid.New().String()
        comment := messages.Comment{
            ID:        commentID,
            PostID:    msg.PostID,
            ParentID:  msg.ParentID,
            Content:   msg.Content,
            AuthorID:  msg.AuthorID,
            Children:  []messages.Comment{},
            CreatedAt: time.Now().Unix(),
        }

        c.comments[commentID] = comment
        c.postComments[msg.PostID] = append(c.postComments[msg.PostID], commentID)

        context.Respond(&messages.CreateCommentResponse{
            Success:   true,
            CommentID: commentID,
            ActorPID:  context.Self(),
        })

    case *messages.GetPostComments:
        comments := c.buildCommentTree(msg.PostID)
        context.Respond(&messages.GetPostCommentsResponse{
            Success:  true,
            Comments: comments,
        })

			case *messages.Vote:
        response := c.handleVote(msg)

        // If vote was successful, notify user actor to update karma
        if response.Success {
            comment := c.comments[msg.TargetID]
            karmaUpdate := &messages.UpdateKarma{
                UserID: comment.AuthorID,
                Change: c.calculateKarmaChange(msg.IsUpvote),
            }
            context.Send(context.Parent(), karmaUpdate)
        }

        context.Respond(response)
    }
}

func (c *CommentActor) buildCommentTree(postID string) []messages.Comment {
    // First, get all top-level comments (no parentID)
    var result []messages.Comment
    for _, commentID := range c.postComments[postID] {
        comment := c.comments[commentID]
        if comment.ParentID == "" {
            // Recursively build the comment tree
            comment.Children = c.getChildComments(commentID)
            result = append(result, comment)
        }
    }
    return result
}

func (c *CommentActor) getChildComments(parentID string) []messages.Comment {
    var children []messages.Comment
    for _, comment := range c.comments {
        if comment.ParentID == parentID {
            comment.Children = c.getChildComments(comment.ID)
            children = append(children, comment)
        }
    }
    return children
}
