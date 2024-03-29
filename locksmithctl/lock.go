// Copyright 2015 CoreOS, Inc.
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

package main

import (
	"fmt"
	"os"

	"github.com/flatcar/locksmith/lock"
	"github.com/flatcar/locksmith/pkg/machineid"
)

var (
	cmdLock = &Command{
		Name:    "lock",
		Summary: "Lock this machine or a given machine-id for reboot.",
		Usage:   "<machine-id>",
		Description: `Lock is for manual locking of the reboot lock for this machine or a given
machine-id. Under normal operation this should not be necessary.`,
		Run: runLock,
	}
)

func runLock(args []string) (exit int) {
	elc, err := getClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing etcd client:", err)
		return 1
	}

	var mID string

	if len(args) == 0 {
		mID = machineid.MachineID("/")
		if mID == "" {
			fmt.Fprintln(os.Stderr, "Cannot read machine-id")
			return 1
		}
	} else {
		mID = args[0]
	}

	l := lock.New(mID, elc)

	err = l.Lock()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error locking:", err)
		return 1
	}

	return 0
}
