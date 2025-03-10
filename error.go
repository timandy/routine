package routine

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
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
	goid       uint64
	gopc       uintptr
	message    string
	stackTrace []uintptr
	cause      RuntimeError
}

func (re *runtimeError) Goid() uint64 {
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

func runtimeErrorNew(cause any) (goid uint64, gopc uintptr, msg string, stackTrace []uintptr, innerErr RuntimeError) {
	runtimeErr, isRuntimeErr := cause.(RuntimeError)
	if !isRuntimeErr {
		if err, isErr := cause.(error); isErr {
			msg = err.Error()
		} else if cause != nil {
			msg = fmt.Sprint(cause)
		}
	}
	gp := getg()
	return gp.goid(), gp.gopc(), msg, captureStackTrace(2, 100), runtimeErr
}

func runtimeErrorNewWithMessage(message string) (goid uint64, gopc uintptr, msg string, stackTrace []uintptr, innerErr RuntimeError) {
	gp := getg()
	return gp.goid(), gp.gopc(), message, captureStackTrace(2, 100), nil
}

func runtimeErrorNewWithMessageCause(message string, cause any) (goid uint64, gopc uintptr, msg string, stackTrace []uintptr, innerErr RuntimeError) {
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
	return gp.goid(), gp.gopc(), message, captureStackTrace(2, 100), runtimeErr
}

func runtimeErrorError(re RuntimeError) string {
	builder := &bytes.Buffer{}
	runtimeErrorPrintStackTrace(re, builder)
	runtimeErrorPrintCreatedBy(re, builder)
	return builder.String()
}

func runtimeErrorPrintStackTrace(re RuntimeError, builder *bytes.Buffer) {
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
		savePoint := builder.Len()
		skippedPanic := false
		frames := runtime.CallersFrames(stackTrace)
		for {
			frame, more := frames.Next()
			if showFrame(frame.Function) {
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
			} else if skipFrame(frame.Function, skippedPanic) {
				builder.Truncate(savePoint)
				skippedPanic = true
			}
			if !more {
				break
			}
		}
	}
}

func runtimeErrorPrintCreatedBy(re RuntimeError, builder *bytes.Buffer) {
	goid := re.Goid()
	if goid == 1 {
		return
	}
	pc := re.Gopc()
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	if frame.Func == nil {
		return
	}
	builder.WriteString(newLine)
	builder.WriteString("   ")
	builder.WriteString(endOfErrorStack)
	builder.WriteString(newLine)
	builder.WriteString("   ")
	builder.WriteString(wordCreatedBy)
	builder.WriteString(" ")
	builder.WriteString(frame.Function)
	builder.WriteString("() ")
	builder.WriteString(wordIn)
	builder.WriteString(" ")
	builder.WriteString(frame.File)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(frame.Line))
}

func runtimeErrorTypeName(re RuntimeError) string {
	typeName := []rune(reflect.TypeOf(re).Elem().Name())
	typeName[0] = unicode.ToUpper(typeName[0])
	return string(typeName)
}
