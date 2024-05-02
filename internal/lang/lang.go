package lang

type Languager interface {
	Pkg() []string
	Tools() []string
}
