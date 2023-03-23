// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "go_tls.h"
#include "tls_arm64.h"
#include "funcdata.h"
#include "textflag.h"


// Windows ARM64 needs an immediate 0xf000 argument.
// See go.dev/issues/53837.
#define BREAK	\
#ifdef GOOS_windows	\
	BRK	$0xf000 	\
#else 				\
	BRK 			\
#endif 				\


TEXT runtime·breakpoint(SB),NOSPLIT|NOFRAME,$0-0
	BREAK
	RET

TEXT runtime·asminit(SB),NOSPLIT|NOFRAME,$0-0
	RET

TEXT runtime·mstart(SB),NOSPLIT|TOPFRAME,$0
	BL	runtime·mstart0(SB)
	RET // not reached

/*
 *  go-routine
 */

// void gogo(Gobuf*)
// restore state from Gobuf; longjmp
TEXT runtime·gogo(SB), NOSPLIT|NOFRAME, $0-8
	RET

TEXT gogo<>(SB), NOSPLIT|NOFRAME, $0
	RET

// void mcall(fn func(*g))
// Switch to m->g0's stack, call fn(g).
// Fn must never return. It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall<ABIInternal>(SB), NOSPLIT|NOFRAME, $0-8
	RET

// systemstack_switch is a dummy routine that systemstack leaves at the bottom
// of the G stack. We need to distinguish the routine that
// lives at the bottom of the G stack from the one that lives
// at the top of the system stack because the one at the top of
// the system stack terminates the stack walk (see topofstack()).
TEXT runtime·systemstack_switch(SB), NOSPLIT, $0-0
	UNDEF
	BL	(LR)	// make sure this function is not leaf
	RET

// func systemstack(fn func())
TEXT runtime·systemstack(SB), NOSPLIT, $0-8
    RET

/*
 * support for morestack
 */

// Called during function prolog when more stack is needed.
// Caller has already loaded:
// R3 prolog's LR (R30)
//
// The traceback routines see morestack on a g0 as being
// the top of a stack (for example, morestack calling newstack
// calling the scheduler calling newm calling gc), so we must
// record an argument size. For that purpose, it has no arguments.
TEXT runtime·morestack(SB),NOSPLIT|NOFRAME,$0-0

	UNDEF

TEXT runtime·morestack_noctxt(SB),NOSPLIT|NOFRAME,$0-0
	RET

// spillArgs stores return values from registers to a *internal/abi.RegArgs in R20.
TEXT ·spillArgs(SB),NOSPLIT,$0-0

	RET

// unspillArgs loads args into registers from a *internal/abi.RegArgs in R20.
TEXT ·unspillArgs(SB),NOSPLIT,$0-0

	RET

// reflectcall: call a function with the given argument list
// func call(stackArgsType *_type, f *FuncVal, stackArgs *byte, stackArgsSize, stackRetOffset, frameSize uint32, regArgs *abi.RegArgs).
// we don't have variable-sized frames, so we use a small number
// of constant-sized-frame functions to encode a few bits of size in the pc.
// Caution: ugly multiline assembly macros in your future!

#define DISPATCH(NAME,MAXSIZE)		\
	MOVD	$MAXSIZE, R27;		\
	CMP	R27, R16;		\
	BGT	3(PC);			\
	MOVD	$NAME(SB), R27;	\
	B	(R27)
// Note: can't just "B NAME(SB)" - bad inlining results.

TEXT ·reflectcall(SB), NOSPLIT|NOFRAME, $0-48
    RET

#define CALLFN(NAME,MAXSIZE)			\
TEXT NAME(SB), WRAPPER, $MAXSIZE-48;		\
	RET

// callRet copies return values back at the end of call*. This is a
// separate function so it can allocate stack space for the arguments
// to reflectcallmove. It does not follow the Go ABI; it expects its
// arguments in registers.
TEXT callRet<>(SB), NOSPLIT, $48-0
	RET

CALLFN(·call16, 16)
CALLFN(·call32, 32)
CALLFN(·call64, 64)
CALLFN(·call128, 128)
CALLFN(·call256, 256)
CALLFN(·call512, 512)
CALLFN(·call1024, 1024)
CALLFN(·call2048, 2048)
CALLFN(·call4096, 4096)
CALLFN(·call8192, 8192)
CALLFN(·call16384, 16384)
CALLFN(·call32768, 32768)
CALLFN(·call65536, 65536)
CALLFN(·call131072, 131072)
CALLFN(·call262144, 262144)
CALLFN(·call524288, 524288)
CALLFN(·call1048576, 1048576)
CALLFN(·call2097152, 2097152)
CALLFN(·call4194304, 4194304)
CALLFN(·call8388608, 8388608)
CALLFN(·call16777216, 16777216)
CALLFN(·call33554432, 33554432)
CALLFN(·call67108864, 67108864)
CALLFN(·call134217728, 134217728)
CALLFN(·call268435456, 268435456)
CALLFN(·call536870912, 536870912)
CALLFN(·call1073741824, 1073741824)

// func memhash32(p unsafe.Pointer, h uintptr) uintptr
TEXT runtime·memhash32<ABIInternal>(SB),NOSPLIT|NOFRAME,$0-24
	RET

// func memhash64(p unsafe.Pointer, h uintptr) uintptr
TEXT runtime·memhash64<ABIInternal>(SB),NOSPLIT|NOFRAME,$0-24

	RET

// func memhash(p unsafe.Pointer, h, size uintptr) uintptr
TEXT runtime·memhash<ABIInternal>(SB),NOSPLIT|NOFRAME,$0-32
	RET

// func strhash(p unsafe.Pointer, h uintptr) uintptr
TEXT runtime·strhash<ABIInternal>(SB),NOSPLIT|NOFRAME,$0-24
	RET

// R0: data
// R1: seed data
// R2: length
// At return, R0 = return value
TEXT aeshashbody<>(SB),NOSPLIT|NOFRAME,$0
	RET

TEXT runtime·procyield(SB),NOSPLIT,$0-0

	RET

// Save state of caller into g->sched,
// but using fake PC from systemstack_switch.
// Must only be called from functions with no locals ($0)
// or else unwinding from systemstack_switch is incorrect.
// Smashes R0.
TEXT gosave_systemstack_switch<>(SB),NOSPLIT|NOFRAME,$0
	RET

// func asmcgocall_no_g(fn, arg unsafe.Pointer)
// Call fn(arg) aligned appropriately for the gcc ABI.
// Called on a system stack, and there may be no g yet (during needm).
TEXT ·asmcgocall_no_g(SB),NOSPLIT,$0-16
	RET

// func asmcgocall(fn, arg unsafe.Pointer) int32
// Call fn(arg) on the scheduler stack,
// aligned appropriately for the gcc ABI.
// See cgocall.go for more details.
TEXT ·asmcgocall(SB),NOSPLIT,$0-20
	RET

// cgocallback(fn, frame unsafe.Pointer, ctxt uintptr)
// See cgocall.go for more details.
TEXT ·cgocallback(SB),NOSPLIT,$24-24
	RET

// Called from cgo wrappers, this function returns g->m->curg.stack.hi.
// Must obey the gcc calling convention.
TEXT _cgo_topofstack(SB),NOSPLIT,$24
	RET

// void setg(G*); set g. for use by needm.
TEXT runtime·setg(SB), NOSPLIT, $0-8
	RET

// void setg_gcc(G*); set g called from gcc
TEXT setg_gcc<>(SB),NOSPLIT,$8
	RET

TEXT runtime·emptyfunc(SB),0,$0-0
	RET

TEXT runtime·abort(SB),NOSPLIT|NOFRAME,$0-0
	MOVD	ZR, R0
	MOVD	(R0), R0
	UNDEF

TEXT runtime·return0(SB), NOSPLIT, $0
	RET

// The top-most function running on a goroutine
// returns to goexit+PCQuantum.
TEXT runtime·goexit(SB),NOSPLIT|NOFRAME|TOPFRAME,$0-0
	RET

// This is called from .init_array and follows the platform, not Go, ABI.
TEXT runtime·addmoduledata(SB),NOSPLIT,$0-0
	RET

TEXT ·checkASM(SB),NOSPLIT,$0-1
	RET

DATA	debugCallFrameTooLarge<>+0x00(SB)/20, $"call frame too large"
GLOBL	debugCallFrameTooLarge<>(SB), RODATA, $20	// Size duplicated below

TEXT runtime·debugCallV2<ABIInternal>(SB),NOSPLIT|NOFRAME,$0-0
    RET

// runtime.debugCallCheck assumes that functions defined with the
// DEBUG_CALL_FN macro are safe points to inject calls.
#define DEBUG_CALL_FN(NAME,MAXSIZE)		\
TEXT NAME(SB),WRAPPER,$MAXSIZE-0;		\
	NO_LOCAL_POINTERS;		\
	MOVD	$0, R20;		\
	BREAK;		\
	MOVD	$1, R20;		\
	BREAK;		\
	RET
DEBUG_CALL_FN(debugCall32<>, 32)
DEBUG_CALL_FN(debugCall64<>, 64)
DEBUG_CALL_FN(debugCall128<>, 128)
DEBUG_CALL_FN(debugCall256<>, 256)
DEBUG_CALL_FN(debugCall512<>, 512)
DEBUG_CALL_FN(debugCall1024<>, 1024)
DEBUG_CALL_FN(debugCall2048<>, 2048)
DEBUG_CALL_FN(debugCall4096<>, 4096)
DEBUG_CALL_FN(debugCall8192<>, 8192)
DEBUG_CALL_FN(debugCall16384<>, 16384)
DEBUG_CALL_FN(debugCall32768<>, 32768)
DEBUG_CALL_FN(debugCall65536<>, 65536)

// func debugCallPanicked(val interface{})
TEXT runtime·debugCallPanicked(SB),NOSPLIT,$16-16
	RET

// Note: these functions use a special calling convention to save generated code space.
// Arguments are passed in registers, but the space for those arguments are allocated
// in the caller's stack frame. These stubs write the args into that stack space and
// then tail call to the corresponding runtime handler.
// The tail call makes these stubs disappear in backtraces.
//
// Defined as ABIInternal since the compiler generates ABIInternal
// calls to it directly and it does not use the stack-based Go ABI.
TEXT runtime·panicIndex<ABIInternal>(SB),NOSPLIT,$0-16
	JMP	runtime·goPanicIndex<ABIInternal>(SB)
TEXT runtime·panicIndexU<ABIInternal>(SB),NOSPLIT,$0-16
	JMP	runtime·goPanicIndexU<ABIInternal>(SB)
TEXT runtime·panicSliceAlen<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R1, R0
	MOVD	R2, R1
	JMP	runtime·goPanicSliceAlen<ABIInternal>(SB)
TEXT runtime·panicSliceAlenU<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R1, R0
	MOVD	R2, R1
	JMP	runtime·goPanicSliceAlenU<ABIInternal>(SB)
TEXT runtime·panicSliceAcap<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R1, R0
	MOVD	R2, R1
	JMP	runtime·goPanicSliceAcap<ABIInternal>(SB)
TEXT runtime·panicSliceAcapU<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R1, R0
	MOVD	R2, R1
	JMP	runtime·goPanicSliceAcapU<ABIInternal>(SB)
TEXT runtime·panicSliceB<ABIInternal>(SB),NOSPLIT,$0-16
	JMP	runtime·goPanicSliceB<ABIInternal>(SB)
TEXT runtime·panicSliceBU<ABIInternal>(SB),NOSPLIT,$0-16
	JMP	runtime·goPanicSliceBU<ABIInternal>(SB)
TEXT runtime·panicSlice3Alen<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R2, R0
	MOVD	R3, R1
	JMP	runtime·goPanicSlice3Alen<ABIInternal>(SB)
TEXT runtime·panicSlice3AlenU<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R2, R0
	MOVD	R3, R1
	JMP	runtime·goPanicSlice3AlenU<ABIInternal>(SB)
TEXT runtime·panicSlice3Acap<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R2, R0
	MOVD	R3, R1
	JMP	runtime·goPanicSlice3Acap<ABIInternal>(SB)
TEXT runtime·panicSlice3AcapU<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R2, R0
	MOVD	R3, R1
	JMP	runtime·goPanicSlice3AcapU<ABIInternal>(SB)
TEXT runtime·panicSlice3B<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R1, R0
	MOVD	R2, R1
	JMP	runtime·goPanicSlice3B<ABIInternal>(SB)
TEXT runtime·panicSlice3BU<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R1, R0
	MOVD	R2, R1
	JMP	runtime·goPanicSlice3BU<ABIInternal>(SB)
TEXT runtime·panicSlice3C<ABIInternal>(SB),NOSPLIT,$0-16
	JMP	runtime·goPanicSlice3C<ABIInternal>(SB)
TEXT runtime·panicSlice3CU<ABIInternal>(SB),NOSPLIT,$0-16
	JMP	runtime·goPanicSlice3CU<ABIInternal>(SB)
TEXT runtime·panicSliceConvert<ABIInternal>(SB),NOSPLIT,$0-16
	MOVD	R2, R0
	MOVD	R3, R1
	JMP	runtime·goPanicSliceConvert<ABIInternal>(SB)
