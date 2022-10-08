// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. It also has a predefined 'standard'
// Logger accessible through helper functions Print[f|ln], Fatal[f|ln], and
// Panic[f|ln], which are easier to use than creating a Logger manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.
package loglv

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// These flags define which text to prefix to each log entry generated by the Logger.
const (
	// Bits or'ed together to control what's printed.
	// There is no control over the order they appear (the order listed
	// here) or the format they present (as described in the comments).
	// The prefix is followed by a colon only when Llongfile or Lshortfile
	// is specified.
	// For example, flags Ldate | Ltime (or LstdFlags) produce,
	//	2009/01/23 01:23:23 message
	// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu   sync.Mutex // ensures atomic writes; protects the following fields
	flag int        // properties
	out  io.Writer  // destination for output
	buf  []byte     // for accumulating text to write
	name string     // logfile's fullpath
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) formatHeader(prefix string, buf *[]byte, t time.Time, file string, line int) {
	if len(prefix) > 0 {
		*buf = append(*buf, prefix...)
		*buf = append(*buf, '/')
	}
	if l.flag&LUTC != 0 {
		t = t.UTC()
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			_, month, day := t.Date()
			// itoa(buf, year, 4)
			// *buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(prefix string, calldepth int, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(prefix, &l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}
func (l *Logger) Name() string { return l.name }

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(out io.Writer, flag int) *Logger {
	lName := ""
	if fout, ok := out.(*os.File); ok {
		if fout != nil && fout != os.Stderr && fout != os.Stdout {
			lName = fout.Name()
		}
	}
	return &Logger{out: out, flag: flag, name: lName}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(bakName, newName string, w io.Writer) (err error) {

	lName, newfile := l.Name(), (*os.File)(nil)

	if lName != newName && newName != "" {
		newfile, err = os.OpenFile(newName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if lName != "" {
		//关闭现有日志文件
		if f, ok := l.out.(*os.File); ok {
			if err = f.Sync(); err != nil {
				if newfile != nil {
					newfile.Close()
				}
				return
			} else if err = f.Close(); err != nil {
				if newfile != nil {
					newfile.Close()
				}
				return
			}
		}
		//更名现有日志文件
		if bakName == "" && lName == newName { //没有传递备份名但新老名字有同名冲突，产生一个日期备份名字
			os.Mkdir(GetLogDir(), 0755)
			bakName = GetLogDir() + "/" + time.Now().Format("06-01-02_15.04.05_") + GetExeBaseName() + ".log"
		}
		if bakName != "" {
			os.Rename(lName, bakName)
		}
	}
	if newfile == nil && newName != "" {
		newfile, err = os.OpenFile(newName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			l.out = os.Stderr
			l.name = ""
			return
		}
	}
	if newfile != nil {
		l.out = newfile
		l.name = newName
	} else {
		l.out = w
		l.name = ""
		if fout, ok := w.(*os.File); ok {
			if fout != nil && fout != os.Stderr && fout != os.Stdout {
				l.name = fout.Name()
			}
		}
	}
	return
}

var std = New(os.Stderr, LstdFlags)

// SetOutput sets the output destination for the standard logger.
func SetOutput(bakName, newName string, w io.Writer) error {
	return std.SetOutput(bakName, newName, w)
}

// These functions write to the standard logger.

// Printf calls std.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(l levlLOG, format string, v ...interface{}) {
	std.Output(l.m_prefix, 2, fmt.Sprintf(format, v...))
}

// Print calls std.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func Print(l levlLOG, v ...interface{}) { std.Output(l.m_prefix, 2, fmt.Sprint(v...)) }

// Println calls std.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func Println(l levlLOG, v ...interface{}) { std.Output(l.m_prefix, 2, fmt.Sprintln(v...)) }

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(Fta.m_prefix, 2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	std.Output(Fta.m_prefix, 2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	std.Output(Fta.m_prefix, 2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(Fta.m_prefix, 2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(Fta.m_prefix, 2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(Fta.m_prefix, 2, s)
	panic(s)
}
