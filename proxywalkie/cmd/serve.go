// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	server "github.com/cyberj/go-proxywalkie/server"
	"github.com/spf13/cobra"
)

var serverSyncInterval int

// lsCmd represents the ls command
var serveCmd = &cobra.Command{
	Use:   "serve [:8080]",
	Short: "Serve a path for improved proxy (server)",
	Run: func(cmd *cobra.Command, args []string) {

		starttime := time.Now()
		proxy, err := server.NewServerParams(workdirPath, time.Duration(serverSyncInterval)*time.Minute)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Infof("Proxy ready (init: %s)", time.Since(starttime))

		listenTo := ":8080"
		if len(args) > 0 {
			listenTo = args[0]
		}

		err = http.ListenAndServe(listenTo, proxy)
		if err != nil {
			logrus.Fatal(err)
			// os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	rootCmd.Flags().IntVarP(&serverSyncInterval, "refresh-interval", "r", 10, "Refresh interval (in minutes)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
