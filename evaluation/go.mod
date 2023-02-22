module github.com/Astenna/Nubes/evaluation

go 1.18

replace github.com/Astenna/Nubes/lib v0.0.0 => ../lib

require (
	github.com/Astenna/Nubes/lib v0.0.0
	github.com/aws/aws-sdk-go v1.44.179
	github.com/google/uuid v1.3.0
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	github.com/aws/aws-lambda-go v1.37.0
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
