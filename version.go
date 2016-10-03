package main

// The git commit that was compiled. This will be filled in by the compiler.
var (
	GitCommit   string
	GitDescribe string
)

// Version contains the current version being run
const Version = "0.0.1"
