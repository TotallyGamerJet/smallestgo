package runtime

import "sync/atomic"

type traceBufPtr uintptr

const (
	traceEvNone              = 0  // unused
	traceEvBatch             = 1  // start of per-P batch of events [pid, timestamp]
	traceEvFrequency         = 2  // contains tracer timer frequency [frequency (ticks per second)]
	traceEvStack             = 3  // stack [stack id, number of PCs, array of {PC, func string ID, file string ID, line}]
	traceEvGomaxprocs        = 4  // current value of GOMAXPROCS [timestamp, GOMAXPROCS, stack id]
	traceEvProcStart         = 5  // start of P [timestamp, thread id]
	traceEvProcStop          = 6  // stop of P [timestamp]
	traceEvGCStart           = 7  // GC start [timestamp, seq, stack id]
	traceEvGCDone            = 8  // GC done [timestamp]
	traceEvGCSTWStart        = 9  // GC STW start [timestamp, kind]
	traceEvGCSTWDone         = 10 // GC STW done [timestamp]
	traceEvGCSweepStart      = 11 // GC sweep start [timestamp, stack id]
	traceEvGCSweepDone       = 12 // GC sweep done [timestamp, swept, reclaimed]
	traceEvGoCreate          = 13 // goroutine creation [timestamp, new goroutine id, new stack id, stack id]
	traceEvGoStart           = 14 // goroutine starts running [timestamp, goroutine id, seq]
	traceEvGoEnd             = 15 // goroutine ends [timestamp]
	traceEvGoStop            = 16 // goroutine stops (like in select{}) [timestamp, stack]
	traceEvGoSched           = 17 // goroutine calls Gosched [timestamp, stack]
	traceEvGoPreempt         = 18 // goroutine is preempted [timestamp, stack]
	traceEvGoSleep           = 19 // goroutine calls Sleep [timestamp, stack]
	traceEvGoBlock           = 20 // goroutine blocks [timestamp, stack]
	traceEvGoUnblock         = 21 // goroutine is unblocked [timestamp, goroutine id, seq, stack]
	traceEvGoBlockSend       = 22 // goroutine blocks on chan send [timestamp, stack]
	traceEvGoBlockRecv       = 23 // goroutine blocks on chan recv [timestamp, stack]
	traceEvGoBlockSelect     = 24 // goroutine blocks on select [timestamp, stack]
	traceEvGoBlockSync       = 25 // goroutine blocks on Mutex/RWMutex [timestamp, stack]
	traceEvGoBlockCond       = 26 // goroutine blocks on Cond [timestamp, stack]
	traceEvGoBlockNet        = 27 // goroutine blocks on network [timestamp, stack]
	traceEvGoSysCall         = 28 // syscall enter [timestamp, stack]
	traceEvGoSysExit         = 29 // syscall exit [timestamp, goroutine id, seq, real timestamp]
	traceEvGoSysBlock        = 30 // syscall blocks [timestamp]
	traceEvGoWaiting         = 31 // denotes that goroutine is blocked when tracing starts [timestamp, goroutine id]
	traceEvGoInSyscall       = 32 // denotes that goroutine is in syscall when tracing starts [timestamp, goroutine id]
	traceEvHeapAlloc         = 33 // gcController.heapLive change [timestamp, heap_alloc]
	traceEvHeapGoal          = 34 // gcController.heapGoal() (formerly next_gc) change [timestamp, heap goal in bytes]
	traceEvTimerGoroutine    = 35 // not currently used; previously denoted timer goroutine [timer goroutine id]
	traceEvFutileWakeup      = 36 // denotes that the previous wakeup of this goroutine was futile [timestamp]
	traceEvString            = 37 // string dictionary entry [ID, length, string]
	traceEvGoStartLocal      = 38 // goroutine starts running on the same P as the last event [timestamp, goroutine id]
	traceEvGoUnblockLocal    = 39 // goroutine is unblocked on the same P as the last event [timestamp, goroutine id, stack]
	traceEvGoSysExitLocal    = 40 // syscall exit on the same P as the last event [timestamp, goroutine id, real timestamp]
	traceEvGoStartLabel      = 41 // goroutine starts running with label [timestamp, goroutine id, seq, label string id]
	traceEvGoBlockGC         = 42 // goroutine blocks on GC assist [timestamp, stack]
	traceEvGCMarkAssistStart = 43 // GC mark assist start [timestamp, stack]
	traceEvGCMarkAssistDone  = 44 // GC mark assist done [timestamp]
	traceEvUserTaskCreate    = 45 // trace.NewContext [timestamp, internal task id, internal parent task id, stack, name string]
	traceEvUserTaskEnd       = 46 // end of a task [timestamp, internal task id, stack]
	traceEvUserRegion        = 47 // trace.WithRegion [timestamp, internal task id, mode(0:start, 1:end), stack, name string]
	traceEvUserLog           = 48 // trace.Log [timestamp, internal task id, key string id, stack, value string]
	traceEvCPUSample         = 49 // CPU profiling sample [timestamp, stack, real timestamp, real P id (-1 when absent), goroutine id]
	traceEvCount             = 50
	// Byte is used but only 6 bits are available for event type.
	// The remaining 2 bits are used to specify the number of arguments.
	// That means, the max event type value is 63.
)

var trace struct {
	// trace.lock must only be acquired on the system stack where
	// stack splits cannot happen while it is held.
	lock          mutex       // protects the following members
	lockOwner     *g          // to avoid deadlocks during recursive lock locks
	enabled       bool        // when set runtime traces events
	shutdown      bool        // set when we are waiting for trace reader to finish after setting enabled to false
	headerWritten bool        // whether ReadTrace has emitted trace header
	footerWritten bool        // whether ReadTrace has emitted trace footer
	shutdownSema  uint32      // used to wait for ReadTrace completion
	seqStart      uint64      // sequence number when tracing was started
	ticksStart    int64       // cputicks when tracing was started
	ticksEnd      int64       // cputicks when tracing was stopped
	timeStart     int64       // nanotime when tracing was started
	timeEnd       int64       // nanotime when tracing was stopped
	seqGC         uint64      // GC start/done sequencer
	reading       traceBufPtr // buffer currently handed off to user
	empty         traceBufPtr // stack of empty buffers
	fullHead      traceBufPtr // queue of full buffers
	fullTail      traceBufPtr
	stackTab      traceStackTable // maps stack traces to unique ids
	// cpuLogRead accepts CPU profile samples from the signal handler where
	// they're generated. It uses a two-word header to hold the IDs of the P and
	// G (respectively) that were active at the time of the sample. Because
	// profBuf uses a record with all zeros in its header to indicate overflow,
	// we make sure to make the P field always non-zero: The ID of a real P will
	// start at bit 1, and bit 0 will be set. Samples that arrive while no P is
	// running (such as near syscalls) will set the first header field to 0b10.
	// This careful handling of the first header field allows us to store ID of
	// the active G directly in the second field, even though that will be 0
	// when sampling g0.
	cpuLogRead *profBuf
	// cpuLogBuf is a trace buffer to hold events corresponding to CPU profile
	// samples, which arrive out of band and not directly connected to a
	// specific P.
	cpuLogBuf traceBufPtr

	reader atomic.Pointer[g] // goroutine that called ReadTrace, or nil

	signalLock  atomic.Uint32 // protects use of the following member, only usable in signal handlers
	cpuLogWrite *profBuf      // copy of cpuLogRead for use in signal handlers, set without signalLock

	// Dictionary for traceEvString.
	//
	// TODO: central lock to access the map is not ideal.
	//   option: pre-assign ids to all user annotation region names and tags
	//   option: per-P cache
	//   option: sync.Map like data structure
	stringsLock mutex
	strings     map[string]uint64
	stringSeq   uint64

	// markWorkerLabels maps gcMarkWorkerMode to string ID.
	markWorkerLabels [len(gcMarkWorkerModeStrings)]uint64

	bufLock mutex       // protects buf
	buf     traceBufPtr // global trace buffer, used when running without a p
}

func traceGoUnpark(gp *g, skip int) {}

func traceGoSysBlock(pp *p) {}

func traceProcStop(pp *p) {}

func traceGCSTWDone()                              {}
func traceGoCreate(newg *g, pc uintptr)            {}
func traceEvent(ev byte, skip int, args ...uint64) {}
func traceReaderAvailable() *g                     { return nil }
func traceGoSysExit(ts int64)                      {}
func traceGoStart()                                {}

//go:systemstack
func traceReader() *g { return nil }

func traceGoPark(traceEv byte, skip int)         {}
func traceGoSched()                              {}
func traceGoPreempt()                            {}
func traceGoSysCall()                            {}
func traceCPUSample(gp *g, pp *p, stk []uintptr) {}

//go:systemstack
func traceProcFree(pp *p) {}

type traceStackTable struct {
	lock mutex // Must be acquired on the system stack
	seq  uint32
	mem  traceAlloc
	tab  [1 << 13]traceStackPtr
}
type traceStackPtr uintptr
type traceAlloc struct {
	head traceAllocBlockPtr
	off  uintptr
}
type traceAllocBlockPtr uintptr

func traceGoEnd()                 {}
func traceGomaxprocs(procs int32) {}
func traceProcStart()             {}
