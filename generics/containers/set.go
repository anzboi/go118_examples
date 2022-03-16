package containers

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](init ...T) Set[T] {
	ret := make(Set[T])
	for _, item := range init {
		ret[item] = struct{}{}
	}
	return ret
}

func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

func Union[T comparable](sets ...Set[T]) Set[T] {
	ret := make(Set[T])
	for _, set := range sets {
		for item := range set {
			ret[item] = struct{}{}
		}
	}
	return ret
}

func Intersection[T comparable](sets ...Set[T]) Set[T] {
	ret := make(Set[T])
	if len(sets) == 0 {
		return ret
	}
	base := sets[0]
	for item := range base {
		allHas := true
		for _, set := range sets[1:] {
			if !set.Has(item) {
				allHas = false
				break
			}
		}
		if allHas {
			ret[item] = struct{}{}
		}
	}
	return ret
}

func Minus[T comparable](a, b Set[T]) Set[T] {
	ret := make(Set[T])
	for item := range a {
		if !b.Has(item) {
			ret[item] = struct{}{}
		}
	}
	return ret
}
