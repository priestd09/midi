package realtime

import (
	"fmt"
	"io"
)

// every realtime.Reader is an io.Reader but not every io.Reader is a realtime.Reader
type Reader interface {
	io.Reader
	realtime()
}

/*
    Each RealTime Category message (ie, Status of 0xF8 to 0xFF) consists of only 1 byte, the Status. These messages are primarily concerned with timing/syncing functions which means that they must be sent and received at specific times without any delays. Because of this, MIDI allows a RealTime message to be sent at any time, even interspersed within some other MIDI message. For example, a RealTime message could be sent inbetween the two data bytes of a Note On message. A device should always be prepared to handle such a situation; processing the 1 byte RealTime message, and then subsequently resume processing the previously interrupted message as if the RealTime message had never occurred.

For more information about RealTime, read the sections Running Status, Ignoring MIDI Messages, and Syncing Sequence Playback.
*/

// reader is is a wrapper around an io.Reader that filters realtime midi events
// when reading it calls Callback for every realtime event and reading everything else into the target buffer
type reader struct {
	input   io.Reader
	handler func(Message)
}

func (r *reader) realtime() {}

// discardReader is an optimized reader that discards realtime messages
type discardReader struct {
	input io.Reader
}

func (r *discardReader) realtime() {}

func (r *discardReader) Read(target []byte) (n int, err error) {
	var bf []byte
	var one int

	for {
		if n == len(target) {
			return

		}
		bf = make([]byte, 1)

		one, err = r.input.Read(bf)

		if err != nil {
			return
		}

		if one != 1 {
			err = fmt.Errorf("could not read %v byte(s)", len(target))
			return
		}

		// => no realtime message
		if bf[0] < 0xF8 {
			target[n] = bf[0]
			n++
			continue
		}

		// don't handle realtime messages, so do nothing here

	}
	return
}

// NewReader returns an io.Reader that filters realtime midi messages.
// For each realtime midi message, rthandler is called (if it is not nil)
// The Reader does no buffering and makes no attempt to close input.
func NewReader(input io.Reader, rthandler func(Message)) Reader {
	if rthandler == nil {
		return &discardReader{input}
	}
	return &reader{input, rthandler}
}

func (r *reader) Read(target []byte) (n int, err error) {
	var bf []byte
	var one int

	for {
		if n == len(target) {
			return

		}
		bf = make([]byte, 1)

		one, err = r.input.Read(bf)

		if err != nil {
			return
		}

		if one != 1 {
			err = fmt.Errorf("could not read %v byte(s)", len(target))
			return
		}

		// => no realtime message
		if bf[0] < 0xF8 {
			target[n] = bf[0]
			n++
			continue
		}

		// error needed here to be able to interrupt the reading from the callback (handler)
		// then an io.EOF error is returned and propagated to midireader.read()
		ev := Dispatch(bf[0])

		// we know that r.handler is not nil (otherwise we would be inside discardReader)
		if ev != nil {
			r.handler(ev)
		}

	}
	return
}