package compiler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/franela/goblin"
)

func validateOutput(fname string, inputJSON string) (bool, error) {
	var testObject map[string]interface{}

	if marshalErr := json.Unmarshal([]byte(inputJSON), &testObject); marshalErr != nil {
		return false, marshalErr
	}

	if reflect.DeepEqual(testObject["test"], testObject["result"]) {
		return true, nil
	}

	validErr := fmt.Errorf("Failed to validate test file %v: Test object did not equal expected result object", fname)

	return false, validErr
}

func TestCompileFile(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("CompileFile()", func() {
		g.It("Should parse valid input", func() {
			testFiles := []string{"../test/json/lock_parenting.json"}
			for _, v := range testFiles {
				result, err := CompileFile(v)
				if err != nil {
					g.Fail(err)
				}

				rv, verr := validateOutput(v, result)
				if verr != nil {
					g.Fail(verr)
				}
				g.Assert(len(result) > 0)
				g.Assert(rv)
			}
		})

		g.It("Should error on bad path", func() {
			testFiles := []string{"../test/json/notathing.json"}
			for _, v := range testFiles {
				result, err := CompileFile(v)
				g.Assert(len(result)).Equal(0)
				g.Assert(err != nil)
			}
		})
	})
}
