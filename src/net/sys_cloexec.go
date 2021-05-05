// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements sysSocket for platforms that do not provide a fast path
// for setting SetNonblock and CloseOnExec.

//go:build aix || darwin || (solaris && !illumos)
// +build aix darwin solaris,!illumos

package net

import (
	"internal/poll"
	"os"
	"syscall"
)

// Wrapper around the socket system call that marks the returned file
// descriptor as nonblocking and close-on-exec.
func sysSocket(family, sotype, proto int) (int, error) { // 这个和老版本有区别了
	// See ../syscall/exec_unix.go for description of ForkLock.
	syscall.ForkLock.RLock()
	s, err := socketFunc(family, sotype, proto) // 通知将socket进行阻塞
	if err == nil {
		syscall.CloseOnExec(s)
	}


	// 这里是内核版本低于2.6.27时，代码会走到这里 ,下面的代码防止描述符溢出
	syscall.ForkLock.RUnlock()
	if err != nil {
		return -1, os.NewSyscallError("socket", err)
	}
	if err = syscall.SetNonblock(s, true); err != nil {
		poll.CloseFunc(s)
		return -1, os.NewSyscallError("setnonblock", err)
	}
	return s, nil
}
