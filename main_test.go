package main

import (
	"bytes"
	"context"
	basev0 "github.com/codefly-dev/core/generated/go/base/v0"
	runtimev0 "github.com/codefly-dev/core/generated/go/services/runtime/v0"
	"github.com/codefly-dev/core/shared"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

func TestRunNoGoMod(t *testing.T) {
	location, err := shared.SolvePath("testdata/nogomod")
	assert.NoError(t, err)
	testRun(t, location)
}

// This doesn't work because of mixing go.mod...
//func TestRunWithGoMod(t *testing.T) {
//	location, err := shared.SolvePath("testdata/regular")
//	assert.NoError(t, err)
//	testRun(t, location)
//}

func testRun(t *testing.T, location string) {
	runner := NewRuntime()
	runner.Settings.Debug = false
	ctx := context.Background()
	_, err := runner.Load(ctx, &runtimev0.LoadRequest{
		Identity: &basev0.ServiceIdentity{
			Location: location,
		},
	})
	assert.NoError(t, err)
	_, err = runner.Init(ctx, &runtimev0.InitRequest{})
	assert.NoError(t, err)

	// We want to override the output
	var buf bytes.Buffer
	runner.WithOutput(&buf)
	_, err = runner.Start(ctx, &runtimev0.StartRequest{})
	assert.NoError(t, err)

	// Create a channel to signal when data is received
	dataReceived := make(chan bool)

	// Start a goroutine to listen for data on the buffer
	go func() {
		for {
			if buf.Len() > 0 {
				dataReceived <- true
				return
			}
			time.Sleep(100 * time.Millisecond) // Check every 100ms
		}
	}()

	defer os.RemoveAll(path.Join(location, ".cache"))

	// Wait for either data to be received or 5 seconds to pass
	select {
	case <-dataReceived:
		assert.Contains(t, buf.String(), "test")
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for data on buffer")
	}

}
