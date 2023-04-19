package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/tasksolo/gosolo"
)

type handler func(ctx context.Context, c *gosolo.Client, args []string) error

var handlers = map[string]handler{
	"task-complete": completeTask,
	"task-done":     completeTask,
	"task-list":     listWrapper[gosolo.TaskResponse]("task", listTaskOpts, printTask),
	"task-ls":       listWrapper[gosolo.TaskResponse]("task", listTaskOpts, printTask),

	"user-list": listWrapper[gosolo.UserResponse]("user", nil, printUser),
	"user-ls":   listWrapper[gosolo.UserResponse]("user", nil, printUser),
}

func usage() {
	fmt.Printf("Usage: %s [ <flag> ... ] <noun> <verb> [ <arg> ... ]\n", os.Args[0])
	fmt.Printf("\ttask: complete, create, get, list, update\n")
	fmt.Printf("\tuser: create, get, list, update\n")
	os.Exit(1)
}

func main() {
	addHandlers[gosolo.TaskResponse, gosolo.TaskRequest]("task")
	addHandlers[gosolo.UserResponse, gosolo.UserRequest]("user")

	base := flag.String("base", "https://api.solø.com", "API base URL")
	debug := flag.Bool("debug", false, "log HTTP details")
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
	}

	noun := flag.Args()[0]
	verb := flag.Args()[1]

	handler := handlers[fmt.Sprintf("%s-%s", noun, verb)]
	if handler == nil {
		fmt.Printf("Unknown command: %s %s\n\n", noun, verb)
		usage()
	}

	c := gosolo.NewClient(*base).
		SetDebug(*debug)

	ctx := context.Background()

	cfg, err := readConfig()
	if err != nil {
		user, err := readline.Line("Login email: ")
		if err != nil {
			fmt.Printf("Error reading login email: %s\n", err)
			os.Exit(1)
		}

		pass, err := readline.Password("Login password: ")
		if err != nil {
			fmt.Printf("Error reading login password: %s\n", err)
			os.Exit(1)
		}

		token, err := c.Auth(ctx, user, string(pass))
		if err != nil {
			fmt.Printf("Error authenticating: %s\n", err)
			os.Exit(1)
		}

		cfg = &config{
			Token: &token,
		}

		err = writeConfig(cfg)
		if err != nil {
			fmt.Printf("Error writing config: %s\n", err)
			os.Exit(1)
		}
	}

	c.SetAuthToken(*cfg.Token)

	err = handler(ctx, c, flag.Args()[2:])
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func addHandlers[TOut, TIn any](name string) {
	for _, cmd := range []string{"create", "make", "mk", "new"} {
		handlers[fmt.Sprintf("%s-%s", name, cmd)] = func(ctx context.Context, c *gosolo.Client, args []string) error {
			return create[TOut, TIn](ctx, c, name, args)
		}
	}

	for _, cmd := range []string{"cat", "get", "show"} {
		handlers[fmt.Sprintf("%s-%s", name, cmd)] = func(ctx context.Context, c *gosolo.Client, args []string) error {
			return get[TOut](ctx, c, name, args)
		}
	}

	for _, cmd := range []string{"change", "modify", "update"} {
		handlers[fmt.Sprintf("%s-%s", name, cmd)] = func(ctx context.Context, c *gosolo.Client, args []string) error {
			return update[TOut, TIn](ctx, c, name, args)
		}
	}
}
