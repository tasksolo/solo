package main

import (
	"context"
	"fmt"

	"github.com/firestuff/patchy/metadata"
	"github.com/tasksolo/gosolo"
)

func create[TOut, TIn any](ctx context.Context, c *gosolo.Client, name string, args []string) error {
	create := new(TIn)

	err := setFields(create, args)
	if err != nil {
		return err
	}

	created, err := gosolo.CreateName[TOut, TIn](ctx, c, name, create)
	if err != nil {
		return err
	}

	printDetails(created)

	return nil
}

func get[TOut any](ctx context.Context, c *gosolo.Client, name string, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("args: <id>")
	}

	shortID := args[0]

	obj, err := gosolo.FindName[TOut](ctx, c, name, shortID)
	if err != nil {
		return err
	}

	printDetails(obj)

	return nil
}

func listWrapper[TOut any](name string, optsCB func() *gosolo.ListOpts, printCB func(*TOut)) func(context.Context, *gosolo.Client, []string) error {
	return func(ctx context.Context, c *gosolo.Client, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("args: none")
		}

		opts := &gosolo.ListOpts{}
		if optsCB != nil {
			opts = optsCB()
		}

		objs, err := gosolo.ListName[TOut](ctx, c, name, opts)
		if err != nil {
			return err
		}

		for _, obj := range objs {
			printCB(obj)
		}

		return nil
	}
}

func update[TOut, TIn any](ctx context.Context, c *gosolo.Client, name string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("args: <id> <field>=<value> ...")
	}

	shortID := args[0]
	fields := args[1:]

	obj, err := gosolo.FindName[TOut](ctx, c, name, shortID)
	if err != nil {
		return err
	}

	md := metadata.GetMetadata(obj)

	update := new(TIn)

	err = setFields(update, fields)
	if err != nil {
		return err
	}

	updated, err := gosolo.UpdateName[TOut, TIn](ctx, c, name, md.ID, update, nil)
	if err != nil {
		return err
	}

	printDetails(updated)

	return nil
}
