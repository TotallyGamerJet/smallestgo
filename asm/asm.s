// Setup the parameters to print hello world
// and then call the Kernel to do it.
TEXT _start(SB), 4, $0
	MOVD $1, R0              // 1 = StdOut
	MOVD $helloworld(SB), R1 // string to print
	MOVD $13, R2             // length of our string
	MOVD $4, R16             // Unix write system call
	SVC                      // Call kernel to output the string

	// Setup the parameters to exit the program
	// and then call the kernel to do it.
	MOVD $0, R0  // Use 0 return code
	MOVD $1, R16 // System call number 1 terminates this program
	SVC          // Call kernel to terminate the program

DATA helloworld+0(SB)/8, $"Hello Wo"
DATA helloworld+8(SB)/5, $"rld!\n"
GLOBL helloworld(SB), 8, $13
