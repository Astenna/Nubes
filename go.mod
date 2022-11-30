module github.com/Astenna/Thesis_PoC

go 1.19

replace github.com/Astenna/Thesis_PoC/faas v0.0.0 => ./faas
replace github.com/Astenna/Thesis_PoC/faas_lib v0.0.0 => ./faas_lib

require github.com/Astenna/Thesis_PoC/faas_lib v0.0.0
require github.com/Astenna/Thesis_PoC/faas v0.0.0
require github.com/aws/aws-sdk-go v1.44.147