module faas_lib_test

go 1.19

replace github.com/Astenna/Nubes/lib v0.0.0 => ../../lib

replace github.com/Astenna/Nubes/faas v0.0.0 => ../faas

require github.com/Astenna/Nubes/lib v0.0.0

require github.com/Astenna/Nubes/faas v0.0.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/aws/aws-sdk-go v1.44.179 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/stretchr/testify v1.8.1
)
