package routine

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

const (
	newLine              = "\n"
	innerErrorPrefix     = " ---> "
	endOfInnerErrorStack = "--- End of inner error stack trace ---"
	endOfErrorStack      = "--- End of error stack trace ---"
	wordAt               = "at"
	wordIn               = "in"
	wordCreatedBy        = "created by"
)

type runtimeError struct {
	goid       int64
	gopc       uintptr
	message    string
	stackTrace []uintptr
	cause      RuntimeError
}

func (re *runtimeError) Goid() int64 {
	return re.goid
}

func (re *runtimeError) Gopc() uintptr {
	return re.gopc
}

func (re *runtimeError) Message() string {
	return re.message
}

func (re *runtimeError) StackTrace() []uintptr {
	return re.stackTrace
}

func (re *runtimeError) Cause() RuntimeError {
	return re.cause
}

func (re *runtimeError) Error() string {
	return runtimeErrorError(re)
}

func runtimeErrorNew(cause any) (goid int64, gopc uintptr, msg string, stackTrace []uintptr, innerErr RuntimeError) {
	runtimeErr, isRuntimeErr := cause.(RuntimeError)
	if !isRuntimeErr {
		if err, isErr := cause.(error); isErr {
			msg = err.Error()
		} else if cause != nil {
			msg = fmt.Sprint(cause)
		}
	}
	gp := getg()
	return gp.goid, *gp.gopc, msg, captureStackTrace(2, 100), runtimeErr
}

func runtimeErrorNewWithMessage(message string) (goid int64, gopc uintptr, msg string, stackTrace []uintptr, innerErr RuntimeError) {
	gp := getg()
	return gp.goid, *gp.gopc, message, captureStackTrace(2, 100), nil
}

func runtimeErrorNewWithMessageCause(message string, cause any) (goid int64, gopc uintptr, msg string, stackTrace []uintptr, innerErr RuntimeError) {
	runtimeErr, isRuntimeErr := cause.(RuntimeError)
	if !isRuntimeErr {
		causeMsg := ""
		if err, isErr := cause.(error); isErr {
			causeMsg = err.Error()
		} else if cause != nil {
			causeMsg = fmt.Sprint(cause)
		}
		if len(message) == 0 {
			message = causeMsg
		} else if len(causeMsg) != 0 {
			message += " - " + causeMsg
		}
	}
	gp := getg()
	return gp.goid, *gp.gopc, message, captureStackTrace(2, 100), runtimeErr
}

func runtimeErrorError(re RuntimeError) string {
	builder := &strings.Builder{}
	runtimeErrorPrintStackTrace(re, builder)
	runtimeErrorPrintCreatedBy(re, builder)
	return builder.String()
}

func runtimeErrorPrintStackTrace(re RuntimeError, builder *strings.Builder) {
	builder.WriteString(runtimeErrorTypeName(re))
	message := re.Message()
	if len(message) > 0 {
		builder.WriteString(": ")
		builder.WriteString(message)
	}
	cause := re.Cause()
	if cause != nil {
		builder.WriteString(newLine)
		builder.WriteString(innerErrorPrefix)
		runtimeErrorPrintStackTrace(cause, builder)
		builder.WriteString(newLine)
		builder.WriteString("   ")
		builder.WriteString(endOfInnerErrorStack)
	}
	stackTrace := re.StackTrace()
	if stackTrace != nil {
		frames := runtime.CallersFrames(stackTrace)
		for {
			frame, more := frames.Next()
			if len(frame.Function) > 0 && frame.Function != "runtime.goexit" {
				builder.WriteString(newLine)
				builder.WriteString("   ")
				builder.WriteString(wordAt)
				builder.WriteString(" ")
				builder.WriteString(frame.Function)
				builder.WriteString("() ")
				builder.WriteString(wordIn)
				builder.WriteString(" ")
				builder.WriteString(frame.File)
				builder.WriteString(":")
				builder.WriteString(strconv.Itoa(frame.Line))
			}
			if !more {
				break
			}
		}
	}
}

func runtimeErrorPrintCreatedBy(re RuntimeError, builder *strings.Builder) {
	goid := re.Goid()
	if goid == 1 {
		return
	}
	gopc := re.Gopc()
	fn := runtime.FuncForPC(gopc)
	if fn == nil {
		return
	}
	file, line := fn.FileLine(gopc - 1)
	builder.WriteString(newLine)
	builder.WriteString("   ")
	builder.WriteString(endOfErrorStack)
	builder.WriteString(newLine)
	builder.WriteString("   ")
	builder.WriteString(wordCreatedBy)
	builder.WriteString(" ")
	builder.WriteString(fn.Name())
	builder.WriteString("() ")
	builder.WriteString(wordIn)
	builder.WriteString(" ")
	builder.WriteString(file)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(line))
}

func runtimeErrorTypeName(re RuntimeError) string {
	typeName := []rune(reflect.TypeOf(re).Elem().Name())
	typeName[0] = unicode.ToUpper(typeName[0])
	return string(typeName)
}
