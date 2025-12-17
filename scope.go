package main

type Env struct {
	local map[string]any
	outer *Env
}

func NewEnv(params []string, args []any, outer *Env) *Env {
	local := make(map[string]any)
	for i := range params {
		local[params[i]] = args[i]
	}
	return &Env{
		local: local,
		outer: outer,
	}
}

func (e *Env) Define(name string, val any) {
	e.local[name] = val
}

// find the innermost Env in which varName appears.
func (e *Env) Find(varName string) any {
	for key := range e.local {
		if key == varName {
			return e.local[key]
		}
	}
	e.outer.Find(varName)
	return nil
}

func standardEnv() *Env {
	params := []string{
		"#t",
		"#f",
		"eq?",
		"+",
		"*",
		"/",
	}
	args := []any{
		true,
		false,
		equalOp,
		addOp,
		multOp,
		divOp,
	}
	return NewEnv(params, args, nil)
}
