module github.com/Astenna/Thesis_PoC

replace github.com/Astenna/Thesis_PoC/FaaSLib v0.0.0 => ./faas_lib
replace github.com/Astenna/Thesis_PoC/FaaS v0.0.0 => ./FaaS

go 1.19

require github.com/Astenna/Thesis_PoC/FaaSLib v0.0.0
require github.com/Astenna/Thesis_PoC/FaaS v0.0.0