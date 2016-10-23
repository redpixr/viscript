package main

import (
	"fmt"
	//"github.com/go-gl/gl/v2.1/gl"
	"math"
	"regexp"
	"strconv"
	//"regexp/syntax"
)

// ISSUES?

/*
only looks for 1 expression per line

no attempt is made to get anything inside a function
which may come after the opening curly brace, on the same line

closing curly brace of function only recognized as a "}" amidst spaces
*/

// NOTES

/*

TODO:
* make sure names are unique within a partic scope
* allow // comments at any position

*/

var types = []string{"bool", "int32", "string"} // FIXME: should allow [] and [42] prefixes
var builtinFuncs = []string{"add32", "sub32", "mult32", "div32"}
var mainBlock = &CodeBlock{Name: "main"} // the root/entry/top/alpha level of the program
var currBlock = mainBlock

// REGEX (raw strings to avoid having to quote backslashes)
var declaredVar = regexp.MustCompile(`^( +)?var( +)?([a-zA-Z]\w*)( +)?int32(( +)?=( +)?([0-9]+))?$`)
var declFuncStart = regexp.MustCompile(`^func ([a-zA-Z]\w*)( +)?\((.*)\)( +)?\{$`)
var declFuncEnd = regexp.MustCompile(`^( +)?\}( +)?$`)
var calledFunc = regexp.MustCompile(`^( +)?([a-zA-Z]\w*)\(([0-9]+|[a-zA-Z]\w*),( +)?([0-9]+|[a-zA-Z]\w*)\)$`)
var comment = regexp.MustCompile(`^//.*`)

type VarBool struct {
	name  string
	value bool
}

type VarInt32 struct {
	name  string
	value int32
}

type VarString struct {
	name  string
	value string
}

type CodeBlock struct {
	Name        string
	VarBools    []VarBool
	VarInt32s   []VarInt32
	VarStrings  []VarString
	CodeBlocks  []*CodeBlock
	Expressions []string
	Parameters  []string // unused atm
}

func initParser() {
	/*
		for _, f := range funcs {
			con.Add(fmt.Sprintf(f))
		}
	*/
	makeHighlyVisibleRuntimeLogHeader(`PARSING`, 5)
	parseAll()
	makeHighlyVisibleRuntimeLogHeader("RUNNING", 5)
	run(mainBlock)
}

func parseAll() {
	for i, line := range rend.Focused.Body {
		parseLine(i, line, false)
	}
}

func parseLine(i int, line string, coloring bool) {
	switch {
	case declaredVar.MatchString(line):
		result := declaredVar.FindStringSubmatch(line)

		if coloring {
			rend.Color(violet)
		} else {
			var s = fmt.Sprintf("%d: var (%s) declared", i, result[3])
			//printIntsFrom(currBlock)

			if result[8] == "" {
				currBlock.VarInt32s = append(currBlock.VarInt32s, VarInt32{result[3], 0})
			} else {
				value, err := strconv.Atoi(result[8])
				if err != nil {
					s = fmt.Sprintf("%s... BUT COULDN'T CONVERT ASSIGNMENT (%s) TO A NUMBER!", s, result[8])
				} else {
					currBlock.VarInt32s = append(currBlock.VarInt32s, VarInt32{result[3], int32(value)})
					s = fmt.Sprintf("%s & assigned: %d", s, value)
				}
			}

			con.Add(fmt.Sprintf("%s\n", s))
		}
	case declFuncStart.MatchString(line):
		result := declFuncStart.FindStringSubmatch(line)

		if coloring {
			rend.Color(fuschia)
		} else {
			con.Add(fmt.Sprintf("%d: func (%s) declared, with params: %s\n", i, result[1], result[3]))

			if currBlock.Name == "main" {
				currBlock = &CodeBlock{Name: result[1]}
				mainBlock.CodeBlocks = append(mainBlock.CodeBlocks, currBlock) // FUTURE FIXME: methods in structs shouldn't be on main/root func
			} else {
				con.Add("Func'y func-ception! CAN'T PUT A FUNC INSIDE A FUNC!\n")
			}
		}
	case declFuncEnd.MatchString(line):
		if coloring {
			rend.Color(fuschia)
		} else {
			con.Add(fmt.Sprintf("func close...\n"))
			//printIntsFrom(mainBlock)
			//printIntsFrom(currBlock)

			if currBlock.Name == "main" {
				con.Add(fmt.Sprintf("ERROR! Main\\Root level function doesn't need enclosure!\n"))
			} else {
				currBlock = mainBlock
			}
		}
	case calledFunc.MatchString(line): // FIXME: hardwired for 2 params each
		result := calledFunc.FindStringSubmatch(line)

		if coloring {
			rend.Color(fuschia)
		} else {
			con.Add(fmt.Sprintf("%d: func call (%s) expressed\n", i, result[2]))
			con.Add(fmt.Sprintf("currBlock: %s\n", currBlock))
			currBlock.Expressions = append(currBlock.Expressions, line)
			/*
				currBlock.Expressions = append(currBlock.Expressions, result[2])
				currBlock.Parameters = append(currBlock.Parameters, result[3])
				currBlock.Parameters = append(currBlock.Parameters, result[5])
			*/
			//printIntsFrom(currBlock)

			/*
				// prints out all captures
				for i, v := range result {
					con.Add(fmt.Sprintf("%d. %s\n", i, v))
				}
			*/
		}
	case comment.MatchString(line): // allow "//" comments    FIXME to allow this at any later point in the line
		if coloring {
			rend.Color(grayDark)
		}
	case line == "":
		// just ignore
	default:
		if coloring {
			rend.Color(white)
		} else {
			con.Add(fmt.Sprintf("SYNTAX ERROR on line %d: \"%s\"\n", i, line))
		}
	}
}

func run(pb *CodeBlock) { // passed block of code
	con.Add(fmt.Sprintf("running function: '%s'\n", pb.Name))

	for i, line := range pb.Expressions {
		con.Add(fmt.Sprintf("running expression: '%s' in function: '%s'\n", line, pb.Name))

		switch {
		case calledFunc.MatchString(line): // FIXME: hardwired for 2 params each
			result := calledFunc.FindStringSubmatch(line)
			con.Add(fmt.Sprintf("%d: calling func (%s) with params: %s, %s\n", i, result[2], result[3], result[5]))

			a := getInt32(result[3])
			if /* not legit num */ a == math.MaxInt32 {
				return
			}
			b := getInt32(result[5])
			if /* not legit num */ b == math.MaxInt32 {
				return
			}

			switch result[2] {
			case "add32":
				con.Add(fmt.Sprintf("%d + %d = %d\n", a, b, a+b))
			case "sub32":
				con.Add(fmt.Sprintf("%d - %d = %d\n", a, b, a-b))
			case "mult32":
				con.Add(fmt.Sprintf("%d * %d = %d\n", a, b, a*b))
			case "div32":
				con.Add(fmt.Sprintf("%d / %d = %d\n", a, b, a/b))
			default:
				for _, fun := range pb.CodeBlocks {
					con.Add((fmt.Sprintf("CodeBlock.Name considered: %s   switching on: %s\n", fun.Name, result[2])))

					if fun.Name == result[2] {
						con.Add((fmt.Sprintf("'%s' matched '%s'\n", fun.Name, result[2])))
						run(fun)
					}
				}
			}
		}
	}
}

func getInt32(s string) int32 {
	value, err := strconv.Atoi(s)

	if err != nil {
		for _, v := range currBlock.VarInt32s {
			if s == v.name {
				return v.value
			}
		}

		if currBlock.Name != "main" {
			for _, v := range mainBlock.VarInt32s {
				if s == v.name {
					return v.value
				}
			}
		}

		con.Add(fmt.Sprintf("ERROR!  '%s' IS NOT A VALID VARIABLE/FUNCTION!\n", s))
		return math.MaxInt32
	}

	return int32(value)
}

func printIntsFrom(f *CodeBlock) {
	if len(f.VarInt32s) == 0 {
		con.Add(fmt.Sprintf("%s has no elements!\n", f.Name))
	} else {
		for i, v := range f.VarInt32s {
			con.Add(fmt.Sprintf("%s.VarInt32s[%d]: %s = %d\n", f.Name, i, v.name, v.value))
		}
	}
}

/*
The FindAllStringSubmatch-function will, for each match, return an array with the
entire match in the first field and the
content of the groups in the remaining fields.
The arrays for all the matches are then captured in a container array.

the number of fields in the resulting array always matches the number of groups plus one.
*/
