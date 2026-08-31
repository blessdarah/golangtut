package main

import (
	"blessdarah/tuts/user"
	"fmt"
)

func main() {
	store := user.NewUserStore()

	addr := user.Address{
		Zip:    "12345",
		City:   "jakarta",
		Street: "jakarta",
	}

	bless := user.User{
		Name:    "blessdarah",
		Age:     19,
		Email:   "blessdarah@gmail.com",
		Address: addr,
	}

	err := store.Add(bless)
	if err != nil {
		fmt.Println(err)
		return
	}

	users := store.GetAll()

	fmt.Println(users)

}
