package types

import "github.com/Astenna/Nubes/lib"

type Category struct {
	CName		string					`nubes:"id" dynamodbav:"Id"`
	Movies		lib.ReferenceNavigationList[Movie]	`nubes:"hasOne-Category" dynamodbav:"-"`
	isInitialized	bool
	invocationDepth	int
}

func (Category) GetTypeName() string {
	return "Category"
}

func (u Category) GetId() string {
	return u.CName
}
func (receiver *Category) Init() {
	receiver.isInitialized = true
	receiver.Movies = *lib.NewReferenceNavigationList[Movie](receiver.CName, receiver.GetTypeName(), "Category", false)
}
func (receiver *Category) saveChangesIfInitialized() error {
	if receiver.isInitialized && receiver.invocationDepth == 1 {
		_libError := lib.Upsert(receiver, receiver.CName)
		if _libError != nil {
			return _libError
		}
	}
	return nil
}
