#ifndef CORE_FOUNDATION_SAFE_H
#define CORE_FOUNDATION_SAFE_H

#include "CoreFoundationSafeTypes.h"

#include <CoreFoundation/CoreFoundation.h>

CFBooleanRefSafe kCFBooleanFalseSafe(void) {
  return (CFBooleanRefSafe)kCFBooleanFalse;
}

CFBooleanRefSafe kCFBooleanTrueSafe(void) {
  return (CFBooleanRefSafe)kCFBooleanTrue;
}

void CFReleaseSafe(CFTypeRefSafe cf) {
  CFRelease((CFTypeRef)cf);
}

CFTypeID CFGetTypeIDSafe(CFTypeRefSafe cf) {
  return CFGetTypeID((CFTypeRef)cf);
}

CFStringRefSafe CFCopyTypeIDDescriptionSafe(CFTypeID type_id) {
  return (CFStringRefSafe)CFCopyTypeIDDescription(type_id);
}

CFArrayRefSafe CFArrayCreateSafe(CFAllocatorRef allocator, const uintptr_t *values, CFIndex numValues, const CFArrayCallBacks *callBacks) {
  return (CFArrayRefSafe)CFArrayCreate(allocator, (const void **)values, numValues, callBacks);
}

void CFArrayGetValuesSafe(CFArrayRefSafe theArray, CFRange range, const uintptr_t *values) {
  return CFArrayGetValues((CFArrayRef)theArray, range, (const void **)values);
}

CFIndex CFArrayGetCountSafe(CFArrayRefSafe theArray) {
  return CFArrayGetCount((CFArrayRef)theArray);
}

CFDictionaryRefSafe CFDictionaryCreateSafe(CFAllocatorRef allocator, const uintptr_t *keys, const uintptr_t *values, CFIndex numValues, const CFDictionaryKeyCallBacks *keyCallBacks, const CFDictionaryValueCallBacks *valueCallBacks) {
  return (CFDictionaryRefSafe)CFDictionaryCreate(allocator, (const void **)keys, (const void **)values, numValues, keyCallBacks, valueCallBacks);
}

void CFDictionaryGetKeysAndValuesSafe(CFDictionaryRefSafe theDict, const uintptr_t *keys, const uintptr_t *values) {
  return CFDictionaryGetKeysAndValues((CFDictionaryRef)theDict, (const void **)keys, (const void **)values);
}

CFIndex CFDictionaryGetCountSafe(CFDictionaryRefSafe theDict) {
  return CFDictionaryGetCount((CFDictionaryRef)theDict);
}

CFStringRefSafe CFStringCreateWithBytesSafe(CFAllocatorRef alloc, const UInt8 *bytes, CFIndex numBytes, CFStringEncoding encoding, Boolean isExternalRepresentation) {
  return (CFStringRefSafe)CFStringCreateWithBytes(alloc, bytes, numBytes, encoding, isExternalRepresentation);
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

CFNumberType CFNumberGetTypeSafe(CFNumberRefSafe number) {
  return CFNumberGetType((CFNumberRef)number);
}

Boolean CFNumberGetValueSafe(CFNumberRefSafe number, CFNumberType theType, void *valuePtr) {
  return CFNumberGetValue((CFNumberRef)number, theType, valuePtr);
}

Boolean CFBooleanGetValueSafe(CFBooleanRefSafe boolean) {
  return CFBooleanGetValue((CFBooleanRef)boolean);
}

CFDataRefSafe CFDataCreateSafe(CFAllocatorRef allocator, const UInt8 *bytes, CFIndex length) {
  return (CFDataRefSafe)CFDataCreateSafe(allocator, bytes, length);
}

const UInt8 * CFDataGetBytePtrSafe(CFDataRefSafe theData) {
  return CFDataGetBytePtr((CFDataRef)theData);
}

CFIndex CFDataGetLengthSafe(CFDataRefSafe theData) {
  return CFDataGetLength((CFDataRef)theData);
}

CFURLRefSafe CFURLCreateWithFileSystemPathSafe(CFAllocatorRef allocator, CFStringRefSafe filePath, CFURLPathStyle pathStyle, Boolean isDirectory) {
  return (CFURLRefSafe)CFURLCreateWithFileSystemPath(allocator, (CFStringRef)filePath, pathStyle, isDirectory);
}

#endif
