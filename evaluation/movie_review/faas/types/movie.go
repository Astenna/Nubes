package types

import "github.com/Astenna/Nubes/lib"

type Movie struct {
	Id		string
	Title		string
	ProductionYear	int
	Category	lib.Reference[Category]
	Reviews		lib.ReferenceNavigationList[Review]	`nubes:"hasOne-Movie" dynamodbav:"-"`
	isInitialized	bool
}

func (Movie) GetTypeName() string {
	return "Movie"
}
func (receiver *Movie) Init() {
	receiver.isInitialized = true
	receiver.Reviews = *lib.NewReferenceNavigationList[Review](receiver.Id, receiver.GetTypeName(), "Movie", false)
}
