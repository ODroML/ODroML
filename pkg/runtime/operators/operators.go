package operators

import (
	"github.com/odroml/odroml/pkg/runtime"
)

func Load(c *runtime.Context) {
	LoadBasicMath(c)
	LoadLogical(c)
	LoadCompare(c)
}
