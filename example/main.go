package main

import (
	"github.com/ngalayko/golangQL"
	"fmt"
)

func main() {
	donald := Duck{
		Name: DuckName{
			FirstName: "Donald",
			LastName:  "Duck",
		},
		Nephews: []*Duck{&louie, &dewey, &huey},
	}

	filtred, err := golangQL.Filter(donald, "{name {firstName lastName} nephews { name { firstName } hat } ")
	if err != nil {
		panic(err)
	}

	fmt.Println(filtred)
}

type DuckName struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Duck struct {
	Name    DuckName `json:"name"`
	Hat     string   `json:"hat"`
	Nephews []*Duck  `json:"nephews"`
}

var (
	huey = Duck{
		Name: DuckName{
			FirstName: "Huey",
			LastName:  "Duck",
		},
		Hat: "red",
	}

	dewey = Duck{
		Name: DuckName{
			FirstName: "Dewey",
			LastName:  "Duck",
		},
		Hat: "blue",
	}

	louie = Duck{
		Name: DuckName{
			FirstName: "Louie",
			LastName:  "Duck",
		},
		Hat: "green",
	}
)
