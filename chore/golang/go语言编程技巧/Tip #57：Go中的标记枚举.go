package main

// flag enum
type BasicInfo int

const (
	IsBoolean  BasicInfo = 1 << iota // 1
	IsInteger                        // 2
	IsUnsigned                       // 4
	IsFloat                          // 8
	IsComplex                        // 16
	IsString                         // 32
	IsUntyped                        // 64

	IsOrdered   = IsInteger | IsFloat | IsComplex  // 42
	IsNumeric   = IsInteger | IsFloat | IsComplex  // 26
	IsConstType = IsBoolean | IsNumeric | IsString // 59
)
