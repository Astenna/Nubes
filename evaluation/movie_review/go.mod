module github.com/Astenna/Nubes/movie_review

replace github.com/Astenna/Nubes/lib v0.0.0 => ../../lib

require (
	github.com/Astenna/Nubes/lib v0.0.0
	github.com/aws/aws-lambda-go v1.37.0
	github.com/aws/aws-sdk-go v1.44.179
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

go 1.19
