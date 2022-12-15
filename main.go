package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, er := New(dir, nil)
	if er {
		fmt.Println("Error", er)
	}

	employees := []User{
		{"Tim", "23", "8899665544", "My tech", Address{"indore", "Madhya Pradesh", "India", "998855663322"}},
		{"Sam", "24", "9966345544", "techBro", Address{"Chennai", "Chennai", "India", "998812663322"}},
		{"Tom", "28", "9924345544", "techBro2", Address{"City", "State", "India", "997612663322"}},
	}

	for _, value := range employees {
		db.write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, er := db.ReadAll("users")
	if er != nil {
		fmt.Println("Error", er)
	}
	fmt.Println(records)
}
