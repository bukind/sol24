package main

import (
	"fmt"
	"os"
	"strconv"
)

type arg interface {
	String() string
	Int() int
	AsArg() string
}

type intArg int

func (i intArg) String() string {
	return fmt.Sprintf("%d", i)
}
func (i intArg) AsArg() string {
	return i.String()
}
func (i intArg) intArg() int {
	return int(i)
}

type exprArg struct {
	Value int
	Str   string
}

func (e exprArg) String() string {
	return e.Str
}
func (e exprArg) AsArg() string {
	return "(" + e.String() + ")"
}
func (e exprArg) intArg() int {
	return e.Value
}

type step []Arg

func (s step) Extract(idx int) (Arg, step) {
	var result step
	for i, x := range a {
		if i != idx {
			result = append(result, x)
		}
	}
	return a[idx], result
}

func (s step) IsDone(want int) (step, bool) {
	if len(s) == 0 {
		return nil, true
	}
	if len(s) == 1 {
		if s[0].intArg() == want {
			return s, true
		}
		return nil, true
	}
	return nil, false
}

func (s step) Clone(a Arg) step {
	return append(step{a}, s...)
}

func (s step) Nextsteps(a, b Arg) []step {
	ai := a.intArg()
	bi := b.intArg()
	as := a.AsArg()
	bs := b.AsArg()
	result := []step{
		s.Clone(exprArg{ai + bi, as + "+" + bs}),
		s.Clone(exprArg{ai * bi, as + "*" + bs}),
	}
	if ai > bi {
		result = append(result, s.Clone(exprArg{ai - bi, as + "-" + bs}))
		if bi != 0 && ai/bi*bi == ai {
			result = append(result, s.Clone(exprArg{ai / bi, as + "/" + bs}))
		}
	} else {
		result = append(result, s.Clone(exprArg{bi - ai, bs + "-" + as}))
		if ai != 0 && bi/ai*ai == bi {
			result = append(result, s.Clone(exprArg{bi / ai, bs + "/" + as}))
		}
	}
	return result
}

func (s step) Solve(want int) step {
	if res, ok := s.IsDone(want); ok {
		return res
	}
	// make a permutation
	for i := 0; i < len(s)-1; i++ {
		arg1, baselist := s.Extract(i)
		for j := i; j < len(baselist); j++ {
			arg2, args := baselist.Extract(j)
			steps := args.Nextsteps(arg1, arg2)
			fmt.Println(steps)
			for _, step := range steps {
				res := step.Solve(want)
				if len(res) == 1 {
					return res
				}
			}
		}
	}
	return nil
}

func main() {
	var args step
	for _, a := range os.Args[1:] {
		arg, err := strconv.ParseUint(a, 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Want positive integer, got %q\n", a)
			os.Exit(1)
		}
		args = append(args, intArg(int(arg)))
	}
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "At least 3 arguments are required\n")
		os.Exit(1)
	}
	init := args[1:]
	fmt.Println("result:", init.Solve(args[0].intArg()))
}
