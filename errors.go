package webtransport

import (
	"errors"
	"fmt"

	"github.com/lucas-clemente/quic-go"
)

type ErrorCode uint8

const (
	firstErrorCode = 0x52e4a40fa8db
	lastErrorCode  = 0x52e4a40fa9e2
)

func webtransportCodeToHTTPCode(n ErrorCode) quic.StreamErrorCode {
	return quic.StreamErrorCode(firstErrorCode) + quic.StreamErrorCode(n) + quic.StreamErrorCode(n/0x1e)
}

func httpCodeToWebtransportCode(h quic.StreamErrorCode) (ErrorCode, error) {
	if h < firstErrorCode || h > lastErrorCode {
		return 0, errors.New("error code outside of expected range")
	}
	if (h-0x21)%0x1f == 0 {
		return 0, errors.New("invalid error code")
	}
	shifted := h - firstErrorCode
	return ErrorCode(shifted - shifted/0x1f), nil
}

// WebTransportBufferedStreamRejectedErrorCode is the error code of the
// H3_WEBTRANSPORT_BUFFERED_STREAM_REJECTED error.
const WebTransportBufferedStreamRejectedErrorCode quic.StreamErrorCode = 0x3994bd84

// StreamError is the error that is returned from stream operations (Read, Write) when the stream is canceled.
type StreamError struct {
	ErrorCode ErrorCode
}

func (e *StreamError) Is(target error) bool {
	_, ok := target.(*StreamError)
	return ok
}

func (e *StreamError) Error() string {
	return fmt.Sprintf("stream canceled with error code %d", e.ErrorCode)
}
