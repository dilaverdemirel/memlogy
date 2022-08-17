package cmd

import (
	"log"
	"memlogy/database"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add your notes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		fdescription, _ := cmd.Flags().GetString("description")

		var description = ""

		if fdescription != "" {
			description = fdescription
		} else if len(args) > 0 {
			for _, arg := range args {
				description = description + " " + arg
			}
		}
		entryRepository := database.EntryRepository()

		entry := database.Entry{
			LogTime:     time.Now(),
			Description: description,
		}

		_, err := entryRepository.Create(entry)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("description", "d", "", "Write something directly or -d \"Write Anything\" or --description \"Write Anything\"")
}
