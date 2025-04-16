package fastlog

// trimCallerPath takes a caller, and returns #count number of parts from end of caller
func trimCallerPath(buf *Buffer, caller string, count int) {
	start := 0
	for i := len(caller) - 1; i >= 0 && count > 0; i-- {
		if caller[i] == '/' {
			start = i + 1
			count--
		}
	}
	buf.store = append(buf.store, caller[start:]...)
	// buf.Append(caller[start:])
}
