package prometheus

import (
	"io"
	"math"
	"strconv"
	"sync"
)

//taken from the original prometheus go client

const (
	initialNumBufSize = 24
)

var (
	numBufPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, initialNumBufSize)
			return &b
		},
	}
)

// writeFloat is equivalent to fmt.Fprint with a float64 argument but hardcodes
// a few common cases for increased efficiency. For non-hardcoded cases, it uses
// strconv.AppendFloat to avoid allocations, similar to writeInt.
func writeFloat(w io.Writer, f float64) (int, error) {
	switch {
	case f == 1:
		return w.Write([]byte("1"))
	case f == 0:
		return w.Write([]byte("0"))
	case f == -1:
		return w.Write([]byte("-1"))
	case math.IsNaN(f):
		return w.Write([]byte("NaN"))
	case math.IsInf(f, +1):
		return w.Write([]byte("+Inf"))
	case math.IsInf(f, -1):
		return w.Write([]byte("-Inf"))
	default:
		bp := numBufPool.Get().(*[]byte)
		*bp = strconv.AppendFloat((*bp)[:0], f, 'g', -1, 64)
		written, err := w.Write(*bp)
		numBufPool.Put(bp)
		return written, err
	}
}

// writeInt is equivalent to fmt.Fprint with an int64 argument but uses
// strconv.AppendInt with a byte slice taken from a sync.Pool to avoid
// allocations.
func writeInt(w io.Writer, i int64) (int, error) {
	bp := numBufPool.Get().(*[]byte)
	*bp = strconv.AppendInt((*bp)[:0], i, 10)
	written, err := w.Write(*bp)
	numBufPool.Put(bp)
	return written, err
}
