package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/tasksolo/gosolo"
)

var (
	includeComplete = flag.Bool("include-complete", false, "include complete tasks")
	includeFuture   = flag.Bool("include-future", false, "include tasks with after > now")
)

func listTaskOpts() *gosolo.ListOpts[gosolo.Task] {
	opts := &gosolo.ListOpts[gosolo.Task]{}

	if *includeComplete {
		opts.Sorts = append(opts.Sorts, "+complete")
	} else {
		opts.Filters = append(opts.Filters, gosolo.Filter{
			Path:  "complete",
			Op:    "eq",
			Value: "false",
		})
	}

	if *includeFuture {
		opts.Sorts = append(opts.Sorts, "+after")
	} else {
		opts.Filters = append(opts.Filters, gosolo.Filter{
			Path:  "after",
			Op:    "lte",
			Value: "now",
		})
	}

	return opts
}

func completeTask(ctx context.Context, c *gosolo.Client, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("args: <id>")
	}

	shortID := args[0]

	task, err := c.FindTask(ctx, shortID)
	if err != nil {
		return err
	}

	complete := true
	updated, err := c.UpdateTask(ctx, task.ID, &gosolo.Task{
		Complete: complete,
	}, nil)
	if err != nil {
		return err
	}

	printDetails(updated)

	return nil
}

func printTask(task *gosolo.Task) {
	complete := "o"
	name := task.Name

	if task.Complete {
		complete = color(2, "ø")
		name = color(7, name)
	} else if task.After.After(time.Now()) {
		name = color(5, name)
	}

	fmt.Printf("%s %s %s\n", color(6, getShortID(task)), complete, name)
}
