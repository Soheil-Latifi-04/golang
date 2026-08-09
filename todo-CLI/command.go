package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CmdFlags struct {
	Add    string
	Delete int
	Edit   string
	Toggle int
	List   bool
}

func NewCmdFlags() *CmdFlags {
	cf := CmdFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo specify title")
	flag.StringVar(&cf.Edit, "edit", "", "Edit a todo by index and specify a new title. id:New_title")
	flag.IntVar(&cf.Delete, "del", -1, "Specify a todo with index to edit")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Specify a todo with index to toggle")
	flag.BoolVar(&cf.List, "list", false, "List all todos")

	flag.Parse()

	return &cf
}

func (cf *CmdFlags) Execute(todos *Todos) {
	switch {
	case cf.List:
		todos.print()
	case cf.Add != "":
		todos.add(cf.Add)
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error: Invalid format for edit. Please use id:new_title")
			os.Exit(1)
		}

		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error: Invalid index for edit")
			os.Exit(1)
		}
		todos.edit(index, parts[1])
	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)
	case cf.Delete != -1:
		todos.del(cf.Delete)

	default:
		fmt.Println("Invalid command")
	}
}
