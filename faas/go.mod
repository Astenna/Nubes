module github.com/Astenna/Nubes/faas

go 1.19

replace github.com/Astenna/Nubes/lib v0.0.0 => ../lib

require github.com/Astenna/Nubes/lib v0.0.0

require github.com/davecgh/go-spew v1.1.1 // indirect

require (
	github.com/aws/aws-lambda-go v1.36.0
	github.com/aws/aws-sdk-go v1.44.147 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)
