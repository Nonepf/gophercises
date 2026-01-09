package main

import (
	"fmt"
	"os"
	"strings"

	"strconv"

	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	// 注册指令
	// rootCmd -> addCmd
	// 		   -> listCmd
	//		   -> completeCmd

	var tasks []string

	var rootCmd = &cobra.Command{
		Use:   "task",
		Short: "Task 是一个 CLI 任务管理器",
		Long:  `哇啦`,
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new task to your TODO list",

		Run: func(cmd *cobra.Command, args []string) {
			task := strings.Join(args, " ")
			tasks = append(tasks, task)
			fmt.Printf("Added \"%s\" to your TODO list\n", task)
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all the tasks that remain to be finished",

		Run: func(cmd *cobra.Command, args []string) {
			for i, task := range tasks {
				fmt.Printf("%d: %s\n", i+1, task)
			}
		},
	}

	var completeCmd = &cobra.Command{
		Use:   "complete",
		Short: "Complete task with certain index",

		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil || (id <= 0 || id > len(tasks)) {
				fmt.Printf("ERROR: Invalid Number!")
				return
			}
			id = id - 1
			tasks[id] = tasks[len(tasks)-1]
			tasks = tasks[:len(tasks)-1]
		},
	}

	rootCmd.AddCommand(addCmd, listCmd, completeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return rootCmd
}

func main() {
	initCmd()
}
