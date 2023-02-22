package models

type Review struct {
	Id          string
	Rating      int
	MovieId     string
	ReviewerId  string
	Text        string
	DownvotedBy map[string]struct{}
	UpvotedBy   map[string]struct{}
}

type ReviewVote struct {
	ReviewId  string
	AccountId string
	IsUpvote  bool
}
