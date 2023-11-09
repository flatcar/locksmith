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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"

	client "go.etcd.io/etcd/client/v3"

	"golang.org/x/net/context"
)

// ErrNotFound is used when a key is not found - which means
// it returns 0 value.
var ErrNotFound = errors.New("key not found")

const (
	keyPrefix       = "coreos.com/updateengine/rebootlock"
	groupBranch     = "groups"
	semaphoreBranch = "semaphore"
	// SemaphorePrefix is the key in etcd where the semaphore will be stored
	SemaphorePrefix = keyPrefix + "/" + semaphoreBranch
)

// KeysAPI is the minimum etcd client.KV interface EtcdLockClient needs
// to do its job.
type KeysAPI interface {
	Get(ctx context.Context, key string, opts ...client.OpOption) (*client.GetResponse, error)
	Txn(ctx context.Context) client.Txn
}

// EtcdLockClient is a wrapper around the etcd client that provides
// simple primitives to operate on the internal semaphore and holders
// structs through etcd.
type EtcdLockClient struct {
	keyapi  KeysAPI
	keypath string
}

// NewEtcdLockClient creates a new EtcdLockClient. The group parameter defines
// the etcd key path in which the client will manipulate the semaphore. If the
// group is the empty string, the default semaphore will be used.
func NewEtcdLockClient(keyapi KeysAPI, group string) (*EtcdLockClient, error) {
	key := SemaphorePrefix
	if group != "" {
		key = path.Join(keyPrefix, groupBranch, url.QueryEscape(group), semaphoreBranch)
	}

	elc := &EtcdLockClient{keyapi, key}
	if err := elc.Init(); err != nil {
		return nil, fmt.Errorf("unable to init etcd lock client: %w", err)
	}

	return elc, nil
}

// Init sets an initial copy of the semaphore if it doesn't exist yet.
// So we first try to get the value, if the value is not found we create the key
// with a default semaphore value.
func (c *EtcdLockClient) Init() error {
	sem := newSemaphore()
	payload, err := json.Marshal(sem)
	if err != nil {
		return fmt.Errorf("unable to marshal initial semaphore: %w", err)
	}

	if _, err := c.keyapi.Txn(context.TODO()).
		If(
			client.Compare(client.Version(c.keypath), "=", 0),
		).
		Then(
			client.OpPut(c.keypath, string(payload)),
		).
		Commit(); err != nil {
		return fmt.Errorf("unable to commit initial transaction: %w", err)
	}

	return nil
}

// Get fetches the Semaphore from etcd.
func (c *EtcdLockClient) Get() (*Semaphore, error) {
	resp, err := c.keyapi.Get(context.Background(), c.keypath, client.WithLastCreate()...)
	if err != nil {
		return nil, err
	}

	// There is no proper way to handle non-existing value for a
	// given key.
	// See https://github.com/etcd-io/etcd/issues/6089 for more details.
	if resp.Count == 0 {
		return nil, ErrNotFound
	}

	kv := resp.Kvs[0]

	sem := &Semaphore{}
	err = json.Unmarshal(kv.Value, sem)
	if err != nil {
		return nil, err
	}

	sem.Index = uint64(kv.Version)

	return sem, nil
}

// Set sets a Semaphore in etcd.
func (c *EtcdLockClient) Set(sem *Semaphore) error {
	if sem == nil {
		return errors.New("cannot set nil semaphore")
	}
	b, err := json.Marshal(sem)
	if err != nil {
		return err
	}

	response, err := c.keyapi.Txn(context.Background()).
		If(
			client.Compare(client.Version(c.keypath), "=", int64(sem.Index)),
		).
		Then(
			client.OpPut(c.keypath, string(b)),
		).
		Commit()
	if err != nil {
		return fmt.Errorf("making transaction: %w", err)
	}

	if !response.Succeeded {
		return fmt.Errorf("failed to set the semaphore - it got updated in the meantime")
	}

	return nil
}
