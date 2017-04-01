package main

import (
	"github.com/ngalayko/golangQL"
	"fmt"
	"encoding/json"
)

func main() {
	donald := Duck{
		Name: DuckName{
			FirstName: "Donald",
			LastName:  "Duck",
		},
		Nephews: []*Duck{&louie, &dewey, &huey},
	}

	jsonString, err := json.Marshal(donald)
	if err != nil {
		panic(err)
	}

	filtred, err := golangQL.Filter(donald, "{name {firstName lastName} nephews { name { firstName } hat } ")
	if err != nil {
		panic(err)
	}

	filteredJsonStr, err := json.Marshal(filtred)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonString))
	//"name":{"firstName":"Donald","lastName":"Duck"},"hat":"","nephews":[{"name":{"firstName":"Louie","lastName":"Duck"},"hat":"green","nephews":null},{"name":{"firstName":"Dewey","lastName":"Duck"},"hat":"blue","nephews":null},{"name":{"firstName":"Huey","lastName":"Duck"},"hat":"red","nephews":null}]}
	fmt.Println(string(filteredJsonStr))
	//{"name":{"firstName":"Donald","lastName":"Duck"},"nephews":[{"hat":"green","name":{"firstName":"Louie"}},{"hat":"blue","name":{"firstName":"Dewey"}},{"hat":"red","name":{"firstName":"Huey"}}]}
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
