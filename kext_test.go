// +build darwin

package kext

import "testing"

func TestKextInfo(t *testing.T) {
	info, err := KextInfoRaw("com.github.osxfuse.filesystems.osxfusefs")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("%v", info)
}

func TestKextInfoForLabel(t *testing.T) {
	info, err := KextInfoForLabel("com.github.osxfuse.filesystems.osxfusefs")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("%v", info)
}

func TestKextInfoNotFound(t *testing.T) {
	info, err := KextInfoForLabel("not.a.kext")
	if err != nil {
		t.Fatal(err)
	}
	if info != nil {
		t.Fatalf("Should have returned nil")
	}
}
