package main

import (
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	pb "dat520/lab2/grpc/proto"
)

func TestLookup(t *testing.T) {
	for _, test := range lookupTests {
		t.Run(test.desc, func(t *testing.T) {
			for index := range test.insertReq {
				gotInsertResp, err := test.server.Insert(context.Background(), test.insertReq[index])
				if err != nil {
					t.Fatal(err)
				}
				if gotInsertResp.Success != test.insertResp[index].Success {
					k, v := test.insertReq[index].Key, test.insertReq[index].Value
					t.Errorf("Insert(%s, %s) = %t, want %t", k, v, test.insertResp[index].Success, gotInsertResp.Success)
				}
			}

			for index := range test.lookupReq {
				gotLookupResp, err := test.server.Lookup(context.Background(), test.lookupReq[index])
				if err != nil {
					t.Fatal(err)
				}
				if gotLookupResp.Value != test.lookupResp[index].Value {
					k := test.lookupReq[index].Key
					t.Errorf("Lookup(%s) = %s, want %s", k, gotLookupResp.Value, test.lookupResp[index].Value)
				}
			}
		})
	}
}

func TestKeys(t *testing.T) {
	for _, test := range keysTests {
		t.Run(test.desc, func(t *testing.T) {
			for index := range test.insertReq {
				gotInsertResp, err := test.server.Insert(context.Background(), test.insertReq[index])
				if err != nil {
					t.Fatal(err)
				}
				if gotInsertResp.Success != test.insertResp[index].Success {
					k, v := test.insertReq[index].Key, test.insertReq[index].Value
					t.Errorf("Insert(%s, %s) = %t, want %t", k, v, test.insertResp[index].Success, gotInsertResp.Success)
				}
			}

			for index := range test.keysReq {
				gotKeysResp, err := test.server.Keys(context.Background(), test.keysReq[index])
				if err != nil {
					t.Fatal(err)
				}
				// Sort the keys in alphabetical order before comparing since the order of the keys is not specified.
				sort.Strings(gotKeysResp.Keys)
				if !cmp.Equal(gotKeysResp.Keys, test.keysResp[index].Keys) {
					t.Errorf("Keys() = %v, want %v", gotKeysResp.Keys, test.keysResp[index].Keys)
				}
			}
		})
	}
}

var lookupTests = []struct {
	server     *keyValueServicesServer
	desc       string
	insertReq  []*pb.InsertRequest
	insertResp []*pb.InsertResponse
	lookupReq  []*pb.LookupRequest
	lookupResp []*pb.LookupResponse
}{
	{
		NewKeyValueServicesServer(),
		"Lookup After Zero Key/Values Inserted",
		nil,
		nil,
		[]*pb.LookupRequest{{Key: "1"}},
		[]*pb.LookupResponse{{Value: ""}},
	},
	{
		NewKeyValueServicesServer(),
		"Lookup After One Key/Value Inserted",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}},
		[]*pb.InsertResponse{{Success: true}},
		[]*pb.LookupRequest{{Key: "1"}},
		[]*pb.LookupResponse{{Value: "one"}},
	},
	{
		NewKeyValueServicesServer(),
		"Lookup After One Key/Value Inserted: Empty Value",
		[]*pb.InsertRequest{{Key: "1", Value: ""}},
		[]*pb.InsertResponse{{Success: true}},
		[]*pb.LookupRequest{{Key: "1"}},
		[]*pb.LookupResponse{{Value: ""}},
	},
	{
		NewKeyValueServicesServer(),
		"Lookup After One Key/Value Inserted: Empty Key",
		[]*pb.InsertRequest{{Key: "", Value: "one"}},
		[]*pb.InsertResponse{{Success: true}},
		[]*pb.LookupRequest{{Key: ""}},
		[]*pb.LookupResponse{{Value: "one"}},
	},
	{
		NewKeyValueServicesServer(),
		"Lookup After Two Key/Values Inserted",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}, {Key: "2", Value: "two"}},
		[]*pb.InsertResponse{{Success: true}, {Success: true}},
		[]*pb.LookupRequest{{Key: "2"}},
		[]*pb.LookupResponse{{Value: "two"}},
	},
	{
		NewKeyValueServicesServer(),
		"Lookup After Two Key/Values Inserted: Same key, different values",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}, {Key: "1", Value: "one again"}},
		[]*pb.InsertResponse{{Success: true}, {Success: true}},
		[]*pb.LookupRequest{{Key: "1"}},
		[]*pb.LookupResponse{{Value: "one again"}},
	},
	{
		NewKeyValueServicesServer(),
		"Lookup After Three Key/Values Inserted",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}, {Key: "2", Value: "two"}, {Key: "3", Value: "three"}},
		[]*pb.InsertResponse{{Success: true}, {Success: true}, {Success: true}},
		[]*pb.LookupRequest{{Key: "3"}},
		[]*pb.LookupResponse{{Value: "three"}},
	},
}

var keysTests = []struct {
	server     *keyValueServicesServer
	desc       string
	insertReq  []*pb.InsertRequest
	insertResp []*pb.InsertResponse
	keysReq    []*pb.KeysRequest
	keysResp   []*pb.KeysResponse
}{
	{
		NewKeyValueServicesServer(),
		"Keys After No Key/Values Inserted",
		nil,
		nil,
		[]*pb.KeysRequest{{}},
		[]*pb.KeysResponse{{Keys: []string{}}},
	},
	{
		NewKeyValueServicesServer(),
		"Keys After One Key/Value Inserted",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}},
		[]*pb.InsertResponse{{Success: true}},
		[]*pb.KeysRequest{{}},
		[]*pb.KeysResponse{{Keys: []string{"1"}}},
	},
	{
		NewKeyValueServicesServer(),
		"Keys After Two Key/Values Inserted",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}, {Key: "2", Value: "two"}},
		[]*pb.InsertResponse{{Success: true}, {Success: true}},
		[]*pb.KeysRequest{{}},
		[]*pb.KeysResponse{{Keys: []string{"1", "2"}}},
	},
	{
		NewKeyValueServicesServer(),
		"Keys After Two Key/Values Inserted: Same Key Twice",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}, {Key: "1", Value: "one again"}},
		[]*pb.InsertResponse{{Success: true}, {Success: true}},
		[]*pb.KeysRequest{{}},
		[]*pb.KeysResponse{{Keys: []string{"1"}}},
	},
	{
		NewKeyValueServicesServer(),
		"Keys After Three Key/Values Inserted",
		[]*pb.InsertRequest{{Key: "1", Value: "one"}, {Key: "2", Value: "two"}, {Key: "3", Value: "three"}},
		[]*pb.InsertResponse{{Success: true}, {Success: true}, {Success: true}},
		[]*pb.KeysRequest{{}},
		[]*pb.KeysResponse{{Keys: []string{"1", "2", "3"}}},
	},
}
