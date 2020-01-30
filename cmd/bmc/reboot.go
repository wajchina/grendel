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

package bmc

import (
	"errors"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ubccr/grendel/client"
	"github.com/ubccr/grendel/nodeset"
)

var (
	rebootCmd = &cobra.Command{
		Use:   "reboot",
		Short: "Reboot hosts",
		Long:  `Reboot hosts`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			ns, err := nodeset.NewNodeSet(strings.Join(args, ","))
			if err != nil {
				return err
			}

			if ns.Len() == 0 {
				return errors.New("Node nodes in nodeset")
			}

			return runReboot(ns)
		},
	}
)

func init() {
	bmcCmd.AddCommand(rebootCmd)
}

func runReboot(ns *nodeset.NodeSet) error {
	gc, err := client.NewClient()
	if err != nil {
		return err
	}

	hostList, err := gc.FindHosts(ns)
	if err != nil {
		return err
	}

	if len(hostList) == 0 {
		return errors.New("No hosts found")
	}

	delay := viper.GetInt("bmc.delay")
	runner := NewJobRunner(viper.GetInt("bmc.fanout"))
	for _, host := range hostList {
		runner.RunReboot(host)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	runner.Wait()

	return nil
}
