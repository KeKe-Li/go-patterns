package chapter02

import "io"

type errWrite struct {
	w   io.Writer
	err error
}

func (e *errWrite) Write(p []byte) {
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(p)
}

func (e *errWrite) Err() error {
	return e.err
}

func do() {
	//ew := errWrite{
	//	w:   nil,
	//	err: nil,
	//}
	//ew.Write(buf)
	//ew.Write(buf)
	//...
	//...
	//if ew.Err() != nil {
	//	return ew.Err()
	//}
	//return nil

}
