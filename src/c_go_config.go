package main

// This file exists solely to help the VSCode Go language server
// find and understand our CGO dependencies.

/*
#cgo CFLAGS: -I${SRCDIR}/c -I${SRCDIR}/../sdk
*/
import "C"
