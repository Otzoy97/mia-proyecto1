package cmd

type command interface {
	run()
	validate() bool
	addOpp()
}
