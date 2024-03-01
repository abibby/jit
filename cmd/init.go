/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abibby/jit/git"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Add hooks to your repo",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		//  applypatch-msg
		//  commit-msg
		//  fsmonitor-watchman
		//  post-update
		//  pre-applypatch
		//  pre-commit
		//  pre-merge-commit
		//  pre-push
		//  pre-rebase
		//  pre-receive
		//  prepare-commit-msg
		//  push-to-checkout
		//  update
		m := map[string]string{
			"prepare-commit-msg": prepareCommitMsgCmd.Use,
		}
		root, err := git.Root()
		if err != nil {
			return err
		}
		ex, err := os.Executable()
		if err != nil {
			return err
		}
		for hook, command := range m {
			file := filepath.Join(root, "hooks", hook)

			err = os.Remove(file)
			if errors.Is(err, os.ErrNotExist) {
			} else if err != nil {
				fmt.Printf("%#v", err)
				return err
			}

			return os.WriteFile(file, []byte(fmt.Sprintf("%s %s \"$@\"\n", ex, command)), 0777)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
