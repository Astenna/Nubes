package types

import (
	"fmt"

	"github.com/Astenna/Nubes/lib"
)

type Review struct {
	Id            string
	Rating        int
	Movie         lib.Reference[Movie]
	Reviewer      lib.Reference[Account]
	Text          string
	DownvotedBy   map[string]struct{} `nubes:"readonly"`
	UpvotedBy     map[string]struct{} `nubes:"readonly"`
	MapField      map[string]string
	isInitialized bool
}

func (Review) GetTypeName() string {
	return "Review"
}

func (m *Review) Downvote(account Account) (int, error) {
	if m.isInitialized {
		tempReceiverName, _libError := lib.GetObjectState[Review](m.Id)
		if _libError != nil {
			return *new(int), _libError
		}
		m = tempReceiverName
		m.Init()
	}

	// _, _ = m.DownvotedBy["account.GetId()"]
	// _, _ = m.DownvotedBy[account.Email]
	// _, _ = m.DownvotedBy[account.GetId()]

	if _, exists := m.DownvotedBy[account.GetId()]; exists {
		return len(m.DownvotedBy), fmt.Errorf("the user have already downvoted")
	}

	delete(m.UpvotedBy, account.GetId())

	m.DownvotedBy[account.GetId()] = struct{}{}
	if m.isInitialized {
		_libError := lib.Upsert(m, m.Id)
		if _libError != nil {
			return *new(int), _libError
		}
	}
	return len(m.DownvotedBy), nil
}

func (m *Review) Upvote(account Account) (int, error) {
	if m.isInitialized {
		tempReceiverName, _libError := lib.GetObjectState[Review](m.Id)
		if _libError != nil {
			return *new(int), _libError
		}
		m = tempReceiverName
		m.Init()
	}
	if _, exists := m.UpvotedBy[account.GetId()]; exists {
		return len(m.UpvotedBy), fmt.Errorf("the user have already upvoted")
	}

	delete(m.DownvotedBy, account.GetId())

	m.UpvotedBy[account.GetId()] = struct{}{}
	if m.isInitialized {
		_libError := lib.Upsert(m, m.Id)
		if _libError != nil {
			return *new(int), _libError
		}
	}
	return len(m.DownvotedBy), nil
}

func (receiver *Review) Init() {
	receiver.isInitialized = true
}
