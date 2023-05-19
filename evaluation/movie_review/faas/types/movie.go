package types

import "github.com/Astenna/Nubes/lib"

type Movie struct {
	Id              string
	Title           string
	ProductionYear  int
	Category        lib.Reference[Category]             `dynamodbav:",omitempty"`
	Reviews         lib.ReferenceNavigationList[Review] `nubes:"hasOne-Movie" dynamodbav:"-"`
	isInitialized   bool
	invocationDepth int
}

func (Movie) GetTypeName() string {
	return "Movie"
}
func (receiver *Movie) Init() {
	receiver.isInitialized = true
	receiver.Reviews = *lib.NewReferenceNavigationList[Review](lib.ReferenceNavigationListParam{OwnerId: receiver.Id, OwnerTypeName: receiver.GetTypeName(), OtherTypeName: (*new(Review)).GetTypeName(), ReferringFieldName: "Movie", IsManyToMany: false})
}
func (receiver *Movie) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.Id)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}
