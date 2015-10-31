// +build darwin

package kext

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#include <IOKit/kext/KextManager.h>
*/
import "C"

func KextInfo(label string) (map[interface{}]interface{}, error) {
	cfLabel, err := StringToCFString(label)
	if cfLabel != nil {
		defer Release(C.CFTypeRef(cfLabel))
	}
	if err != nil {
		return nil, err
	}
	cfLabels := ArrayToCFArray([]C.CFTypeRef{C.CFTypeRef(cfLabel)})
	if cfLabels != nil {
		defer Release(C.CFTypeRef(cfLabels))
	}
	cfDict := C.KextManagerCopyLoadedKextInfo(cfLabels, nil)

	m, err := ConvertCFDictionary(cfDict)
	if err != nil {
		return nil, err
	}

	return m, nil
}
