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

func TestCompileParseFileRefs(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("ParseFileRefs()", func() {
		g.It("Should parse local refs", func() {
			testFiles := []string{"../test/json/extract_dict.json", "../test/json/extract_extvalue.json", "../test/json/extract_list.json", "../test/json/extract_value.json"}
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
	})

	g.Describe("ParseFileRefs()", func() {
		g.It("Should obey @doc resolution directive", func() {
			testFiles := []string{"../test/json/doc_resolve.json"}
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
	})
}

func TestCompileParseFileParents(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("ParseFileParent() ", func() {
		g.It("Should support single parents and parent lists", func() {
			testFiles := []string{"../test/json/grandparent_list1.json", "../test/json/grandparent_list2.json", "../test/json/multiple_inheritance.json"}
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
	})

	g.Describe("ParseFileParent() ", func() {
		g.It("Should preserve locked name regardless of child values", func() {
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
	})
}

func TestCompileParseFileSchema(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("ParseFileSchema() ", func() {
		g.It("Should catch schema validation failures", func() {
			testFiles := []string{"../test/json/nested_schema_fail.json", "../test/json/schema_fail.json"}
			for _, v := range testFiles {
				result, err := CompileFile(v)
				g.Assert(len(result)).Equal(0)
				g.Assert(err != nil)
			}
		})
	})

	g.Describe("ParseFileSchema() ", func() {
		g.It("Should validate files that conform to their schemas", func() {
			testFiles := []string{"../test/json/nested_schema_pass.json", "../test/json/schema_pass.json"}
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
	})
}
