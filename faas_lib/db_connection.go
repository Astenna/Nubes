package faas_lib

import "github.com/aws/aws-sdk-go/aws/session"

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
