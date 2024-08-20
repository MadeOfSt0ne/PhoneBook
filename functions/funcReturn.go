package main

import "fmt"

func funcRet(i int) func(int) int {
	if i < 0 {
		return func(k int) int {
			k = -k
			return k + k
		}
	}

	return func(k int) int {
		return k * k
	}
}

func main() {
	n := 10
	i := funcRet(n)
	j := funcRet(-4)

	fmt.Printf("%T\n", i)
	fmt.Printf("%T %v\n", j, j)
	fmt.Println("j", j, j(-5))

	fmt.Println(i(10))
	fmt.Println(j(10))
}
