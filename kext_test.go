// +build darwin

package kext

import (
	"fmt"
	"testing"
)

func TestKextInfo(t *testing.T) {
	info, err := KextInfo("com.github.osxfuse.filesystems.osxfusefs")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", info)
}
