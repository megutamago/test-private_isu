package main

import (
	"fmt"
	//"reflect"
)

type Person struct {
    Name    string
    Age     int
    Address string
}

func main() {
    personData := map[string]map[string]string{
        "Alice": {
            "Age":     "25",
            "Address": "123 Main St",
        },
        "Bob": {
            "Age":     "30",
            "Address": "456 High St",
        },
        "Charlie": {
            "Age":     "35",
            "Address": "789 Elm St",
        },
    }

    var people []Person
    for name, data := range personData {
        var person Person
        person.Name = name
        //person.Age = parseInt(data["Age"])
        person.Address = data["Address"]
        people = append(people, person)
    }

    fmt.Println(people)
    //fmt.Println(personData)
}

func parseInt(str string) int {
    var result int
    _, err := fmt.Sscanf(str, "%d", &result)
    if err != nil {
        panic(err)
    }
    return result
}
