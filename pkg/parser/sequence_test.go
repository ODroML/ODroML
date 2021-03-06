package parser

import (
	"testing"

	"github.com/odroml/odroml/pkg/lexer/tokens"
)

func Test_sequenceParser_Parse(t *testing.T) {
	tests := []struct {
		name  string
		input []tokens.Token
		want  int
	}{
		{
			"Single value",
			[]tokens.Token{
				{Type: tokens.Number, Value: "3"},
			},
			1,
		},
		{
			"Statement",
			[]tokens.Token{
				{Type: tokens.Statement},
			},
			1,
		},
		{
			"Empty block",
			[]tokens.Token{
				{Type: tokens.LeftBlock},
				{Type: tokens.RightBlock},
			},
			2,
		},
		{
			"Empty block with semicolon",
			[]tokens.Token{
				{Type: tokens.LeftBlock},
				{Type: tokens.Statement},
				{Type: tokens.RightBlock},
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := newSequenceParser(false, 0)
			_, got, err := sp.Parse(tt.input)
			if err != nil {
				t.Errorf("sequenceParser.Parse() got unexpected error %v", err)
			}
			if got != tt.want {
				t.Errorf("sequenceParser.Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
