package main

import (
	"fmt"

	"github.com/Astenna/Nubes/movie_review_baseline/faas/db"
	"github.com/Astenna/Nubes/movie_review_baseline/faas/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func VoteOnReviewHandler(input models.ReviewVote) error {
	if input.ReviewId != "" && input.AccountId != "" {
		return fmt.Errorf("missing ReviewId or AccountId")
	}

	review, err := db.GetById[models.Review](input.ReviewId, "Review")
	if err != nil {
		return err
	}

	var previouslyUpvoted bool
	var previouslyDownvoted bool

	if _, exists := review.UpvotedBy[input.AccountId]; exists {
		previouslyUpvoted = true
	} else if _, exists := review.DownvotedBy[input.AccountId]; exists {
		previouslyDownvoted = true
	}

	if input.IsUpvote {
		if previouslyUpvoted {
			return fmt.Errorf("already upvoted")
		}
		if previouslyDownvoted {
			delete(review.DownvotedBy, input.AccountId)
		}

		review.UpvotedBy[input.AccountId] = struct{}{}
	} else {
		if previouslyDownvoted {
			return fmt.Errorf("already downvoted")
		}
		if previouslyUpvoted {
			delete(review.UpvotedBy, input.AccountId)
		}

		review.DownvotedBy[input.AccountId] = struct{}{}
	}

	return db.Upsert(review, input.ReviewId, "Review")
}

func main() {
	lambda.Start(VoteOnReviewHandler)
}
