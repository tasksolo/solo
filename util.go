package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopatchy/metadata"
	"github.com/gopatchy/path"
)

func getShortID(obj any) string {
	md := metadata.GetMetadata(obj)
	return md.ID[0:4]
}

func ifEmpty(s *string) string {
	if s == nil {
		return ""
	} else {
		return *s
	}
}

func color(c int, str string) string {
	b := ""

	if c > 7 {
		c = c % 8
		b = ";1"
	}

	return fmt.Sprintf("\033[1;3%d%sm%s\033[0m", c, b, str)
}

func setFields(obj any, args []string) error {
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf(`"%s" missing =`, arg)
		}

		err := path.Set(obj, parts[0], parts[1])
		if err != nil {
			return err
		}
	}

	return nil
}

func printDetails(obj any) {
	js, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(js))
}
