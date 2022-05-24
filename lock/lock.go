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

package lock

import (
	"errors"
)

// Lock takes care of locking in generic clients
type Lock struct {
	id     string
	client LockClient
}

var (
	// ErrExist is the error returned if a holder with the specified id
	// is already holding the semaphore
	ErrExist = errors.New("holder exists")
	// ErrNotExist is the error returned if there is no holder with the
	// specified id holding the semaphore
	ErrNotExist = errors.New("holder does not exist")
)

// New returns a new lock with the provided arguments
func New(id string, client LockClient) (lock *Lock) {
	return &Lock{id, client}
}

// Lock adds this lock id as a holder to the semaphore
// it will return an error if there is a problem getting or setting the
// semaphore, or if the maximum number of holders has been reached, or if a lock
// with this id is already a holder
func (l *Lock) Lock() (err error) {
	return l.client.RecursiveLock()
}

// Unlock removes this lock id as a holder of the semaphore
// it returns an error if there is a problem getting or setting the semaphore,
// or if this lock is not locked.
func (l *Lock) Unlock() error {
	return l.client.UnlockIfHeld()
}
