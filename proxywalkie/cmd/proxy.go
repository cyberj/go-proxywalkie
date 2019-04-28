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
	"github.com/cyberj/go-proxywalkie/proxy"
	"github.com/spf13/cobra"
)

var proxyDelete bool
var proxyBackground bool
var proxySyncInterval int
var proxyServer string
var proxyPort string

// lsCmd represents the ls command
var proxyCmd = &cobra.Command{
	Use:   "proxy SERVER",
	Short: "Proxy (client) for Intuiface. You need the full httpurl for server (like \"http://1.2.3.4:8080/\")",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		starttime := time.Now()

		proxyServer = args[0]

		logrus.Infof("Initializing Proxy")
		proxy, err := proxy.NewProxyParams(workdirPath, proxyServer, time.Duration(proxySyncInterval)*time.Minute, proxyDelete, proxyBackground)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Infof("Proxy initialized (%s)", time.Since(starttime))
		proxy.Ready()
		logrus.Infof("Proxy ready (%s)", time.Since(starttime))

		err = http.ListenAndServe(":"+proxyPort, proxy.Router())
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	proxyCmd.PersistentFlags().StringVarP(&proxyPort, "port", "p", "8081", "Local server URL")
	proxyCmd.PersistentFlags().BoolVarP(&proxyDelete, "delete", "d", false, "Delete files")
	proxyCmd.PersistentFlags().BoolVarP(&proxyBackground, "background", "b", false, "Background Sync")
	proxyCmd.PersistentFlags().IntVarP(&proxySyncInterval, "sync-interval", "u", 5, "Sync interval (in minutes)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
