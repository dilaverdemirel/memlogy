package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"memlogy/database"
	"memlogy/util"

	"github.com/eiannone/keyboard"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your notes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		entryRepository := database.EntryRepository()
		day := time.Now()
		currentDataList := load(day, *entryRepository)

		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
		}()

		editing := false
		for {

			if !editing {
				_, key, err := keyboard.GetKey()
				if err != nil {
					panic(err)
				}
				if key == keyboard.KeyEsc {
					break
				}

				if key == keyboard.KeyArrowLeft {
					day = day.AddDate(0, 0, -1)
					currentDataList = load(day, *entryRepository)
				}

				if key == keyboard.KeyArrowRight {
					day = day.AddDate(0, 0, 1)
					currentDataList = load(day, *entryRepository)
				}

				if key == keyboard.KeyDelete {
					editing = true
					fmt.Print("Please select entry index : ")
				}
			} else {
				editing = handleEditing(currentDataList, *entryRepository, day)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func load(day time.Time, entryRepository database.SQLiteRepository) [][]string {
	util.ClearScreen()
	fmt.Println("------------------------------------------------------")
	fmt.Println("                    Selected Day")
	fmt.Println("                    " + day.Format("2 Jan 2006"))
	fmt.Println("------------------------------------------------------")

	fmt.Println("")

	all, err := entryRepository.GetByDay(day)
	if err != nil {
		log.Fatal(err)
	}

	var data = make([][]string, 0)
	for index, website := range all {
		line := make([]string, 0)
		line = append(line, fmt.Sprintf("%d", index+1))
		line = append(line, website.LogTime.Format("2 Jan 2006 15:04"))
		line = append(line, website.Description)
		line = append(line, fmt.Sprintf("%d", website.Id))
		data = append(data, line)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Index", "Log Time", "Description", "ID"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgWhiteColor, tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.FgWhiteColor, tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.FgWhiteColor, tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.FgWhiteColor, tablewriter.Bold, tablewriter.BgGreenColor})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	fmt.Println("- Press ESC to quit")
	fmt.Println("- You can use left and right arrow")
	fmt.Println("- Press Delete to delete an entry")
	fmt.Println("")

	return data
}

func handleEditing(
	currentDataList [][]string,
	entryRepository database.SQLiteRepository,
	day time.Time) bool {

	editing := true
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	index := ""
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		index = index + fmt.Sprintf("%c", event.Rune)

		//Delete selected row
		if event.Key == keyboard.KeyEnter {
			selectedRow, err1 := strconv.Atoi(index[:len(index)-1])
			if err1 != nil {
				panic(err1)
			}
			if selectedRow <= len(currentDataList) {
				selectedRow = selectedRow - 1
				rowId, _ := strconv.ParseInt(currentDataList[selectedRow][3], 10, 64)
				entryRepository.Delete(rowId)
			}
			load(day, entryRepository)
			editing = false
			break
		}

		//Back to list
		if event.Key == keyboard.KeyEsc {
			editing = false
			break
		}

		//Delete entered index
		if event.Key == keyboard.KeyBackspace || event.Key == keyboard.KeyBackspace2 {
			length := len(index)
			if length == 1 {
				index = ""
			} else if length > 1 {
				index = index[:length-2]
			}
		}
		util.ClearCurrentLine()
		fmt.Print("Please select entry index : " + index)
	}

	return editing
}
