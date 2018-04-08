// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/tealeg/xlsx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type card struct {
	ID string `json:"id"`
}

type cardDetails struct {
	ID     string `json:"id"`
	Badges struct {
		Votes             int `json:"votes"`
		AttachmentsByType struct {
			Trello struct {
				Board int `json:"board"`
				Card  int `json:"card"`
			} `json:"trello"`
		} `json:"attachmentsByType"`
		ViewingMemberVoted bool      `json:"viewingMemberVoted"`
		Subscribed         bool      `json:"subscribed"`
		Fogbugz            string    `json:"fogbugz"`
		CheckItems         int       `json:"checkItems"`
		CheckItemsChecked  int       `json:"checkItemsChecked"`
		Comments           int       `json:"comments"`
		Attachments        int       `json:"attachments"`
		Description        bool      `json:"description"`
		Due                time.Time `json:"due"`
		DueComplete        bool      `json:"dueComplete"`
	} `json:"badges"`
	CheckItemStates  []interface{} `json:"checkItemStates"`
	Closed           bool          `json:"closed"`
	DueComplete      bool          `json:"dueComplete"`
	DateLastActivity time.Time     `json:"dateLastActivity"`
	Desc             string        `json:"desc"`
	DescData         struct {
		Emoji struct {
		} `json:"emoji"`
	} `json:"descData"`
	Due          time.Time     `json:"due"`
	Email        interface{}   `json:"email"`
	IDBoard      string        `json:"idBoard"`
	IDChecklists []interface{} `json:"idChecklists"`
	Members      []struct {
		ID         string      `json:"id"`
		AvatarHash interface{} `json:"avatarHash"`
		FullName   string      `json:"fullName"`
		Initials   string      `json:"initials"`
		Username   string      `json:"username"`
	} `json:"members"`
	IDList            string        `json:"idList"`
	IDMembers         []string      `json:"idMembers"`
	IDMembersVoted    []interface{} `json:"idMembersVoted"`
	IDShort           int           `json:"idShort"`
	IDAttachmentCover interface{}   `json:"idAttachmentCover"`
	Labels            []struct {
		ID      string `json:"id"`
		IDBoard string `json:"idBoard"`
		Name    string `json:"name"`
		Color   string `json:"color"`
	} `json:"labels"`
	Limits struct {
		Attachments struct {
			PerCard struct {
				Status    string `json:"status"`
				DisableAt int    `json:"disableAt"`
				WarnAt    int    `json:"warnAt"`
			} `json:"perCard"`
		} `json:"attachments"`
		Checklists struct {
			PerCard struct {
				Status    string `json:"status"`
				DisableAt int    `json:"disableAt"`
				WarnAt    int    `json:"warnAt"`
			} `json:"perCard"`
		} `json:"checklists"`
		Stickers struct {
			PerCard struct {
				Status    string `json:"status"`
				DisableAt int    `json:"disableAt"`
				WarnAt    int    `json:"warnAt"`
			} `json:"perCard"`
		} `json:"stickers"`
	} `json:"limits"`
	IDLabels              []string `json:"idLabels"`
	ManualCoverAttachment bool     `json:"manualCoverAttachment"`
	Name                  string   `json:"name"`
	Pos                   float32  `json:"pos"`
	ShortLink             string   `json:"shortLink"`
	ShortURL              string   `json:"shortUrl"`
	Subscribed            bool     `json:"subscribed"`
	URL                   string   `json:"url"`
}

var (
	boardName string
	listName  string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cards := getCards(args)
		writeExcel(cards)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exportCmd.PersistentFlags().StringVar(&boardName, "b", "", "name of the board")
	exportCmd.PersistentFlags().StringVar(&listName, "l", "", "name of the list")

}

func getCards(args []string) []card {
	var trelloCards []card

	if boardName == "" {
		log.Fatal("Specified board name is empty")
	}
	if listName == "" {
		log.Fatal("Specified list name is empty")
	}
	listID := getListID(boardName, listName)

	apiKey = viper.Get("apiKey")
	clientToken = viper.Get("token")
	url := fmt.Sprintf("%s/1/lists/%s/cards?key=%s&token=%s", trelloURL, listID, apiKey, clientToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error occured: ", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("bad request sent to trello api")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&trelloCards); err != nil {
		log.Fatal("Error decoding trello response: ", err)
	}
	return trelloCards
}

func writeExcel(cards []card) {
	apiKey = viper.Get("apiKey")
	clientToken = viper.Get("token")
	client := &http.Client{}
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Trello Tasks")
	if err != nil {
		log.Fatal("Could not create xcel file")
	}
	row = sheet.AddRow()
	headers := []string{"Name", "Labels", "Description", "Last Activity Date", "Closed", "Due Date", "Due complete", "No.of comments", "Members", "ShortUrl"}
	for _, header := range headers {
		cell = row.AddCell()
		cell.Value = header
	}

	for _, card := range cards {
		row = sheet.AddRow()
		url := fmt.Sprintf("%s/1/cards/%s?key=%s&token=%s&fields=all&attachments=false&attachment_fields=all&members=true&checkItemStates=true&checklists=none&checklist_fields=all&sticker_fields=all", trelloURL, card.ID, apiKey, clientToken)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal("Error occured: ", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
		}

		if resp.StatusCode != 200 {
			log.Fatal("bad request sent to trello api")
		}
		defer resp.Body.Close()
		var cardInfo cardDetails
		if err := json.NewDecoder(resp.Body).Decode(&cardInfo); err != nil {
			log.Fatal("Error decoding trello response: ", err)
		}
		cell = row.AddCell()
		//Add card name
		cell.Value = cardInfo.Name
		cell = row.AddCell()
		//Add card label
		var label string
		for _, labelInfo := range cardInfo.Labels {
			label += labelInfo.Name + ","
		}
		cell.Value = label
		cell = row.AddCell()
		//Add card description
		cell.Value = cardInfo.Desc
		//Add card Last used date
		cell = row.AddCell()
		cell.Value = cardInfo.DateLastActivity.String()
		//Add card closed status
		cell = row.AddCell()
		cell.Value = strconv.FormatBool(cardInfo.Closed)
		//Add card due date
		cell = row.AddCell()
		cell.Value = cardInfo.Due.String()
		//Add card due completed status
		cell = row.AddCell()
		cell.Value = strconv.FormatBool(cardInfo.DueComplete)
		//Add card no.of comments
		cell = row.AddCell()
		cell.Value = strconv.Itoa(cardInfo.Badges.Comments)
		//Add card members
		cell = row.AddCell()
		var member string
		for _, memberInfo := range cardInfo.Members {
			member += memberInfo.FullName
		}
		cell.Value = member
		//Add card short url
		cell = row.AddCell()
		cell.Value = cardInfo.ShortURL
	}
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	path := fmt.Sprintf("%s/Trello.xlsx", home)
	err = file.Save(path)
	if err != nil {
		log.Fatal("Could not save excelfile")
	}
	fmt.Printf("Please check your file in %s/Trello.xlsx\n", home)
}
