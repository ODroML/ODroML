package parser

import (
	"github.com/odroml/odroml/pkg/lexer/tokens"
	"github.com/odroml/odroml/pkg/runtime/nodes"
	"github.com/pkg/errors"
)

type assignmentParser struct {
	index, size int
	input       []tokens.Token
	assignment  *nodes.Assignment
}

func newAssignmentParser() *assignmentParser {
	return &assignmentParser{
		assignment: nodes.NewMultiAssignment([]string{}),
	}
}

func (ap *assignmentParser) collectAliases() error {
	collectAliases := true
	for ap.index < ap.size && collectAliases {
		active := ap.input[ap.index]
		switch active.Type {
		case tokens.Separator:
			break
		case tokens.Operator:
			if tokens.AssignmentOperator.Match(active.Value) {
				ap.assignment.Operator = active.Value[:len(active.Value)-1]
			} else if active.Value != "=" {
				return errors.Errorf("expected assignment operator, got %s", active.Value)
			}
			collectAliases = false
		case tokens.Identifier:
			ap.assignment.Alias = append(ap.assignment.Alias, active.Value)
		default:
			return errors.Errorf("did not expect token %s", active.Type)
		}
		ap.index++
	}
	if collectAliases {
		return errors.New("expected assignment operator")
	}
	return nil
}

func (ap *assignmentParser) assignValues() error {
	for i := 0; i < len(ap.assignment.Alias); i++ {
		term, n, err := newTermParser().Parse(ap.input[ap.index:])
		if err != nil {
			return err
		}
		ap.index += n
		ap.assignment.AddBack(term)

		if ap.index < ap.size && ap.input[ap.index].Type == tokens.Separator {
			ap.index++
		}
	}
	return nil
}

func (ap *assignmentParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	ap.index, ap.size = 0, len(input)
	ap.input = input

	if err := ap.collectAliases(); err != nil {
		return ap.assignment, ap.index, errors.Wrap(err, "collecting aliases")
	}

	if err := ap.assignValues(); err != nil {
		return ap.assignment, ap.index, errors.Wrap(err, "assigning values")
	}

	return ap.assignment, ap.index, nil
}
