package jsonbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"ntoolkit/errors"
	"ntoolkit/linereader"
	"time"
)

// Constants
const defaultBufferSize = 1024
const defaultTimeout = 100

// Bridge is a basic network json reader.
type Bridge struct {

	// Chunk size is the size of data read from the bridge each read.
	ChunkSize int

	// Timeout is the timeout for read and write operations in ms
	Timeout int

	// The line reader to read data with
	reader *linereader.LineReader

	// Read buffer
	buffer *bytes.Buffer

	// Current data
	data *bytes.Buffer

	// Network conenctions
	input  net.Conn
	output net.Conn
}

// New returns a new Bridge that looks at the given network
// connections input and output for communication.
func New(input net.Conn, output net.Conn) *Bridge {
	return &Bridge{
		defaultBufferSize,
		defaultTimeout,
		linereader.New(),
		bytes.NewBuffer(make([]byte, defaultBufferSize, defaultBufferSize)),
		bytes.NewBuffer(make([]byte, 1, 1)),
		input,
		output}
}

// Read chunks into the internal buffer until we hit timeout or maximum read interval.
// If reading fails due to an error, the pending buffer is loaded.
func (bridge *Bridge) Read() error {
	bridge.input.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(bridge.Timeout)))
	count, err := bridge.input.Read(bridge.buffer.Bytes())
	if err != nil {
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				return nil // Timeouts mean no pending messages, not an error
			}
			return errors.Fail(ErrRead{}, err, "Failed to read from stream")
		}
	} else {
		fmt.Printf("Read: %v\n", string(bridge.buffer.Bytes()[:count]))
		bridge.reader.Write(bridge.buffer.Bytes()[:count])
	}
	return nil
}

// Write an object as json to the bridge
func (bridge *Bridge) Write(data interface{}) error {

	// To string
	jstr, err := json.Marshal(data)
	if err != nil {
		return errors.Fail(ErrMarshal{}, err, "Failed to convert object to data")
	}

	// Setup
	bridge.output.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(bridge.Timeout)))
	jdata := []byte(jstr)

	// Push to stream
	sent := 0
	total := len(jdata)
	for sent < total {
		count, werr := bridge.output.Write(jdata[sent:])
		if werr != nil {
			return errors.Fail(ErrWrite{}, werr, "Failed to write to stream")
		}
		sent += count
	}

	// End of data marker
	_, err = bridge.output.Write([]byte("\n"))
	if err != nil {
		return errors.Fail(ErrWrite{}, err, "Failed to write to stream")
	}

	fmt.Printf("Wrote all bytes to stream!\n")

	return nil
}

// Len returns the number of pending messages on the bridge
func (bridge *Bridge) Len() int {
	return bridge.reader.Len()
}

// Next updates the internal buffer to point at the next object
func (bridge *Bridge) Next() error {
	if bridge.reader.Len() > 0 {
		bridge.data.Reset()
		bridge.data.Write([]byte(bridge.reader.Next()))
		return nil
	}
	return errors.Fail(ErrNoData{}, nil, "No objects available")
}

// As attempts to parse the current buffered object into the given data object.
func (bridge *Bridge) As(data interface{}) error {
	if err := json.Unmarshal(bridge.data.Bytes(), data); err != nil {
		return errors.Fail(ErrUnmarshal{}, err, "Failed to convert data to object")
	}
	return nil
}

// Raw returns the current active chunk as a string
func (bridge *Bridge) Raw() string {
	return string(bridge.data.Bytes())
}
