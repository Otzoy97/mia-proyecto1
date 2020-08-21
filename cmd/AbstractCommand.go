package cmd

//Command ...
type Command interface {
	Run()
	Validate() bool
	//AddOp ...
	AddOp(string, interface{})
}
