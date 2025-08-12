package object

type Enviorment struct {
	store map[string]Object
}

func NewEnv() *Enviorment {
	s := make(map[string]Object)
	return &Enviorment{store: s}
}

func (e *Enviorment) Put(name string, obj Object) Object {
	e.store[name] = obj
	return obj

}

func (e *Enviorment) Get(name string) (Object, bool) {
	v, ok := e.store[name]
	return v, ok

}
