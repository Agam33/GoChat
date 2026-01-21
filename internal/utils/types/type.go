package types

type Meta map[string]any

type Pagination struct {
	Limit int
	Page  int
}
