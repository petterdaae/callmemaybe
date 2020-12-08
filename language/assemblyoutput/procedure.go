package assemblyoutput

type procedure struct {
	name                              string
	NumberOfArgs                      int
	StackSizeBeforeFunctionGeneration int
	operations                        []string
}
