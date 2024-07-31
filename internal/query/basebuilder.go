package query

type BaseBuilder interface {
	//	SetParent(parent interface{})
	Build() (string, []interface{}, error)
}
