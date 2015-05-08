package echo

type ntype uint8
type children []*node

type node struct {
	typ      ntype
	label    byte
	prefix   string
	parent   *node
	children children
	handler  HandlerFunc
	pnames   []string
	echo     *Echo
}

type router struct {
	trees map[string]*node
	echo  *Echo
}

const (
	stype ntype = iota
	ptype
	mtype
)

func NewRouter(e *Echo) (r *router) {
	r = &router{
		trees: make(map[string]*node),
		echo:  e,
	}

	for _, m := range methods {
		r.trees[m] = &node{
			prefix:   "",
			children: children{},
		}
	}

	return
}

func (r *router) Add(method, path string, h HandlerFunc, echo *Echo) {
	var pnames []string

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			r.insert(method, path[:i], nil, stype, nil, echo)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, ptype, pnames, echo)
				return
			}

			r.insert(method, path[:i], nil, ptype, pnames, echo)
		} else if path[i] == '*' {
			r.insert(method, path[:i], nil, stype, nil, echo)
			pnames = append(pnames, "_name")
			r.insert(method, path[:i+1], h, mtype, pnames, echo)
			return
		}
	}

	r.insert(method, path, h, stype, pnames, echo)

}

func (r *router) insert(method, path string, h HandlerFunc, t ntype, pnames []string, echo *Echo) {

	cn := r.trees[method]
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := lcp(search, cn.prefix)

		if l == 0 {
			cn.label = search[0]
			cn.prefix = search
			if h != nil {
				cn.thyp = t
				cn.handler = h
				cn.pnames = pnames
				cn.echo = echo
			}
		} else if l < pl {

		}
	}

}
