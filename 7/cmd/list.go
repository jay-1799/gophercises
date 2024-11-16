package cmd

import (
	"fmt"
	"gophercises/7/db"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong", err.Error())
			os.Exit(1)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete!")
			return
		}
		fmt.Println("you have the following tasks:")
		for i, task := range tasks {
			fmt.Printf("$d. %s \n", i+1, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
