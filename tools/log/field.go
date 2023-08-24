package log

import (
	"go.uber.org/zap"
)

type Field = zap.Field

var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	BoolP       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128P = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64P  = zap.Complex64p
	Float64     = zap.Float64
	Float64P    = zap.Float64p
	Float32     = zap.Float32
	Float32P    = zap.Float32p
	Int         = zap.Int
	IntP        = zap.Intp
	Int64       = zap.Int64
	Int64P      = zap.Int64p
	Int32       = zap.Int32
	Int32P      = zap.Int32p
	Int16       = zap.Int16
	Int16P      = zap.Int16p
	Int8        = zap.Int8
	Int8P       = zap.Int8p
	String      = zap.String
	StringP     = zap.Stringp
	Uint        = zap.Uint
	UintP       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64P     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32P     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16P     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8P      = zap.Uint8p
	Uintptr     = zap.Uintptr
	UintptrP    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	TimeP       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	DurationP   = zap.Durationp
	Object      = zap.Object
	Inline      = zap.Inline
	Any         = zap.Any
	Err         = zap.Error
)
