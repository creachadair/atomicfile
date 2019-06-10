# atomicfile

http://godoc.org/github.com/creachadair/atomicfile

[![Go Report Card](https://goreportcard.com/badge/github.com/creachadair/atomicfile)](https://goreportcard.com/report/github.com/creachadair/atomicfile)

Package atomicfile implements all-or-nothing file replacement by staging output
to a temporary file adjacent to the specified target file, and renaming over
the target when the temporary is closed.
