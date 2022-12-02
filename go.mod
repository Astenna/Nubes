module github.com/Astenna/Thesis_PoC

go 1.19

replace github.com/Astenna/Thesis_PoC/faas v0.0.0 => ./faas

replace github.com/Astenna/Thesis_PoC/faas_lib v0.0.0 => ./faas_lib

require github.com/aws/aws-lambda-go v1.35.0
