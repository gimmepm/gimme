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
	"fmt"
	"os"

	"github.com/gimmepm/gimme/pkg/gh"
	"github.com/spf13/cobra"
)

// getReposCmd represents the repos command
var getReposCmd = &cobra.Command{
	Use:   "repos",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := GetToken(cmd)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		starredRepos, err := gh.ListStarredRepos(token)
		if err != nil {
			fmt.Printf("getReposCmd.Run: Error listing repos: %v\n", err)
			os.Exit(1)
		}

		for _, repo := range starredRepos {
			fmt.Printf(
				"%s/%s\n",
				repo.GetRepository().GetOwner().GetLogin(),
				repo.GetRepository().GetName(),
			)
		}
	},
}

func init() {
	getCmd.AddCommand(getReposCmd)
}
