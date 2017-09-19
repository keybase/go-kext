// +build darwin ios

package kext

/*
#cgo LDFLAGS: -framework CoreFoundation

#include <stdint.h>
#include <CoreFoundation/CoreFoundation.h>

typedef uintptr_t CFTypeRefSafe;
typedef uintptr_t CFStringRefSafe;

void CFReleaseSafe(CFTypeRefSafe cf) {
  CFRelease((CFTypeRef)cf);
}

CFArrayRef CFArrayCreateSafe(CFAllocatorRef allocator, const uintptr_t *values, CFIndex numValues, const CFArrayCallBacks *callBacks) {
  return CFArrayCreate(allocator, (const void **)values, numValues, callBacks);
}

CFDictionaryRef CFDictionaryCreateSafe(CFAllocatorRef allocator, const uintptr_t *keys, const uintptr_t *values, CFIndex numValues, const CFDictionaryKeyCallBacks *keyCallBacks, const CFDictionaryValueCallBacks *valueCallBacks) {
  return CFDictionaryCreate(allocator, (const void **)keys, (const void **)values, numValues, keyCallBacks, valueCallBacks);
}

const char * CFStringGetCStringPtrSafe(CFStringRefSafe theString, CFStringEncoding encoding) {
  return CFStringGetCStringPtr((CFStringRef)theString, encoding);
}

CFIndex CFStringGetLengthSafe(CFStringRefSafe theString) {
  return CFStringGetLength((CFStringRef)theString);
}

CFIndex CFStringGetBytesSafe(CFStringRefSafe theString, CFRange range, CFStringEncoding encoding, UInt8 lossByte, Boolean isExternalRepresentation, UInt8 *buffer, CFIndex maxBufLen, CFIndex *usedBufLen) {
  return CFStringGetBytes((CFStringRef)theString, range, encoding, lossByte, isExternalRepresentation, buffer, maxBufLen, usedBufLen);
}
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

func ReleaseSafe(ref CFTypeRefSafe) {
	C.CFReleaseSafe(C.CFTypeRefSafe(ref))
}

// BytesToCFData will return a CFDataRef and if non-nil, must be released with
// ReleaseSafe(CFTypeRefSafe(unsafe.Pointer(ref))).
func BytesToCFData(b []byte) (C.CFDataRef, error) {
	if uint64(len(b)) > math.MaxUint32 {
		return nil, errors.New("Data is too large")
	}
	var p *C.UInt8
	if len(b) > 0 {
		p = (*C.UInt8)(&b[0])
	}
	cfData := C.CFDataCreate(nil, p, C.CFIndex(len(b)))
	if cfData == nil {
		return nil, fmt.Errorf("CFDataCreate failed")
	}
	return cfData, nil
}

// CFDataToBytes converts CFData to bytes.
func CFDataToBytes(cfData C.CFDataRef) ([]byte, error) {
	return C.GoBytes(unsafe.Pointer(C.CFDataGetBytePtr(cfData)), C.int(C.CFDataGetLength(cfData))), nil
}

// MapToCFDictionary will return a CFDictionaryRef and if non-nil, must be
// released with Release(ref).
func MapToCFDictionary(m map[CFTypeRefSafe]CFTypeRefSafe) (C.CFDictionaryRef, error) {
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
	if cfDict == nil {
		return nil, fmt.Errorf("CFDictionaryCreate failed")
	}
	return cfDict, nil
}

// cfDictionaryToMap converts CFDictionaryRef to a map.
func cfDictionaryToMap(cfDict C.CFDictionaryRef) (m map[CFTypeRefSafe]CFTypeRefSafe) {
	count := C.CFDictionaryGetCount(cfDict)
	if count > 0 {
		keys := make([]CFTypeRefSafe, count)
		values := make([]CFTypeRefSafe, count)
		C.CFDictionaryGetKeysAndValues(cfDict, (*unsafe.Pointer)(unsafe.Pointer(&keys[0])), (*unsafe.Pointer)(unsafe.Pointer(&values[0])))
		m = make(map[CFTypeRefSafe]CFTypeRefSafe, count)
		for i := C.CFIndex(0); i < count; i++ {
			k := keys[i]
			v := values[i]
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
func CFArrayToArray(cfArray C.CFArrayRef) (a []C.CFTypeRef) {
	count := C.CFArrayGetCount(cfArray)
	if count > 0 {
		a = make([]C.CFTypeRef, count)
		C.CFArrayGetValues(cfArray, C.CFRange{0, count}, (*unsafe.Pointer)(&a[0]))
	}
	return
}

// Convertable knows how to convert an instance to a CFTypeRef.
type Convertable interface {
	Convert() (C.CFTypeRef, error)
}

// ConvertMapToCFDictionary converts a map to a CFDictionary and if non-nil,
// must be released with Release(ref).
func ConvertMapToCFDictionary(attr map[string]interface{}) (C.CFDictionaryRef, error) {
	m := make(map[CFTypeRefSafe]CFTypeRefSafe)
	for key, i := range attr {
		var valueRef CFTypeRefSafe
		switch v := i.(type) {
		default:
			return nil, fmt.Errorf("Unsupported value type: %v", reflect.TypeOf(i))
		case C.CFTypeRef:
			valueRef = CFTypeRefSafe(v)
		case bool:
			if i == true {
				valueRef = CFTypeRefSafe(unsafe.Pointer(C.kCFBooleanTrue))
			} else {
				valueRef = CFTypeRefSafe(unsafe.Pointer(C.kCFBooleanFalse))
			}
		case []byte:
			bytesRef, err := BytesToCFData(i.([]byte))
			if err != nil {
				return nil, err
			}
			defer ReleaseSafe(CFTypeRefSafe(unsafe.Pointer(bytesRef)))
		case string:
			stringRef, err := StringToCFString(i.(string))
			if err != nil {
				return nil, err
			}
			defer ReleaseSafe(CFTypeRefSafe(stringRef))
		case Convertable:
			convertedRef, err := (i.(Convertable)).Convert()
			if err != nil {
				return nil, err
			}
			defer ReleaseSafe(CFTypeRefSafe(unsafe.Pointer(convertedRef)))
		}
		keyRef, err := StringToCFString(key)
		if err != nil {
			return nil, err
		}
		m[CFTypeRefSafe(unsafe.Pointer(keyRef))] = valueRef
	}

	cfDict, err := MapToCFDictionary(m)
	if err != nil {
		return nil, err
	}
	return cfDict, nil
}

// CFTypeDescription returns type string for CFTypeRef.
func CFTypeDescription(ref C.CFTypeRef) string {
	typeID := C.CFGetTypeID(ref)
	typeDesc := CFStringRefSafe(unsafe.Pointer(C.CFCopyTypeIDDescription(typeID)))
	defer ReleaseSafe(CFTypeRefSafe(typeDesc))
	return CFStringToString(typeDesc)
}

// Convert converts a CFTypeRef to a go instance.
func Convert(ref C.CFTypeRef) (interface{}, error) {
	typeID := C.CFGetTypeID(ref)
	if typeID == C.CFStringGetTypeID() {
		return CFStringToString(CFStringRefSafe(ref)), nil
	} else if typeID == C.CFDictionaryGetTypeID() {
		return ConvertCFDictionary(C.CFDictionaryRef(ref))
	} else if typeID == C.CFArrayGetTypeID() {
		arr := CFArrayToArray(C.CFArrayRef(ref))
		results := make([]interface{}, 0, len(arr))
		for _, ref := range arr {
			v, err := Convert(ref)
			if err != nil {
				return nil, err
			}
			results = append(results, v)
			return results, nil
		}
	} else if typeID == C.CFDataGetTypeID() {
		b, err := CFDataToBytes(C.CFDataRef(ref))
		if err != nil {
			return nil, err
		}
		return b, nil
	} else if typeID == C.CFNumberGetTypeID() {
		return CFNumberToInterface(C.CFNumberRef(ref)), nil
	} else if typeID == C.CFBooleanGetTypeID() {
		if C.CFBooleanGetValue(C.CFBooleanRef(ref)) != 0 {
			return true, nil
		}
		return false, nil
	}

	return nil, fmt.Errorf("Invalid type: %s", CFTypeDescription(ref))
}

// ConvertCFDictionary converts a CFDictionary to map (deep).
func ConvertCFDictionary(d C.CFDictionaryRef) (map[interface{}]interface{}, error) {
	m := cfDictionaryToMap(C.CFDictionaryRef(d))
	result := make(map[interface{}]interface{})

	for k, v := range m {
		gk, err := Convert(C.CFTypeRef(k))
		if err != nil {
			return nil, err
		}
		gv, err := Convert(C.CFTypeRef(v))
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
func CFNumberToInterface(cfNumber C.CFNumberRef) interface{} {
	typ := C.CFNumberGetType(cfNumber)
	switch typ {
	case C.kCFNumberSInt8Type:
		var sint C.SInt8
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&sint))
		return int8(sint)
	case C.kCFNumberSInt16Type:
		var sint C.SInt16
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&sint))
		return int16(sint)
	case C.kCFNumberSInt32Type:
		var sint C.SInt32
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&sint))
		return int32(sint)
	case C.kCFNumberSInt64Type:
		var sint C.SInt64
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&sint))
		return int64(sint)
	case C.kCFNumberFloat32Type:
		var float C.Float32
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&float))
		return float32(float)
	case C.kCFNumberFloat64Type:
		var float C.Float64
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&float))
		return float64(float)
	case C.kCFNumberCharType:
		var char C.char
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&char))
		return byte(char)
	case C.kCFNumberShortType:
		var short C.short
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&short))
		return int16(short)
	case C.kCFNumberIntType:
		var i C.int
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&i))
		return int32(i)
	case C.kCFNumberLongType:
		var long C.long
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&long))
		return int(long)
	case C.kCFNumberLongLongType:
		// This is the only type that may actually overflow us
		var longlong C.longlong
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&longlong))
		return int64(longlong)
	case C.kCFNumberFloatType:
		var float C.float
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&float))
		return float32(float)
	case C.kCFNumberDoubleType:
		var double C.double
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&double))
		return float64(double)
	case C.kCFNumberCFIndexType:
		// CFIndex is a long
		var index C.CFIndex
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&index))
		return int(index)
	case C.kCFNumberNSIntegerType:
		// We don't have a definition of NSInteger, but we know it's either an int or a long
		var nsInt C.long
		C.CFNumberGetValue(cfNumber, typ, unsafe.Pointer(&nsInt))
		return int(nsInt)
	}
	panic("Unknown CFNumber type")
}
