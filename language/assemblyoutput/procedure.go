package assemblyoutput

type procedure struct {
	Name                              string
	NumberOfArgs                      int
	StackSizeBeforeFunctionGeneration int
	Operations                        []string
}
