package parser

import (
	"github.com/odroml/odroml/pkg/lexer/tokens"
	"github.com/odroml/odroml/pkg/runtime"
	"github.com/odroml/odroml/pkg/runtime/nodes"
	"github.com/pkg/errors"
)

func newSequenceParser(substitute bool, cap int) *sequenceParser {
	sp := &sequenceParser{substitute: substitute, cap: cap}
	sp.handlers = map[*tokens.Type]func() error{
		tokens.LeftBlock:  sp.handleLeftBlock,
		tokens.Identifier: sp.handleIdentifier,
	}
	return sp
}

type sequenceParser struct {
	substitute       bool
	statement        bool
	index, size, cap int
	sequence         *nodes.Sequence
	active           tokens.Token
	input            []tokens.Token
	handlers         map[*tokens.Type]func() error
}

// handleLeftBlock generates a subsequence that runs in a substitute namespace.
func (sp *sequenceParser) handleLeftBlock() error {
	item, n, err := newSequenceParser(true, 0).Parse(sp.inputSegment(1))
	if err != nil {
		return err
	}
	// ignore both closing and opening block
	sp.index += n + 2
	sp.sequence.AddBack(item)
	sp.statement = false
	return nil
}

// inputSegment generates a slice of the active input sequence.
func (sp *sequenceParser) inputSegment(offset int) []tokens.Token {
	return sp.input[sp.index+offset:]
}

func (sp *sequenceParser) checkForAssignment() bool {
	for i := sp.index; i < sp.size; i++ {
		sp.active = sp.input[i]
		switch sp.active.Type {
		case tokens.Identifier, tokens.Separator:
		case tokens.Operator:
			if !tokens.AssignmentOperator.Match(sp.active.Value) {
				return false
			}
			return true
		default:
			return false
		}
	}
	return false
}

func (sp *sequenceParser) handleIdentifier() error {
	switch sp.active.Value {
	case variableKeyword, constantKeyword:
		stmt, n, err := newDeclarationParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
	case returnKeyword:
		stmt, n, err := newReturnParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
	case breakKeyword:
		sp.sequence.AddBack(nodes.NewController(runtime.BehaviorBreak))
		sp.index++
	case continueKeyword:
		sp.sequence.AddBack(nodes.NewController(runtime.BehaviorContinue))
		sp.index++
	case fallthroughKeyword:
		sp.sequence.AddBack(nodes.NewController(runtime.BehaviorFallthrough))
		sp.index++
	case functionKeyword:
		stmt, n, err := newFunctionParser(false).Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
		sp.statement = false
	case operatorKeyword:
		stmt, n, err := newOperatorParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
		sp.statement = false
	case ifKeyword:
		stmt, n, err := newBranchParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
		sp.statement = false
	case forKeyword:
		stmt, n, err := newLoopParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
		sp.statement = false
	case matchKeyword:
		stmt, n, err := newMatchParser().Parse(sp.inputSegment(0))
		if err != nil {
			return err
		}
		sp.sequence.AddBack(stmt)
		sp.index += n
		sp.statement = false
	default:
		if sp.checkForAssignment() {
			stmt, n, err := newAssignmentParser().Parse(sp.inputSegment(0))
			if err != nil {
				return err
			}
			sp.sequence.AddBack(stmt)
			sp.index += n
		} else {
			return sp.handleTerm()
		}
	}
	return nil
}

func (sp *sequenceParser) handleTerm() error {
	term, n, err := newTermParser().Parse(sp.inputSegment(0))
	if err != nil {
		return err
	}
	if term != nil {
		sp.sequence.AddBack(term)
	}
	sp.index += n
	return nil
}

func (sp *sequenceParser) Parse(input []tokens.Token) (nodes.Node, int, error) {
	sp.index, sp.size = 0, len(input)
	sp.sequence = nodes.NewSequence(sp.substitute)
	sp.input = input
	for sp.index < sp.size {
		sp.statement = true
		sp.active = sp.input[sp.index]

		//fmt.Printf("[%d:%d] %s\n", sp.index, sp.size, sp.active)
		switch sp.active.Type {
		case nil, tokens.RightBlock:
			return sp.sequence, sp.index, nil
		default:
			handler, ok := sp.handlers[sp.active.Type]
			if !ok {
				handler = sp.handleTerm
			}
			if err := handler(); err != nil {
				return sp.sequence, sp.index, errors.Wrapf(err, "failed handling token %s", sp.active.Type.Name)
			}
		}
		if sp.cap != 0 && len(sp.sequence.Childs) >= sp.cap {
			return sp.sequence, sp.index, nil
		}
		if sp.index < sp.size && sp.statement {
			if sp.input[sp.index].Type != tokens.Statement {
				return sp.sequence, sp.index, errors.New("expected end statement")
			}
			sp.index++
		}
	}
	return sp.sequence, sp.index, nil
}
