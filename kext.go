// +build darwin

package kext

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#include <IOKit/kext/KextManager.h>
*/
import "C"
import "fmt"

type KextInfo struct {
	Version string
	Started bool
}

func KextInfoForLabel(label string) (*KextInfo, error) {
	info, err := KextInfoRaw(label)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return &KextInfo{
		Version: info["CFBundleVersion"].(string),
		Started: info["OSBundleStarted"].(bool),
	}, nil
}

func KextInfoRaw(label string) (map[interface{}]interface{}, error) {
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

	info, hasKey := m[label]
	if !hasKey {
		return nil, nil
	}

	var ret, cast = info.(map[interface{}]interface{})
	if !cast {
		return nil, fmt.Errorf("Unexpected value for kext info")
	}

	return ret, nil
}
