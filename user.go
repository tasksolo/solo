package main

import (
	"fmt"

	"github.com/tasksolo/gosolo"
)

func printUser(user *gosolo.User) {
	fmt.Printf("%s %30s %s\n", getShortID(user), user.Email, user.Name)
}
