# atomicfile

http://godoc.org/bitbucket.org/creachadair/atomicfile

[![Go Report Card](https://goreportcard.com/badge/bitbucket.org/creachadair/atomicfile)](https://goreportcard.com/report/bitbucket.org/creachadair/atomicfile)

Package atomicfile implements all-or-nothing file replacement by staging output
to a temporary file adjacent to the specified target file, and renaming over
the target when the temporary is closed.
