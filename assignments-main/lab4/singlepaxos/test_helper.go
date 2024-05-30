// Test helper functions - DO NOT EDIT

package singlepaxos

const (
	valueFromClientOne = "A client command"
	valueFromClientTwo = "Another client command"
)

type msgPair[Req any, Resp any] struct {
	req      Req
	wantResp Resp
}

type (
	promiseAccept  msgPair[Promise, Accept]
	preparePromise msgPair[Prepare, Promise]
	acceptLearn    msgPair[Accept, Learn]
	learnValue     msgPair[Learn, Value]
)
