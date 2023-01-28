package client_lib

type MovieStub struct {
	Id string

	Title string

	ProductionYear int

	Category Reference[CategoryStub]
}

func (MovieStub) GetTypeName() string {
	return "Movie"
}

type ReviewStub struct {
	Id string

	Rating int

	Movie Reference[MovieStub]

	Reviewer Reference[AccountStub]

	Text string

	DownvotedBy map[string]struct{} `nubes:"readonly"`

	UpvotedBy map[string]struct{} `nubes:"readonly"`

	MapField map[string]string
}

func (ReviewStub) GetTypeName() string {
	return "Review"
}

type AccountStub struct {
	Nickname string

	Email string `nubes:"id,readonly" dynamodbav:"Id"`

	Password string `nubes:"readonly"`
}

func (AccountStub) GetTypeName() string {
	return "Account"
}

type CategoryStub struct {
	CName string `nubes:"id" dynamodbav:"Id"`
}

func (CategoryStub) GetTypeName() string {
	return "Category"
}
