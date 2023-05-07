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
	"task-list":     listWrapper[gosolo.Task]("task", listTaskOpts, printTask),
	"task-ls":       listWrapper[gosolo.Task]("task", listTaskOpts, printTask),

	"user-list": listWrapper[gosolo.User]("user", nil, printUser),
	"user-ls":   listWrapper[gosolo.User]("user", nil, printUser),
	// TODO: Add "user getshard"
}

func usage() {
	fmt.Printf("Usage: %s [ <flag> ... ] <noun> <verb> [ <arg> ... ]\n", os.Args[0])
	fmt.Printf("\ttask: complete, create, get, list, update\n")
	fmt.Printf("\tuser: create, get, list, update\n")
	os.Exit(1)
}

func main() {
	var err error

	ctx := context.Background()

	addHandlers[gosolo.Task]("task")
	addHandlers[gosolo.User]("user")

	profile := flag.String("profile", "default", "configuration profile name")
	base := flag.String("base", "", "API base URL")
	debug := flag.Bool("debug", false, "log HTTP details")
	insecure := flag.Bool("insecure", false, "allow invalid TLS certs")

	flag.Parse()

	cfgMap, err := readConfig()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if *profile == "" {
		for name := range cfgMap {
			fmt.Printf("%s\n", name)
		}

		return
	}

	if len(flag.Args()) < 2 {
		usage()
	}

	cfg := cfgMap[*profile]

	if cfg == nil {
		cfg = &gosolo.Config{}
	}

	if *debug {
		cfg.Debug = true
	}

	if *insecure {
		cfg.Insecure = true
	}

	if *base != "" {
		cfg.BaseURL = *base
	}

	getUserPass := func() (string, string, error) {
		user, err := readline.Line("Login email: ")
		if err != nil {
			return "", "", err
		}

		pass, err := readline.Password("Login password: ")
		if err != nil {
			return "", "", err
		}

		return user, string(pass), nil
	}

	c, err := gosolo.NewClient(ctx, cfg, getUserPass)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	cfgMap[*profile] = cfg

	err = writeConfig(cfgMap)
	if err != nil {
		fmt.Printf("Error writing config: %s\n", err)
		os.Exit(1)
	}

	noun := flag.Args()[0]
	verb := flag.Args()[1]

	handler := handlers[fmt.Sprintf("%s-%s", noun, verb)]
	if handler == nil {
		fmt.Printf("Unknown command: %s %s\n\n", noun, verb)
		usage()
	}

	err = handler(ctx, c, flag.Args()[2:])
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func addHandlers[T any](name string) {
	for _, cmd := range []string{"create", "make", "mk", "new"} {
		handlers[fmt.Sprintf("%s-%s", name, cmd)] = func(ctx context.Context, c *gosolo.Client, args []string) error {
			return create[T](ctx, c, name, args)
		}
	}

	for _, cmd := range []string{"cat", "get", "show"} {
		handlers[fmt.Sprintf("%s-%s", name, cmd)] = func(ctx context.Context, c *gosolo.Client, args []string) error {
			return get[T](ctx, c, name, args)
		}
	}

	for _, cmd := range []string{"change", "modify", "update"} {
		handlers[fmt.Sprintf("%s-%s", name, cmd)] = func(ctx context.Context, c *gosolo.Client, args []string) error {
			return update[T](ctx, c, name, args)
		}
	}
}
