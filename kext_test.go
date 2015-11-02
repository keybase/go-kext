// +build darwin

package kext

import "testing"

func TestInfoRaw(t *testing.T) {
	info, err := LoadInfoRaw("com.github.osxfuse.filesystems.osxfusefs")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("%v", info)
}

func TestInfo(t *testing.T) {
	info, err := LoadInfo("com.github.osxfuse.filesystems.osxfusefs")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("%v", info)
}

func TestInfoNotFound(t *testing.T) {
	info, err := LoadInfo("not.a.kext")
	if err != nil {
		t.Fatal(err)
	}
	if info != nil {
		t.Fatalf("Should have returned nil")
	}
}

/*
func TestLoad(t *testing.T) {
	err := Load("com.github.osxfuse.filesystems.osxfusefs", []string{"/Library/Filesystems/osxfusefs.fs/Support/osxfusefs.kext"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnload(t *testing.T) {
	err := Unload("com.github.osxfuse.filesystems.osxfusefs")
	if err != nil {
		t.Fatal(err)
	}
}
*/
