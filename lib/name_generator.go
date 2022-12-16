package lib

// GetTableName generates DynamoDB's
// table name according to object's type
// cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// BenchmarkReflection-8           422089638                2.837 ns/op
// BenchmarkNew-8                  1000000000               0.2520 ns/op
//func GetTableName[T any]() string {
//	return reflect.TypeOf((*T)(nil)).Name()
//}
