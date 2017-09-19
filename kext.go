// +build darwin,!ios

package kext

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#include <IOKit/kext/KextManager.h>

typedef uintptr_t CFStringRefSafe;

CFURLRef CFURLCreateWithFileSystemPathSafe(CFAllocatorRef allocator, CFStringRefSafe filePath, CFURLPathStyle pathStyle, Boolean isDirectory) {
  return CFURLCreateWithFileSystemPath(allocator, (CFStringRef)filePath, pathStyle, isDirectory);
}

OSReturn KextManagerLoadKextWithIdentifierSafe(CFStringRefSafe kextIdentifier, CFArrayRef dependencyKextAndFolderURLs) {
  return KextManagerLoadKextWithIdentifier((CFStringRef)kextIdentifier, dependencyKextAndFolderURLs);
}

OSReturn KextManagerUnloadKextWithIdentifierSafe(CFStringRefSafe kextIdentifier) {
  return KextManagerUnloadKextWithIdentifier((CFStringRef)kextIdentifier);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Info struct {
	Version string
	Started bool
}

func LoadInfo(kextID string) (*Info, error) {
	info, err := LoadInfoRaw(kextID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return &Info{
		Version: info["CFBundleVersion"].(string),
		Started: info["OSBundleStarted"].(bool),
	}, nil
}

func LoadInfoRaw(kextID string) (map[interface{}]interface{}, error) {
	cfKextID, err := StringToCFString(kextID)
	if cfKextID != 0 {
		defer ReleaseSafe(CFTypeRefSafe(cfKextID))
	}
	if err != nil {
		return nil, err
	}
	cfKextIDs := ArrayToCFArray([]C.CFTypeRef{C.CFTypeRef(cfKextID)})
	if cfKextIDs != nil {
		defer ReleaseSafe(CFTypeRefSafe(unsafe.Pointer(cfKextIDs)))
	}

	cfDict := C.KextManagerCopyLoadedKextInfo(cfKextIDs, nil)

	m, err := ConvertCFDictionary(cfDict)
	if err != nil {
		return nil, err
	}

	info, hasKey := m[kextID]
	if !hasKey {
		return nil, nil
	}

	var ret, cast = info.(map[interface{}]interface{})
	if !cast {
		return nil, fmt.Errorf("Unexpected value for kext info")
	}

	return ret, nil
}

func Load(kextID string, paths []string) error {
	cfKextID, err := StringToCFString(kextID)
	if cfKextID != 0 {
		defer ReleaseSafe(CFTypeRefSafe(cfKextID))
	}
	if err != nil {
		return err
	}

	var urls []C.CFTypeRef
	for _, p := range paths {
		cfPath, err := StringToCFString(p)
		if cfPath != 0 {
			defer ReleaseSafe(CFTypeRefSafe(cfPath))
		}
		if err != nil {
			return err
		}
		cfURL := C.CFURLCreateWithFileSystemPathSafe(nil, C.CFStringRefSafe(cfPath), 0, 1)
		if cfURL != nil {
			defer ReleaseSafe(CFTypeRefSafe(unsafe.Pointer(cfURL)))
		}

		urls = append(urls, C.CFTypeRef(cfURL))
	}

	cfURLs := ArrayToCFArray(urls)
	if cfURLs != nil {
		defer ReleaseSafe(CFTypeRefSafe(unsafe.Pointer(cfURLs)))
	}

	ret := C.KextManagerLoadKextWithIdentifierSafe(C.CFStringRefSafe(cfKextID), cfURLs)
	if ret != 0 {
		return fmt.Errorf("Error loading kext(%d)", ret)
	}
	return nil
}

func Unload(kextID string) error {
	cfKextID, err := StringToCFString(kextID)
	if cfKextID != 0 {
		defer ReleaseSafe(CFTypeRefSafe(cfKextID))
	}
	if err != nil {
		return err
	}
	ret := C.KextManagerUnloadKextWithIdentifierSafe(C.CFStringRefSafe(cfKextID))
	if ret != 0 {
		return fmt.Errorf("Error unloading kext (%d)", ret)
	}
	return nil
}
