package runtime

type semaProfileFlags int

func semacquire1(addr *uint32, lifo bool, profile semaProfileFlags, skipframes int, reason waitReason) {
}
func semacquire(addr *uint32)                                {}
func semrelease(addr *uint32)                                {}
func semrelease1(addr *uint32, handoff bool, skipframes int) {}
