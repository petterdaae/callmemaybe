package assemblyoutput

import "lang/language/memorymodel"

type procedure struct {
	Name                              string
	NumberOfArgs                      int
	StackSizeBeforeFunctionGeneration int
	Operations                        []string
	ReturnKind                        memorymodel.ContextElementKind
}
