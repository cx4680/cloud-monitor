package generate

import (
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	generate()
}
func TestEqual(t *testing.T) {
	s := ""
	var b string
	print(len(b))
	println(reflect.DeepEqual(s, b))
}

func BenchmarkStrEqual1(b *testing.B) {
	s := "a"
	var b1 string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		print(b1 == s)
	}
}

func BenchmarkStrEqual2(b *testing.B) {
	/*s := "a"
	var b1 string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		print(reflect.DeepEqual(s,b1))
	}*/
	if true && false {

	}
}
