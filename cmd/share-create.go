// Copyright © 2016 Dropbox, Inc.
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
	"errors"
	"fmt"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/sharing"
	"github.com/spf13/cobra"
)

func shareCreate(cmd *cobra.Command, args []string) (err error) {
	path := ""
	if len(args) == 1 {
		if path, err = validatePath(args[0]); err != nil {
			return err
		}
	} else {
			return errors.New("share create: missing operand")
	}

	linkAge, _ := time.ParseDuration("240h")

	arg := sharing.NewCreateSharedLinkWithSettingsArg(path)

	settings := sharing.NewSharedLinkSettings()
	settings.Expires = time.Now().UTC().Round(time.Second).Add(linkAge)

	arg.Settings = settings	

	dbx := sharing.New(config)
	res, err := dbx.CreateSharedLinkWithSettings(arg)
	if err != nil {
		switch e := err.(type) {
		case sharing.CreateSharedLinkWithSettingsAPIError:
			fmt.Printf("%v", e.EndpointError)
		default:
			return err
		}
	}

	prepareSharedLink(res)

	return
}

func prepareSharedLink(link sharing.IsSharedLinkMetadata) {
		switch sl := link.(type) {
		case *sharing.FileLinkMetadata:
			printLink(sl.SharedLinkMetadata)
		case *sharing.FolderLinkMetadata:
			printLink(sl.SharedLinkMetadata)
		default:
			fmt.Printf("found unknown shared link type")
	}
}

func printShareLink(sl sharing.SharedLinkMetadata) {
	fmt.Printf("%v\t%v\n", sl.Name, sl.Url)
}

var shareCreateCmd = &cobra.Command{
	Use:   "create <path>",
	Short: "Create shared link",
	RunE:  shareCreate,
}

func init() {
	shareCmd.AddCommand(shareCreateCmd)
}
