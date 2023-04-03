module github.com/Astenna/Nubes/example

go 1.18

replace github.com/Astenna/Nubes/lib v0.0.0 => ../lib

require github.com/Astenna/Nubes/lib v0.0.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jftuga/geodist v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/aws/aws-sdk-go v1.44.179
	github.com/google/uuid v1.3.0
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/stretchr/testify v1.8.1
)
