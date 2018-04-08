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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type lists struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "fetch all lists of a board",
	Long:  `Usage: trello list "board name"`,
	Run: func(cmd *cobra.Command, args []string) {
		lists := getLists(args)
		fmt.Printf("%s\n------------------------\n", "Lists Available")
		for _, list := range lists {
			fmt.Printf("%s\n", list.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getLists(args []string) []lists {
	if len(args) == 0 {
		log.Fatalln("No board name specified")
	}
	requestedBoardName := args[0]
	availableBoards := getBoards()
	boardID := ""
	for _, boards := range availableBoards {
		if strings.ToLower(boards.Name) == strings.ToLower(requestedBoardName) {
			boardID = boards.ID
		}
	}
	apiKey = viper.Get("apiKey")
	clientToken = viper.Get("token")
	url := fmt.Sprintf("%s/1/boards/%s/lists?key=%s&token=%s", trelloURL, boardID, apiKey, clientToken)
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
	var trelloLists []lists
	if err := json.NewDecoder(resp.Body).Decode(&trelloLists); err != nil {
		log.Fatal("Error decoding trello response: ", err)
	}

	return trelloLists
}
