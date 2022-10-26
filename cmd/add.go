package cmd

import (
	"log"
	"memlogy/database"
	"strings"
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
		fday, _ := cmd.Flags().GetInt16("go-to-day")

		if fday > 0 {
			panic("You can only add entry to previous days!")
		}

		var description = ""
		if fdescription != "" {
			description = fdescription
		} else if len(args) > 0 {
			for _, arg := range args {
				description = description + " " + arg
			}
		}
		entryRepository := database.EntryRepository()

		var entryDay = time.Now()
		if fday != 0 {
			entryDay = entryDay.AddDate(0, 0, int(fday))
		}

		entry := database.Entry{
			LogTime:     entryDay,
			Description: strings.TrimLeft(description, "*"),
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
	addCmd.Flags().Int16P("go-to-day", "g", 0, "Specify the entry day with positive or negative numbers for example \"--go-to-day=-1\" = yesterday ")
}
