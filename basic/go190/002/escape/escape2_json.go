package main

import (
	"encoding/json"
	"fmt"
)

type user_json struct {
	ID   int
	Name string
}

func main() {
	retriFunc := []func(int)(*user_json, error){ retrieveUserV1, retrieveUserV2}
	for _, f := range retriFunc {
		u, err := f(1234)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%+v\n",*u)
	}
}

func retrieveUserV1(id int) (*user_json, error) {
	r, err := getUser(id)
	if err != nil {
		return nil, err
	}

	var u *user_json
	err = json.Unmarshal([]byte(r), &u)  // must share the pointer variable with the json.Unmarshal call
	                                     // The json.Unmarshal call will create the user value and assign
	                                     // its address to the pointer variable.
	return u, err
}

func retrieveUserV2(id int) (*user_json, error) {
	r, err := getUser(id)
	if err != nil {
		return nil, err
	}

	var u user_json
	err = json.Unmarshal([]byte(r), &u)  // Comparing with V1, better readability.
	return &u, err
}

func getUser(id int) (string, error) {
	response := fmt.Sprintf(`{"id": %d, "name": "sally"}`, id)
	return response, nil
}