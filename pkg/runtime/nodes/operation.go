package nodes

import (
	"fmt"

	"github.com/odroml/odroml/pkg/runtime"
	"github.com/pkg/errors"
)

// Operation calls an operator on the results of its children as arguments.
type Operation struct {
	BasicNode
	Symbol   string
	ArgCount int
}

// Name returns the name of the AST node.
func (Operation) Name() string {
	return "Operation"
}

// Eval executes the operator by first collecting the parameters and then calling associated function.
func (o *Operation) Eval(c *runtime.Context) (runtime.Value, error) {
	item, err := c.Namespace.Find(runtime.SearchOperator, o.Symbol)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "undefined operator")
	}
	op, ok := item.(runtime.Operator)
	if !ok {
		return runtime.Value{}, errors.Errorf("expected operator, got item %s", item)
	}
	args := make([]runtime.Value, len(o.Childs))
	for i, n := range o.Childs {
		v, err := n.Eval(c)
		if err != nil {
			return runtime.Value{}, errors.Wrap(err, "could not execute operation")
		}
		args[i] = v
	}
	result, err := op.Eval(c, args)
	if err != nil {
		return runtime.Value{}, errors.Wrap(err, "operation "+o.Symbol+" failed")
	}
	c.Behavior = runtime.BehaviorDefault
	return result, nil
}

// NewOperation constructs a new operation of the given symbol.
func NewOperation(symbol string, argCount int, args ...Node) *Operation {
	op := &Operation{
		BasicNode: NewBasic(args...),
		Symbol:    symbol,
		ArgCount:  argCount,
	}
	op.Metadata["label"] = fmt.Sprintf("%s (%d)", symbol, argCount)
	op.Metadata["shape"] = "octagon"
	return op
}
