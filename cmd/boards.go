// Copyright © 2018 NAME HERE athavankanapuli@gmail.com
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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// boards is the response for trello boards api
type boards struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// boardCmd represents the board command
var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "trello board helps to list the boards and their board ids",
	Long:  `trello board helps to list the boards and their board ids`,
	Run: func(cmd *cobra.Command, args []string) {
		//trelloURL = "https://api.trello.com"
		apiKey = viper.Get("apiKey")
		clientToken = viper.Get("token")
		url := fmt.Sprintf("%s/1/members/me/boards?key=%s&token=%s", trelloURL, apiKey, clientToken)

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
		var trelloBoards []boards
		if err := json.NewDecoder(resp.Body).Decode(&trelloBoards); err != nil {
			log.Fatal("Error decoding trello response: ", err)
		}
		fmt.Printf("%s\n------------------------\n", "List of Boards")
		for _, board := range trelloBoards {
			fmt.Printf("%s\n", board.Name)
		}

	},
}

func init() {
	rootCmd.AddCommand(boardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// boardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// boardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
