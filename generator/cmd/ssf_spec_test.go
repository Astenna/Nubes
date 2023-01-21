package cmd

import (
	"testing"

	"github.com/Astenna/Nubes/generator/parser"
	tp "github.com/Astenna/Nubes/generator/template_parser"
)

// // cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// // BenchmarkPrimeNumbers-8   	     740	   1565620 ns/op	  269853 B/op	    6702 allocs/op
// // PASS
// // ok  	github.com/Astenna/Nubes/generator/cmd	2.834s
// // 1.56562 MILISECONDS
// func BenchmarkOrginal(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		typesPath := "../../faas/types"
// 		moduleName := "test"
// 		parsedPackage := parser.GetPackageTypes(typesPath, moduleName)
// 		stateChangingFuncs := parser.ParseStateChangingHandlers(typesPath, parsedPackage)
// 		_ = stateChangingFuncs
// 		parser.AddDBOperationsToMethods(typesPath, parsedPackage)
// 	}
// }

// cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// BenchmarkObjectOrientedDistributedSaves-8   	    1798	    832809 ns/op	  106037 B/op	    2645 allocs/op
// PASS
// ok  	github.com/Astenna/Nubes/generator/cmd	3.191s
// 0.832809 MILISECONDS, OR THEN 0.615473 MILISECONDS
// func BenchmarkObjectOrientedDistributedSaves(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		typesPath := "../../faas/types"
// 		moduleName := "test"
// 		parsedPackage, _ := parser.NewTypeSpecParser(typesPath)
// 		parsedPackage.Run(moduleName)
// 	}
// }

// cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// BenchmarkObjectOrientedOneSaveAll-8   	    2196	    585752 ns/op	  105078 B/op	    2623 allocs/op
// PASS
// ok  	github.com/Astenna/Nubes/generator/cmd	2.937s
// 2.80053 MILISECONDS
// func BenchmarkObjectOrientedOneSaveAll(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		typesPath := "../../faas/types"
// 		moduleName := "test"
// 		parsedPackage, _ := parser.NewTypeSpecParser(typesPath)
// 		parsedPackage.Run(moduleName)
// 	}
// }

// cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// BenchmarkObjectOrientedOneSaveSelected-8   	    1555	    655843 ns/op	  105285 B/op	    2625 allocs/op
// PASS
// ok  	github.com/Astenna/Nubes/generator/cmd	2.627s
// 0.655843 MILISECONDS
func BenchmarkObjectOrientedOneSaveSelected(b *testing.B) {
	for i := 0; i < b.N; i++ {
		typesPath := "../../faas/types"
		moduleName := "test"
		parsedPackage, _ := parser.NewTypeSpecParser(typesPath)
		parsedPackage.Run(moduleName)
	}
}

// cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// BenchmarkObjectOrientedOneSaveSelected-8   	    1555	    655843 ns/op	  105285 B/op	    2625 allocs/op
// PASS
// ok  	github.com/Astenna/Nubes/generator/cmd	2.627s
// 0.744222 MILISECONDS
func BenchmarkObjectOrientedOneSaveSelectedAfterInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		typesPath := "../../faas/types"
		moduleName := "test"
		parsedPackage, _ := parser.NewTypeSpecParser(typesPath)
		parsedPackage.Run(moduleName)
	}
}

// cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// BenchmarkObjectOrientedClient-8   	    1534	    757694 ns/op	  138601 B/op	    3411 allocs/op
// 0.757694 MS
func BenchmarkObjectOrientedClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		typesPath := "../../faas/types"
		clientTypesParser, _ := parser.NewClientTypesParser(tp.MakePathAbosoluteOrExitOnError(typesPath))
		clientTypesParser.Run()
	}
}

func BenchmarkObjectOrientedClientSeparateLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		typesPath := "../../faas/types"
		clientTypesParser, _ := parser.NewClientTypesParser(tp.MakePathAbosoluteOrExitOnError(typesPath))
		clientTypesParser.Run()
	}
}
