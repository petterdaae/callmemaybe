package assemblyoutput

type procedure struct {
	Name                              string
	StackSizeBeforeFunctionGeneration int
	NumberOfArgs                      int
	Operations                        []string
}
