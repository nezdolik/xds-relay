package transport

import (
	"strconv"
	"sync"
	"testing"

	gcp "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetChannel(t *testing.T) {
	w := NewWatchV2(NewRequestV2(&gcp.Request{}))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, more := <-w.GetCh().(chan gcp.Response)
		assert.False(t, more)
		wg.Done()
	}()

	w.Close()
	wg.Wait()
}

func TestSendSuccessful(t *testing.T) {
	w := NewWatchV2(NewRequestV2(&gcp.Request{}))
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		got, more := <-w.GetCh().(chan gcp.Response)
		assert.True(t, more)
		assert.Equal(t, "1", got.Version)
		wg.Done()
	}()
	go func() {
		res := w.Req.CreateResponse("1", nil, nil)
		ok, err := w.Send(res)
		assert.True(t, ok)
		assert.Nil(t, err)
		wg.Done()
	}()
	wg.Wait()
}

func TestSendCastFailure(t *testing.T) {
	w := NewWatchV2(NewRequestV2(&gcp.Request{}))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ok, err := w.Send(&mockResponse{})
		assert.False(t, ok)
		assert.NotNil(t, err)
		wg.Done()
	}()
	wg.Wait()
}

func TestSendFalseWhenBlocked(t *testing.T) {
	w := NewWatchV2(NewRequestV2(&gcp.Request{}))
	var wg sync.WaitGroup
	resp := newResponseV2(w.Req, "1", nil, nil)
	wg.Add(2)
	// We perform 2 sends with no receive on w.Out .
	// One of the send gets blocked because of no recipient.
	// The second send goes goes to default case due to channel full.
	// The second send closes the channel when blocked.
	// The closed channel terminates the blocked send to exit the test case.
	go sendWithCloseChannelOnFailure(t, w, &wg, resp)
	go sendWithCloseChannelOnFailure(t, w, &wg, resp)

	wg.Wait()
}

func TestDropRedundant(t *testing.T) {
	receivedOne := make(chan bool)
	done := make(chan bool)
	w := NewWatchV2(NewRequestV2(&gcp.Request{}))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var i int64 = 1
		for {
			got, more := <-w.GetCh().(chan gcp.Response)
			if !more {
				return
			}
			receivedOne <- true
			assert.True(t, more)
			v, _ := strconv.ParseInt(got.Version, 10, 64)
			assert.Equal(t, v, i)
			i++
			<-done
		}
	}()
	go func() {
		// send the first response
		ok, err := w.Send(newResponseV2(w.Req, "1", nil, nil))
		assert.True(t, ok)
		assert.Nil(t, err)
		// since the channel length is 1, wait for the goroutine to receive the sent response
		<-receivedOne

		// There are no redundants at this point
		d := w.DropRedundant()
		assert.False(t, d)
		// send another response. Since the previous send is blocked on the done channel and channel
		// length is 1, send will go through and will be buffered.
		ok, err = w.Send(newResponseV2(w.Req, "5", nil, nil))
		assert.True(t, ok)
		assert.Nil(t, err)

		// Before sending the next response, drop the buffered response
		d = w.DropRedundant()
		assert.True(t, d)
		// send the next response
		ok, err = w.Send(newResponseV2(w.Req, "2", nil, nil))
		assert.True(t, ok)
		assert.Nil(t, err)
		// let the for select loop continue
		done <- true
		<-receivedOne

		close(done)
		w.Close()
		wg.Done()
	}()

	wg.Wait()
}

func sendWithCloseChannelOnFailure(t *testing.T, w Watch, wg *sync.WaitGroup, r Response) {
	ok, err := w.Send(r)
	assert.Nil(t, err)
	if !ok {
		w.Close()
	}
	wg.Done()
}

var _ Response = &mockResponse{}

type mockResponse struct {
}

// GetPayload gets the api version
func (r *mockResponse) GetPayload() interface{} {
	return ""
}
