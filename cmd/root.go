/*
Copyright © 2022 Harry Law <hrryslw@pm.me>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/h5law/bush/walk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ignorePattern string
	levels        uint

	rootCmd = &cobra.Command{
		Use:   "bush",
		Short: "Recursively list the contents of a directory",
		Long: `Walk though the current directory or the those given as
arguments and display the contents recursively either in
a tree like structure or in plain text.

Use the help subcommand for info about the different flags
available to alter how this program runs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var dir []string
			if len(args) > 0 {
				dir = args
			} else {
				dir = []string{"."}
			}

			// Check if dirs provided are valid path
			for _, path := range dir {
				_, err := os.Stat(path)

				// If path is not valid print error and move to next
				if os.IsNotExist(err) {
					fmt.Printf("\"%s\" [error opening dir]\n", path)
					continue
				}
				if err != nil {
					return err
				}

				// Path is valid
				var dirCount, fileCount int
				if err := walk.Walk(path, &dirCount, &fileCount); err != nil {
					fmt.Println(err)
				}
				fmt.Printf("\n%d directories, %d files\n", dirCount, fileCount)
			}

			return nil
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().UintVarP(
		&levels,
		"levels",
		"L",
		0, "Levels to walk into directory",
	)

	rootCmd.Flags().StringVarP(
		&ignorePattern,
		"ignore",
		"I",
		"", "Pattern to ignore while walking directories",
	)

	viper.BindPFlag("levels", rootCmd.Flags().Lookup("levels"))
	viper.BindPFlag("ignore", rootCmd.Flags().Lookup("ignore"))
	viper.SetDefault("levels", 0)
	viper.SetDefault("ignore", "")
}
