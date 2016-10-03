package main

// The git commit that was compiled. This will be filled in by the compiler.
var (
	GitCommit   string
	GitDescribe string
)

// The main version number that is being run at the moment.
const Version = "0.0.1"
