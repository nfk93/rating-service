package tests

import (
	"math/rand"
)

type TestData struct {
	Name string
}

func defaultData() TestData {
	return TestData{
		Name: randStr(60),
	}
}

func randStr(n int) string {
	str := make([]rune, n)
	for i := range str {
		str[i] = rune(65 + rand.Intn(26))
	}

	return string(str)
}
