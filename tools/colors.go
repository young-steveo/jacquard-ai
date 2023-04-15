package tools

import (
	"io"

	"github.com/fatih/color"
)

var (
	red     = color.New(color.FgRed)
	blue    = color.New(color.FgBlue)
	yellow  = color.New(color.FgYellow)
	green   = color.New(color.FgGreen)
	white   = color.New(color.FgWhite)
	cyan    = color.New(color.FgCyan)
	magenta = color.New(color.FgMagenta)
)

func FprintRed(w io.Writer, s string) {
	red.Fprint(w, s)
}

func FprintfRed(w io.Writer, s string, a ...interface{}) {
	red.Fprintf(w, s, a...)
}

func FprintlnRed(w io.Writer, s string) {
	red.Fprintln(w, s)
}

func FprintBlue(w io.Writer, s string) {
	blue.Fprint(w, s)
}

func FprintfBlue(w io.Writer, s string, a ...interface{}) {
	blue.Fprintf(w, s, a...)
}

func FprintlnBlue(w io.Writer, s string) {
	blue.Fprintln(w, s)
}

func FprintYellow(w io.Writer, s string) {
	yellow.Fprint(w, s)
}

func FprintfYellow(w io.Writer, s string, a ...interface{}) {
	yellow.Fprintf(w, s, a...)
}

func FprintlnYellow(w io.Writer, s string) {
	yellow.Fprintln(w, s)
}

func FprintGreen(w io.Writer, s string) {
	green.Fprint(w, s)
}

func FprintfGreen(w io.Writer, s string, a ...interface{}) {
	green.Fprintf(w, s, a...)
}
func FprintlnGreen(w io.Writer, s string) {
	green.Fprintln(w, s)
}

func FprintWhite(w io.Writer, s string) {
	white.Fprint(w, s)
}

func FprintfWhite(w io.Writer, s string, a ...interface{}) {
	white.Fprintf(w, s, a...)
}

func FprintlnWhite(w io.Writer, s string) {
	white.Fprintln(w, s)
}

func FprintCyan(w io.Writer, s string) {
	cyan.Fprint(w, s)
}

func FprintfCyan(w io.Writer, s string, a ...interface{}) {
	cyan.Fprintf(w, s, a...)
}

func FprintlnCyan(w io.Writer, s string) {
	cyan.Fprintln(w, s)
}

func FprintMagenta(w io.Writer, s string) {
	magenta.Fprint(w, s)
}

func FprintfMagenta(w io.Writer, s string, a ...interface{}) {
	magenta.Fprintf(w, s, a...)
}

func FprintlnMagenta(w io.Writer, s string) {
	magenta.Fprintln(w, s)
}
