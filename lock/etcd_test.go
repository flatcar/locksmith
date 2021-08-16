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
	"testing"

	pb "go.etcd.io/etcd/api/v3/mvccpb"
	client "go.etcd.io/etcd/client/v3"

	"golang.org/x/net/context"
)

const testGroup = ""

// testEtcdClient is the struct used to mock the `etcd`
// client
type testEtcdClient struct {
	err     error
	getResp *client.GetResponse
	txn     testTxn
}

// testTxn implements the client.Txn interface
type testTxn struct {
	err        error
	txnSuccess bool
}

func (t testTxn) If(cs ...client.Cmp) client.Txn {
	return t
}

func (t testTxn) Then(ops ...client.Op) client.Txn {
	return t
}

func (t testTxn) Else(ops ...client.Op) client.Txn {
	return t
}

func (t testTxn) Commit() (*client.TxnResponse, error) {
	return &client.TxnResponse{Succeeded: t.txnSuccess}, t.err
}

func (t *testEtcdClient) Get(ctx context.Context, key string, opts ...client.OpOption) (*client.GetResponse, error) {
	return t.getResp, t.err
}

func (t *testEtcdClient) Txn(ctx context.Context) client.Txn {
	return t.txn
}

func TestEtcdLockClientInit(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		elc, err := NewEtcdLockClient(&testEtcdClient{
			err:     nil,
			getResp: &client.GetResponse{},
		},
			testGroup,
		)
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		if elc == nil {
			t.Fatalf("etcd lock client should not be nil")
		}
	})
	t.Run("Error", func(t *testing.T) {
		_, err := NewEtcdLockClient(&testEtcdClient{
			txn:     testTxn{err: errors.New("connection refused")},
			getResp: &client.GetResponse{Count: 0},
		},
			testGroup,
		)
		if err == nil {
			t.Fatal("error should not be nil")
		}

		if err.Error() != "unable to init etcd lock client: unable to commit initial transaction: connection refused" {
			t.Fatalf("error should be 'unable to init etcd lock client: unable to commit initial transaction: connection refused', got: %v", err)
		}
	})
}

func TestEtcdLockClientGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		elc, err := NewEtcdLockClient(&testEtcdClient{
			err: nil,
			getResp: &client.GetResponse{
				Count: 1,
				Kvs: []*pb.KeyValue{
					&pb.KeyValue{
						Key: []byte(SemaphorePrefix),
						// index should be set from etcd, not json (backported from legacy test)
						Value:   []byte(`{"index": 12, "semaphore": 1, "max": 2, "holders": ["foo", "bar"]}`),
						Version: 1234,
					},
				},
			},
		},
			testGroup,
		)
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		res, err := elc.Get()
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		if res.Index != uint64(1234) {
			t.Fatalf("index should be 1234, got: %d", res.Index)
		}
	})
	t.Run("SuccessNotFound", func(t *testing.T) {
		elc, err := NewEtcdLockClient(&testEtcdClient{
			err:     nil,
			getResp: &client.GetResponse{Count: 0},
		},
			testGroup,
		)
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		_, err = elc.Get()
		if err == nil {
			t.Fatal("error should not be nil")
		}

		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("error should be ErrNotFound, got: %v", err)
		}
	})
	t.Run("ErrorWithMalformedJSON", func(t *testing.T) {
		elc, err := NewEtcdLockClient(&testEtcdClient{
			err: nil,
			getResp: &client.GetResponse{
				Count: 1,
				Kvs: []*pb.KeyValue{
					&pb.KeyValue{
						Key: []byte(SemaphorePrefix),
						// notice the missing `,` in the array
						Value: []byte(`{"semaphore": 1, "max": 2, "holders": ["foo" "bar"]}`),
					},
				},
			},
		},
			testGroup,
		)
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		_, err = elc.Get()
		if err == nil {
			t.Fatal("error should not be nil")
		}

		if err.Error() != "invalid character '\"' after array element" {
			t.Fatalf("error should mention invalid character, got: %v", err)
		}
	})
}

func TestEtcdLockClientSet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		elc, err := NewEtcdLockClient(&testEtcdClient{
			err:     nil,
			getResp: &client.GetResponse{Count: 0},
			txn:     testTxn{txnSuccess: true},
		},
			testGroup,
		)
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		err = elc.Set(&Semaphore{})
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

	})
	t.Run("ErrorNilSemaphore", func(t *testing.T) {
		elc, err := NewEtcdLockClient(&testEtcdClient{
			err:     nil,
			getResp: &client.GetResponse{Count: 0},
		},
			testGroup,
		)
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		err = elc.Set(nil)
		if err == nil {
			t.Fatal("error should not be nil")
		}

		if err.Error() != "cannot set nil semaphore" {
			t.Fatalf("error should 'cannot set nil semaphore', got: %v", err)
		}
	})
	t.Run("ErrorTransaction", func(t *testing.T) {
		elc, err := NewEtcdLockClient(&testEtcdClient{
			err:     nil,
			getResp: &client.GetResponse{Count: 0},
		},
			testGroup,
		)
		if err != nil {
			t.Fatalf("error should be nil, got: %v", err)
		}

		err = elc.Set(&Semaphore{})
		if err == nil {
			t.Fatal("error should not be nil")
		}

		if err.Error() != "failed to set the semaphore - it got updated in the meantime" {
			t.Fatalf("error should be 'failed to set the semaphore - it got updated in the meantime', got: %v", err)
		}
	})
}
