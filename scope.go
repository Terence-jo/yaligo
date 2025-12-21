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

func (e *Env) SetVar(name string, val any) {
	e.local[name] = val
}

// find the innermost Env in which varName appears.
func (e *Env) FindEnvWith(varName string) *Env {
	for key := range e.local {
		if key == varName {
			return e
		}
	}
	if e.outer == nil {
		return nil
	}
	return e.outer.FindEnvWith(varName)
}

func (e *Env) FindVar(varName string) any {
	val, ok := e.local[varName]
	if !ok {
		return nil
	}
	return val
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
		&BuiltIn{equalOp},
		&BuiltIn{addOp},
		&BuiltIn{multOp},
		&BuiltIn{divOp},
	}
	return NewEnv(params, args, nil)
}
