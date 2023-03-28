package types

import (
	"fmt"

	"github.com/Astenna/Nubes/lib"
)

type Review struct {
	Id		string
	Rating		int
	Movie		lib.Reference[Movie]	`dynamodbav:",omitempty"`
	Reviewer	lib.Reference[Account]
	Text		string
	DownvotedBy	map[string]struct{}	`nubes:"readonly"`
	UpvotedBy	map[string]struct{}	`nubes:"readonly"`
	MapField	map[string]string
	isInitialized	bool
	invocationDepth	int
}

func (Review) GetTypeName() string {
	return "Review"
}

func (m *Review) Downvote(account Account) (int, error) {
	m.invocationDepth++
	if m.isInitialized && m.invocationDepth == 1 {
		_libError := lib.GetObjectState(m.Id, m)
		if _libError != nil {
			m.invocationDepth--
			return *new(int), _libError
		}
	}

	if _, exists := m.DownvotedBy[account.GetId()]; exists {
		m.invocationDepth--
		return len(m.DownvotedBy), fmt.Errorf("the user have already downvoted")
	}

	delete(m.UpvotedBy, account.GetId())
	if m.isInitialized {
		_libError := lib.Upsert(m, m.Id)
		if _libError != nil {
			m.invocationDepth--
			return *new(int), _libError
		}
	}
	_libUpsertError := m.saveChangesIfInitialized()
	m.invocationDepth--

	return len(m.DownvotedBy), _libUpsertError
}

func (m *Review) Upvote(account Account) (int, error) {
	m.invocationDepth++
	if m.isInitialized && m.invocationDepth == 1 {
		_libError := lib.GetObjectState(m.Id, m)
		if _libError != nil {
			m.invocationDepth--
			return *new(int), _libError
		}
	}
	if _, exists := m.UpvotedBy[account.GetId()]; exists {
		m.invocationDepth--
		return len(m.UpvotedBy), fmt.Errorf("the user have already upvoted")
	}

	delete(m.DownvotedBy, account.GetId())
	if m.isInitialized {
		_libError := lib.Upsert(m, m.Id)
		if _libError != nil {
			m.invocationDepth--
			return *new(int), _libError
		}
	}
	_libUpsertError := m.saveChangesIfInitialized()
	m.invocationDepth--
	return len(m.DownvotedBy), _libUpsertError
}
func (receiver *Review) Init() {
	receiver.isInitialized = true
}
func (receiver *Review) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}
