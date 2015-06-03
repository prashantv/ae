package ae

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
)

const (
	// skips runtime related functions on the stack.
	skipPastRuntimeCall = true

	// StackSeparator is the separator between stack frames.
	StackSeparator = "; "
)

func (ae *appError) Stack() string {
	var buf bytes.Buffer
	for _, f := range ae.getStack() {
		if buf.Len() > 0 {
			buf.WriteString(StackSeparator)
		}
		buf.WriteString(fmt.Sprintf("%s:%v %s: %s", f.file, f.line, f.funcName, f.lineContents))
	}
	return buf.String()
}

type stackFrame struct {
	funcName     string
	file         string
	line         int
	lineContents string
}

func (ae *appError) getStack() []stackFrame {
	if ae.frameCache == nil {
		// Use non-nil value in case we store an empty list
		ae.frameCache = []stackFrame{}
		for _, pc := range ae.stack {
			f := runtime.FuncForPC(pc)
			file, line := f.FileLine(pc)

			if skipPastRuntimeCall && strings.Contains(f.Name(), "runtime.call64") {
				break
			}
			ae.frameCache = append(ae.frameCache, frame(f, file, line))
		}
	}
	return ae.frameCache
}

func frame(f *runtime.Func, file string, line int) stackFrame {
	// Format is [file]:[line] [fName]: [lineContents]
	return stackFrame{f.Name(), file, line, getFileLine(file, line)}
}

func getFileLine(file string, line int) string {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return "[file not found]"
	}
	bss := string(bs)
	lines := strings.Split(bss, "\n")
	if line >= len(lines) {
		return "[line not found]"
	}

	return strings.Trim(lines[line-1], " \t")
}

// PrintTolog prints the error and stack using the default logger.
func (ae *appError) PrintTolog() {
	var buf bytes.Buffer

	buf.WriteString(ae.errorMsgs())
	buf.WriteString("\n")
	for _, f := range ae.getStack() {
		buf.WriteString(fmt.Sprintf("  %s:%v\n    %s: %s\n", f.file, f.line, f.funcName, f.lineContents))
	}
	log.Println(buf.String())
}

// getStack returns the list of callers (ignoring the most recent skip callers)
func getStackPC(skip int) []uintptr {
	const bufSize = 20
	// +2 since we want to ignore the getStack and the runtime.Callers call.
	skip = skip + 2
	var pc []uintptr
	for numCallers := bufSize; numCallers == bufSize; skip += numCallers {
		buf := make([]uintptr, bufSize)
		numCallers = runtime.Callers(skip, buf)
		pc = append(pc, buf[:numCallers]...)
	}
	return pc
}
