package loglv

import (
	"fmt"
	"runtime"
)

type LOGlevl int

const (
	ALL LOGlevl = iota
	DBG LOGlevl = iota
	INF LOGlevl = iota
	WAR LOGlevl = iota
	ERR LOGlevl = iota
	FTA LOGlevl = iota
	OFF LOGlevl = iota
)

var G_currlogLevl LOGlevl = 0

type levlLOG struct {
	m_loglevl LOGlevl
	m_prefix  string
}

var Dbg levlLOG = levlLOG{m_loglevl: DBG, m_prefix: "DBG"} //DEBUG
var Inf levlLOG = levlLOG{m_loglevl: INF, m_prefix: "INF"} //INFO
var War levlLOG = levlLOG{m_loglevl: WAR, m_prefix: "WAR"} //WARN
var Err levlLOG = levlLOG{m_loglevl: ERR, m_prefix: "ERR"} //ERROR
var Fta levlLOG = levlLOG{m_loglevl: FTA, m_prefix: "FTA"} //FATAL

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func (this *levlLOG) Print(v ...interface{}) {
	if this.m_loglevl < G_currlogLevl {
		return
	}
	std.Output(this.m_prefix, 2, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (this *levlLOG) Printf(format string, v ...interface{}) {
	if this.m_loglevl < G_currlogLevl {
		return
	}
	std.Output(this.m_prefix, 2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func (this *levlLOG) Println(v ...interface{}) {
	if this.m_loglevl < G_currlogLevl {
		return
	}
	std.Output(this.m_prefix, 2, fmt.Sprintln(v...))
}

//v可以额外传递最后一个参数panicfunc func()
func (this *levlLOG) Recoverf(format string, v ...interface{}) (recover_err interface{}) {
	if e := recover(); e != nil {
		var panicfun func()
		if vlen := len(v); vlen > 0 {
			if panicfun, _ = v[vlen-1].(func()); panicfun != nil {
				v = v[:vlen-1]
			}
		}
		if this.m_loglevl >= G_currlogLevl {
			l_buf := [2048]byte{'\n'}
			l_stack := l_buf[:runtime.Stack(l_buf[1:len(l_buf)-2], false)+2]
			l_stack[len(l_stack)-1] = '\n'
			std.Output("\n"+this.m_prefix, 2, fmt.Sprintf(format, v...)+fmt.Sprintf(":%v", e)+string(l_stack))
		}
		if panicfun != nil {
			panicfun()
		}
		return e
	}
	return nil
}
func (this *levlLOG) RecoverCTX(panicfun func(err interface{}), set_context func() string) (recover_err interface{}) {
	if e := recover(); e != nil {
		if this.m_loglevl >= G_currlogLevl {
			l_buf := [2048]byte{'\n'}
			l_stack := l_buf[:runtime.Stack(l_buf[1:len(l_buf)-2], false)+2]
			l_stack[len(l_stack)-1] = '\n'
			if set_context != nil {
				std.Output("\n"+this.m_prefix, 2, set_context()+fmt.Sprintf(":%v", e)+string(l_stack))
			} else {
				std.Output("\n"+this.m_prefix, 2, fmt.Sprintf("%v", e)+string(l_stack))
			}
		}
		if panicfun != nil {
			panicfun(e)
		}
		return e
	}
	return nil
}
func (this *levlLOG) Recoverln(format string, v ...interface{}) (recover_err interface{}) {
	if e := recover(); e != nil {
		if this.m_loglevl >= G_currlogLevl {
			std.Output(this.m_prefix, 2, fmt.Sprintf(format, v...)+fmt.Sprintf(":%v\n", e))
		}
		return e
	}
	return nil
}

func (this *levlLOG) Stackf(format string, v ...interface{}) {
	if this.m_loglevl >= G_currlogLevl {
		l_buf := [2048]byte{'\n'}
		l_stack := l_buf[:runtime.Stack(l_buf[1:len(l_buf)-2], false)+2]
		l_stack[len(l_stack)-1] = '\n'
		std.Output("\n"+this.m_prefix, 2, fmt.Sprintf(format, v...)+string(l_stack))
	}
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if Llongfile or Lshortfile is set; a value of 1 will print the details
// for the caller of Output.
func (this *levlLOG) Output(calldepth int, s string) error {
	return std.Output(this.m_prefix, calldepth+1, s) // +1 for this frame.
}
