package exitcode

const (
	// Usage exit because of wrong params
	Usage = iota + 1
	//Panic reserve 2 for panic
	Panic
	// Execution encounter error when executing share or recover
	Execution
	// SrcErr source file or parts problems
	SrcErr
	// DstErr destination output file problems, like recover destination already exists, sharing destination directory not empty
	DstErr
)
