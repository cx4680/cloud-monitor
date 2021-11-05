package controllers

import (
	"testing"
)

type R1 struct {
	B int
}

type R2 struct {
	A   R1
	Age R3
}

type R3 struct {
	A    R1
	Name string
}

func TestDD(t *testing.T) {
	/*a:=[]int{1,2,3,4,5}
		var resources =make([]int,len(a))
		for index, item := range a {
			resources[index]=item
	        // resources=append(resources,item)
		}
		println(len(resources))*/
	a := R1{
		B: 1,
	}
	r := &R2{A: a, Age: R3{
		A: R1{
			B: 222,
		},
		Name: "XXXX",
	}}
	change(r)
	println(r)
}

func change(rss *R2) {
	rss.A.B = 3
	rss.Age.Name = "sssss"
	rss.Age.A.B = 666666
}

func TestAA(t *testing.T) {

}
