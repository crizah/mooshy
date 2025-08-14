package object

type Enviorment struct {
	store map[string]Object
	outer *Enviorment
}

func SetNewEnclosedEnv(outer *Enviorment) *Enviorment {
	env := NewEnv()
	env.outer = outer
	return env

}

func NewEnv() *Enviorment {
	s := make(map[string]Object)
	return &Enviorment{store: s, outer: nil}
}

func (e *Enviorment) Put(name string, obj Object) Object {
	e.store[name] = obj
	return obj

}

func (e *Enviorment) Get(name string) (Object, bool) {
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		v, ok = e.outer.Get(name)
	}
	return v, ok

}
