package cmd

//Command ...
type Command interface {
	Run()
	Validate()
	//AddOp ...
	AddOp(string, interface{})
}
