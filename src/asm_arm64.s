// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "go_tls.h"
#include "tls_arm64.h"
#include "funcdata.h"
#include "textflag.h"


TEXT runtime·morestack_noctxt(SB),NOSPLIT|NOFRAME,$0-0
	RET

// func memhash32(p unsafe.Pointer, h uintptr) uintptr
TEXT runtime·memhash32<ABIInternal>(SB),NOSPLIT|NOFRAME,$0-24
	RET


