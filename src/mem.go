package runtime

import "unsafe"

func sysAlloc(n uintptr, sysStat *sysMemStat) unsafe.Pointer { return nil }

func sysMap(v unsafe.Pointer, n uintptr, sysStat *sysMemStat) {}

func sysUsed(v unsafe.Pointer, n, prepared uintptr) {}

func sysFree(v unsafe.Pointer, n uintptr, sysStat *sysMemStat) {}

func sysFault(v unsafe.Pointer, n uintptr) {}

func sysReserve(v unsafe.Pointer, n uintptr) unsafe.Pointer { return nil }

func sysUnused(v unsafe.Pointer, n uintptr) {}
