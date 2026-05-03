package main

import "io"

func discard(r io.ReadCloser) {
	if r == nil {
		return
	}
	io.Copy(io.Discard, r)
	r.Close()
}
