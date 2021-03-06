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
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

// getUpdatesCmd represents the updates command
var getUpdatesCmd = &cobra.Command{
	Use:   "updates",
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

		since, err := cmd.Flags().GetString("since")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		latestReleasesByRepo := map[*github.Repository]*github.RepositoryRelease{}
		if since != "" {
			latestReleasesByRepo, err = gh.ListStarredReposLatestReleasesSince(token, since)
		} else {
			latestReleasesByRepo, err = gh.ListStarredReposLatestReleases(token)
		}
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		for repo, release := range latestReleasesByRepo {
			if release != nil {
				fmt.Printf("(%s) %s -- (%s) %s\n", release.GetPublishedAt().Local().Format("Mon Jan 2 3:04 PM"), repo.GetFullName(), release.GetTagName(), release.GetName())
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getUpdatesCmd)
	getUpdatesCmd.PersistentFlags().String("since", "", "Show updates since a particular date/time (uses systemd.time format)")
}
