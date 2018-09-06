package main

import (
	"fmt"
	"os"
	"strconv"
)

type Arg interface {
	String() string
	Int() int
	AsArg() string
}

type Int int

func (i Int) String() string {
	return fmt.Sprintf("%d", i)
}
func (i Int) AsArg() string {
	return i.String()
}
func (i Int) Int() int {
	return int(i)
}

type Expr struct {
	Value int
	Str   string
}

func (e Expr) String() string {
	return e.Str
}
func (e Expr) AsArg() string {
	return "(" + e.String() + ")"
}
func (e Expr) Int() int {
	return e.Value
}

type Step []Arg

func (a Step) Extract(idx int) (Arg, Step) {
	var result Step
	for i, x := range a {
		if i != idx {
			result = append(result, x)
		}
	}
	return a[idx], result
}

func (s Step) IsDone(want int) (Step, bool) {
	if len(s) == 0 {
		return nil, true
	}
	if len(s) == 1 {
		if s[0].Int() == want {
			return s, true
		}
		return nil, true
	}
	return nil, false
}

func (s Step) Clone(a Arg) Step {
	return append(Step{a}, s...)
}

func (s Step) NextSteps(a, b Arg) []Step {
	ai := a.Int()
	bi := b.Int()
	as := a.AsArg()
	bs := b.AsArg()
	result := []Step{
		s.Clone(Expr{ai + bi, as + "+" + bs}),
		s.Clone(Expr{ai * bi, as + "*" + bs}),
	}
	if ai > bi {
		result = append(result, s.Clone(Expr{ai - bi, as + "-" + bs}))
		if bi != 0 && ai/bi*bi == ai {
			result = append(result, s.Clone(Expr{ai / bi, as + "/" + bs}))
		}
	} else {
		result = append(result, s.Clone(Expr{bi - ai, bs + "-" + as}))
		if ai != 0 && bi/ai*ai == bi {
			result = append(result, s.Clone(Expr{bi / ai, bs + "/" + as}))
		}
	}
	return result
}

func (s Step) Solve(want int) Step {
	if res, ok := s.IsDone(want); ok {
		return res
	}
	// make a permutation
	for i := 0; i < len(s)-1; i++ {
		arg1, baselist := s.Extract(i)
		for j := i; j < len(baselist); j++ {
			arg2, args := baselist.Extract(j)
			steps := args.NextSteps(arg1, arg2)
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
	var args Step
	for _, a := range os.Args[1:] {
		arg, err := strconv.ParseUint(a, 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Want positive integer, got %q\n", a)
			os.Exit(1)
		}
		args = append(args, Int(int(arg)))
	}
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "At least 3 arguments are required\n")
		os.Exit(1)
	}
	init := args[1:]
	fmt.Println("result:", init.Solve(args[0].Int()))
}
