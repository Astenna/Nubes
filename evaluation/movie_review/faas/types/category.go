package types

import "github.com/Astenna/Nubes/lib"

type Category struct {
	CName		string					`nubes:"id" dynamodbav:"Id"`
	Movies		lib.ReferenceNavigationList[Movie]	`nubes:"hasOne-Category" dynamodbav:"-"`
	isInitialized	bool
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
