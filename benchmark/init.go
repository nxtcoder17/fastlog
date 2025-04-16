package benchmark

import "fmt"

type user struct {
	FirstName string
	LastName  string
	Age       int
}

var attrs = []any{
	"string", "benchmark",
	"number", 17,
	"floating.number", 17.17,
	"bool", true,
	"map",
	map[string]any{"hello": "world"},
	"array",
	[]any{1, 2, 3, 4},
	"err", fmt.Errorf("this is an error"),
	"user",
	user{FirstName: "sample", LastName: "kumar", Age: 17},
	"large-map",
	map[string]any{
		"k1":  "v1",
		"k2":  "v1",
		"k3":  "v1",
		"k4":  "v1",
		"k5":  "v1",
		"k6":  "v1",
		"k8":  "v1",
		"k9":  "v1",
		"k10": "v1",
		"k11": "v1",
		"k12": "v1",
		"k13": "v1",
	},
}
