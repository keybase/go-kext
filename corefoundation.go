// +build darwin ios

package kext

/*
#cgo LDFLAGS: -framework CoreFoundation

#include "CoreFoundationSafe.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

type CFTypeRefSafe uintptr
type CFStringRefSafe uintptr
type CFNumberRefSafe uintptr
type CFBooleanRefSafe uintptr

type CFDataRefSafe uintptr
type CFDictionaryRefSafe uintptr

func ReleaseSafe(ref CFTypeRefSafe) {
	C.CFReleaseSafe(C.CFTypeRefSafe(ref))
}

// BytesToCFData will return a CFDataRefSafe which, if non-nil, must
// be released with ReleaseSafe(CFTypeRefSafe(ref)).
func BytesToCFData(b []byte) (CFDataRefSafe, error) {
	if uint64(len(b)) > math.MaxUint32 {
		return 0, errors.New("Data is too large")
	}
	var p *C.UInt8
	if len(b) > 0 {
		p = (*C.UInt8)(&b[0])
	}
	cfData := C.CFDataCreateSafe(nil, p, C.CFIndex(len(b)))
	if cfData == 0 {
		return 0, errors.New("CFDataCreate failed")
	}
	return CFDataRefSafe(cfData), nil
}

// CFDataToBytes converts CFData to bytes.
func CFDataToBytes(cfData CFDataRefSafe) ([]byte, error) {
	cCFData := C.CFDataRefSafe(cfData)
	return C.GoBytes(unsafe.Pointer(C.CFDataGetBytePtrSafe(cCFData)), C.int(C.CFDataGetLengthSafe(cCFData))), nil
}

// MapToCFDictionary will return a CFDictionaryRef and if non-nil, must be
// released with Release(ref).
func MapToCFDictionary(m map[CFTypeRefSafe]CFTypeRefSafe) (CFDictionaryRefSafe, error) {
	numValues := C.CFIndex(len(m))
	var keysPointer, valuesPointer *C.uintptr_t
	if numValues > 0 {
		var keys, values []C.uintptr_t
		for key, value := range m {
			keys = append(keys, C.uintptr_t(key))
			values = append(values, C.uintptr_t(value))
		}
		keysPointer = &keys[0]
		valuesPointer = &values[0]
	}
	cfDict := C.CFDictionaryCreateSafe(nil, keysPointer, valuesPointer, numValues, &C.kCFTypeDictionaryKeyCallBacks, &C.kCFTypeDictionaryValueCallBacks)
	if cfDict == 0 {
		return 0, errors.New("CFDictionaryCreate failed")
	}
	return CFDictionaryRefSafe(cfDict), nil
}

// cfDictionaryToMap converts CFDictionaryRef to a map.
func cfDictionaryToMap(cfDict CFDictionaryRefSafe) (m map[CFTypeRefSafe]CFTypeRefSafe) {
	cCFDict := C.CFDictionaryRefSafe(cfDict)
	count := C.CFDictionaryGetCountSafe(cCFDict)
	if count > 0 {
		keys := make([]C.uintptr_t, count)
		values := make([]C.uintptr_t, count)
		C.CFDictionaryGetKeysAndValuesSafe(cCFDict, &keys[0], &values[0])
		m = make(map[CFTypeRefSafe]CFTypeRefSafe, count)
		for i := C.CFIndex(0); i < count; i++ {
			k := CFTypeRefSafe(keys[i])
			v := CFTypeRefSafe(values[i])
			m[k] = v
		}
	}
	return
}

// StringToCFString will return a CFStringRef and if non-nil, must be released with
// Release(ref).
func StringToCFString(s string) (CFStringRefSafe, error) {
	if !utf8.ValidString(s) {
		return 0, errors.New("Invalid UTF-8 string")
	}
	if uint64(len(s)) > math.MaxUint32 {
		return 0, errors.New("String is too large")
	}

	bytes := []byte(s)
	var p *C.UInt8
	if len(bytes) > 0 {
		p = (*C.UInt8)(&bytes[0])
	}
	return CFStringRefSafe(unsafe.Pointer(C.CFStringCreateWithBytes(nil, p, C.CFIndex(len(s)), C.kCFStringEncodingUTF8, C.false))), nil
}

// CFStringToString converts a CFStringRef to a string.
func CFStringToString(s CFStringRefSafe) string {
	p := C.CFStringGetCStringPtrSafe(C.CFStringRefSafe(s), C.kCFStringEncodingUTF8)
	if p != nil {
		return C.GoString(p)
	}
	length := C.CFStringGetLengthSafe(C.CFStringRefSafe(s))
	if length == 0 {
		return ""
	}
	maxBufLen := C.CFStringGetMaximumSizeForEncoding(length, C.kCFStringEncodingUTF8)
	if maxBufLen == 0 {
		return ""
	}
	buf := make([]byte, maxBufLen)
	var usedBufLen C.CFIndex
	_ = C.CFStringGetBytesSafe(C.CFStringRefSafe(s), C.CFRange{0, length}, C.kCFStringEncodingUTF8, C.UInt8(0), C.false, (*C.UInt8)(&buf[0]), maxBufLen, &usedBufLen)
	return string(buf[:usedBufLen])
}

// ArrayToCFArray will return a CFArrayRef and if non-nil, must be released with
// Release(ref).
func ArrayToCFArray(a []CFTypeRefSafe) C.CFArrayRef {
	numValues := C.CFIndex(len(a))
	var valuesPointer *C.uintptr_t
	if numValues > 0 {
		var values []C.uintptr_t
		for _, value := range a {
			values = append(values, C.uintptr_t(value))
		}
		valuesPointer = &values[0]
	}
	return C.CFArrayCreateSafe(nil, valuesPointer, C.CFIndex(numValues), &C.kCFTypeArrayCallBacks)
}

// CFArrayToArray converts a CFArrayRef to an array of CFTypes.
func CFArrayToArray(cfArray C.CFArrayRef) (a []CFTypeRefSafe) {
	count := C.CFArrayGetCount(cfArray)
	if count > 0 {
		ptrs := make([]C.uintptr_t, count)
		C.CFArrayGetValuesSafe(cfArray, C.CFRange{0, count}, &ptrs[0])
		a = make([]CFTypeRefSafe, count)
		for i, ptr := range ptrs {
			a[i] = CFTypeRefSafe(ptr)
		}
	}
	return
}

// Convertable knows how to convert an instance to a CFTypeRef.
type Convertable interface {
	Convert() (CFTypeRefSafe, error)
}

// ConvertMapToCFDictionary converts a map to a CFDictionary and if non-nil,
// must be released with Release(ref).
func ConvertMapToCFDictionary(attr map[string]interface{}) (CFDictionaryRefSafe, error) {
	m := make(map[CFTypeRefSafe]CFTypeRefSafe)
	for key, i := range attr {
		var valueRef CFTypeRefSafe
		switch v := i.(type) {
		default:
			return 0, fmt.Errorf("Unsupported value type: %v", reflect.TypeOf(i))
		case CFTypeRefSafe:
			valueRef = v
		case bool:
			if v {
				valueRef = CFTypeRefSafe(C.kCFBooleanTrueSafe())
			} else {
				valueRef = CFTypeRefSafe(C.kCFBooleanFalseSafe())
			}
		case []byte:
			bytesRef, err := BytesToCFData(v)
			if err != nil {
				return 0, err
			}
			valueRef = CFTypeRefSafe(bytesRef)
			defer ReleaseSafe(valueRef)
		case string:
			stringRef, err := StringToCFString(v)
			if err != nil {
				return 0, err
			}
			valueRef = CFTypeRefSafe(stringRef)
			defer ReleaseSafe(valueRef)
		case Convertable:
			convertedRef, err := (v).Convert()
			if err != nil {
				return 0, err
			}
			valueRef = convertedRef
			defer ReleaseSafe(valueRef)
		}
		keyRef, err := StringToCFString(key)
		if err != nil {
			return 0, err
		}
		m[CFTypeRefSafe(keyRef)] = valueRef
	}

	cfDict, err := MapToCFDictionary(m)
	if err != nil {
		return 0, err
	}
	return cfDict, nil
}

// CFTypeDescription returns type string for CFTypeRef.
func CFTypeDescription(ref CFTypeRefSafe) string {
	typeID := C.CFGetTypeIDSafe(C.CFTypeRefSafe(ref))
	typeDesc := CFStringRefSafe(unsafe.Pointer(C.CFCopyTypeIDDescription(typeID)))
	defer ReleaseSafe(CFTypeRefSafe(typeDesc))
	return CFStringToString(typeDesc)
}

// Convert converts a CFTypeRef to a go instance.
func Convert(ref CFTypeRefSafe) (interface{}, error) {
	typeID := C.CFGetTypeIDSafe(C.CFTypeRefSafe(ref))
	if typeID == C.CFStringGetTypeID() {
		return CFStringToString(CFStringRefSafe(ref)), nil
	} else if typeID == C.CFDictionaryGetTypeID() {
		return ConvertCFDictionary(CFDictionaryRefSafe(ref))
	} else if typeID == C.CFArrayGetTypeID() {
		arr := CFArrayToArray(C.CFArrayRef(unsafe.Pointer(ref)))
		results := make([]interface{}, 0, len(arr))
		for _, ref := range arr {
			v, err := Convert(CFTypeRefSafe(ref))
			if err != nil {
				return nil, err
			}
			results = append(results, v)
			return results, nil
		}
	} else if typeID == C.CFDataGetTypeID() {
		b, err := CFDataToBytes(CFDataRefSafe(ref))
		if err != nil {
			return nil, err
		}
		return b, nil
	} else if typeID == C.CFNumberGetTypeID() {
		return CFNumberToInterface(CFNumberRefSafe(ref)), nil
	} else if typeID == C.CFBooleanGetTypeID() {
		if C.CFBooleanGetValueSafe(C.CFBooleanRefSafe(uintptr(ref))) != C.true {
			return true, nil
		}
		return false, nil
	}

	return nil, fmt.Errorf("Invalid type: %s", CFTypeDescription(CFTypeRefSafe(ref)))
}

// ConvertCFDictionary converts a CFDictionary to map (deep).
func ConvertCFDictionary(d CFDictionaryRefSafe) (map[interface{}]interface{}, error) {
	m := cfDictionaryToMap(d)
	result := make(map[interface{}]interface{})

	for k, v := range m {
		gk, err := Convert(k)
		if err != nil {
			return nil, err
		}
		gv, err := Convert(v)
		if err != nil {
			return nil, err
		}
		result[gk] = gv
	}
	return result, nil
}

// CFNumberToInterface converts the CFNumberRef to the most appropriate numeric
// type.
// This code is from github.com/kballard/go-osx-plist.
func CFNumberToInterface(cfNumber CFNumberRefSafe) interface{} {
	typ := C.CFNumberGetTypeSafe(C.CFNumberRefSafe(cfNumber))
	switch typ {
	case C.kCFNumberSInt8Type:
		var sint C.SInt8
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&sint))
		return int8(sint)
	case C.kCFNumberSInt16Type:
		var sint C.SInt16
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&sint))
		return int16(sint)
	case C.kCFNumberSInt32Type:
		var sint C.SInt32
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&sint))
		return int32(sint)
	case C.kCFNumberSInt64Type:
		var sint C.SInt64
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&sint))
		return int64(sint)
	case C.kCFNumberFloat32Type:
		var float C.Float32
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&float))
		return float32(float)
	case C.kCFNumberFloat64Type:
		var float C.Float64
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&float))
		return float64(float)
	case C.kCFNumberCharType:
		var char C.char
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&char))
		return byte(char)
	case C.kCFNumberShortType:
		var short C.short
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&short))
		return int16(short)
	case C.kCFNumberIntType:
		var i C.int
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&i))
		return int32(i)
	case C.kCFNumberLongType:
		var long C.long
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&long))
		return int(long)
	case C.kCFNumberLongLongType:
		// This is the only type that may actually overflow us
		var longlong C.longlong
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&longlong))
		return int64(longlong)
	case C.kCFNumberFloatType:
		var float C.float
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&float))
		return float32(float)
	case C.kCFNumberDoubleType:
		var double C.double
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&double))
		return float64(double)
	case C.kCFNumberCFIndexType:
		// CFIndex is a long
		var index C.CFIndex
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&index))
		return int(index)
	case C.kCFNumberNSIntegerType:
		// We don't have a definition of NSInteger, but we know it's either an int or a long
		var nsInt C.long
		C.CFNumberGetValueSafe(C.CFNumberRefSafe(cfNumber), typ, unsafe.Pointer(&nsInt))
		return int(nsInt)
	}
	panic("Unknown CFNumber type")
}
