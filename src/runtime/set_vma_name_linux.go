// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package runtime

import (
	"internal/runtime/atomic"
	"internal/runtime/syscall/linux"
	"unsafe"
)

var prSetVMAUnsupported atomic.Bool

// setVMAName calls prctl(PR_SET_VMA, PR_SET_VMA_ANON_NAME, start, len, name)
func setVMAName(start unsafe.Pointer, length uintptr, name string) {
	if debug.decoratemappings == 0 || prSetVMAUnsupported.Load() {
		return
	}

	var sysName [80]byte
	n := copy(sysName[:], " Go: ")
	copy(sysName[n:79], name) // leave final byte zero

	_, _, err := linux.Syscall6(linux.SYS_PRCTL, linux.PR_SET_VMA, linux.PR_SET_VMA_ANON_NAME, uintptr(start), length, uintptr(unsafe.Pointer(&sysName[0])), 0)
	if err == _EINVAL {
		prSetVMAUnsupported.Store(true)
	}
	// ignore other errors
}
