package syscommon

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMessagesString(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			MIDITimingCode(3),
			"syscommon.MIDITimingCode: 3",
		},
		{
			SongPositionPointer(4),
			"syscommon.SongPositionPointer: 4",
		},
		{
			SongSelect(2),
			"syscommon.SongSelect: 2",
		},
		{
			TuneRequest,
			"syscommon.tuneRequest",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.WriteString(test.input.String())

		if got, want := bf.String(), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}

func TestMessagesSongPositionPointer(t *testing.T) {

	tests := []struct {
		expected uint16
	}{
		{8},
		{4},
		{32},
		{320},
		{13320},
	}

	for _, test := range tests {
		bt := SongPositionPointer(test.expected).Raw()
		rd := bytes.NewReader(bt)
		typ, _ := rd.ReadByte()
		r := NewReader(rd, typ)
		msg, _ := r.Read()

		if got, want := msg.(SongPositionPointer).Number(), test.expected; got != want {
			t.Errorf("SongPositionPointer(%v) = %v; want %v", test.expected, got, want)
		}
	}

}

func TestMessagesRaw(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			MIDITimingCode(3),
			"F1 03",
		},
		{
			SongPositionPointer(8),
			"F2 00 08",
		},
		{
			SongSelect(2),
			"F3 02",
		},
		{
			TuneRequest,
			"F6",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.Write(test.input.Raw())

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}
