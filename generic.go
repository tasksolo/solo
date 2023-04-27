package main

import (
	"context"
	"fmt"

	"github.com/gopatchy/metadata"
	"github.com/tasksolo/gosolo"
)

func create[T any](ctx context.Context, c *gosolo.Client, name string, args []string) error {
	create := new(T)

	err := setFields(create, args)
	if err != nil {
		return err
	}

	created, err := gosolo.CreateName[T](ctx, c, name, create)
	if err != nil {
		return err
	}

	printDetails(created)

	return nil
}

func get[T any](ctx context.Context, c *gosolo.Client, name string, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("args: <id>")
	}

	shortID := args[0]

	obj, err := gosolo.FindName[T](ctx, c, name, shortID)
	if err != nil {
		return err
	}

	printDetails(obj)

	return nil
}

func listWrapper[T any](name string, optsCB func() *gosolo.ListOpts[T], printCB func(*T)) func(context.Context, *gosolo.Client, []string) error {
	return func(ctx context.Context, c *gosolo.Client, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("args: none")
		}

		opts := &gosolo.ListOpts[T]{}
		if optsCB != nil {
			opts = optsCB()
		}

		objs, err := gosolo.ListName[T](ctx, c, name, opts)
		if err != nil {
			return err
		}

		for _, obj := range objs {
			printCB(obj)
		}

		return nil
	}
}

func update[T any](ctx context.Context, c *gosolo.Client, name string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("args: <id> <field>=<value> ...") //nolint:revive,stylecheck
	}

	shortID := args[0]
	fields := args[1:]

	obj, err := gosolo.FindName[T](ctx, c, name, shortID)
	if err != nil {
		return err
	}

	md := metadata.GetMetadata(obj)

	update := new(T)

	err = setFields(update, fields)
	if err != nil {
		return err
	}

	updated, err := gosolo.UpdateName[T](ctx, c, name, md.ID, update, nil)
	if err != nil {
		return err
	}

	printDetails(updated)

	return nil
}
