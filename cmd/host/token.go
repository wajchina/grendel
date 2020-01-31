// Copyright 2019 Grendel Authors. All rights reserved.
//
// This file is part of Grendel.
//
// Grendel is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Grendel is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Grendel. If not, see <https://www.gnu.org/licenses/>.

package host

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ubccr/grendel/client"
	"github.com/ubccr/grendel/model"
	"github.com/ubccr/grendel/nodeset"
)

var (
	tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "Generate boot token for hosts",
		Long:  `Generate boot token for hosts`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			gc, err := client.NewClient()
			if err != nil {
				return err
			}

			ns, err := nodeset.NewNodeSet(strings.Join(args, ","))
			if err != nil {
				return err
			}

			if ns.Len() == 0 {
				return errors.New("Node nodes in nodeset")
			}

			hostList, err := gc.FindHosts(ns)
			if err != nil {
				return err
			}

			for _, host := range hostList {
				if len(host.Interfaces) == 0 {
					continue
				}

				token, err := model.NewBootToken(host.ID.String(), host.Interfaces[0].MAC.String())
				if err != nil {
					return fmt.Errorf("Failed to generate signed boot token for host %s: %s", host.Name, err)
				}

				fmt.Printf("%s: %s\n", host.Name, token)
			}

			return nil

		},
	}
)

func init() {
	hostCmd.AddCommand(tokenCmd)
}
