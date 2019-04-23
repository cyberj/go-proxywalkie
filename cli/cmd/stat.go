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
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cyberj/go-proxywalkie/walkie"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Give informations about a directory",
	Run: func(cmd *cobra.Command, args []string) {

		starttime := time.Now()

		w, err := walkie.NewWalkie(workdirPath)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Infof("Walkie init : %s", time.Since(starttime))
		intermediate := time.Now()

		err = w.Explore()
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Infof("Walkie Explore: %s", time.Since(starttime))
		intermediate = time.Now()

		nbdir, nbfiles := w.Directory.Stat()
		logrus.Infof("Directories : %v", nbdir)
		logrus.Infof("Files : %v", nbfiles)

		// fmt.Printf("Directories: %d\n", unsafe.Sizeof(w.Directory))

		logrus.Infof("Stat : %s", time.Since(intermediate))

		// endtime := time.Now()

		fmt.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(statCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
