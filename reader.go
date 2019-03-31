package ratelimit

import "io"

type limitedReader struct {
	r io.Reader
	b *Bucket
}

func (r *limitedReader) Read(buf []byte) (written int, err error) {
	bytesToRead := int64(len(buf))
	for bytesToRead > 0 {
		n := r.b.Take(bytesToRead)
		tmpBuf := buf[written : int64(written)+n]
		ret, err := r.r.Read(tmpBuf)
		if err != nil {
			if rest := n - int64(ret); rest > 0 {
				r.b.Return(rest)
			}
			return written + ret, err
		}
		bytesToRead -= int64(ret)
		written += ret
	}
	return
}

// NewReader wraps an io.Reader and applies transfer rate limitation on it.
func (b *Bucket) NewReader(reader io.Reader) io.Reader {
	return &limitedReader{reader, b}
}
