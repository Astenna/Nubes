module github.com/Astenna/Thesis_PoC/faas

go 1.19

replace github.com/Astenna/Thesis_PoC/faas_lib v0.0.0 => ../faas_lib

require github.com/Astenna/Thesis_PoC/faas_lib v0.0.0

require (
	github.com/aws/aws-sdk-go v1.44.147 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/aws/aws-lambda-go v1.35.0 // indirect
)
