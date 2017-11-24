package dataframe

type ByStringAscending []Value

func (vs ByStringAscending) Len() int {
	return len(vs)
}

func (vs ByStringAscending) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ByStringAscending) Less(i, j int) bool {
	vs1, _ := vs[i].String()
	vs2, _ := vs[j].String()
	return vs1 < vs2
}

type ByStringDescending []Value

func (vs ByStringDescending) Len() int {
	return len(vs)
}

func (vs ByStringDescending) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ByStringDescending) Less(i, j int) bool {
	vs1, _ := vs[i].String()
	vs2, _ := vs[j].String()
	return vs1 > vs2
}

type ByFloat64Ascending []Value

func (vs ByFloat64Ascending) Len() int {
	return len(vs)
}

func (vs ByFloat64Ascending) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ByFloat64Ascending) Less(i, j int) bool {
	vs1, _ := vs[i].Float64()
	vs2, _ := vs[j].Float64()
	return vs1 < vs2
}

type ByFloat64Descending []Value

func (vs ByFloat64Descending) Len() int {
	return len(vs)
}

func (vs ByFloat64Descending) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ByFloat64Descending) Less(i, j int) bool {
	vs1, _ := vs[i].Float64()
	vs2, _ := vs[j].Float64()
	return vs1 > vs2
}

type ByDurationAscending []Value

func (vs ByDurationAscending) Len() int {
	return len(vs)
}

func (vs ByDurationAscending) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ByDurationAscending) Less(i, j int) bool {
	vs1, _ := vs[i].Duration()
	vs2, _ := vs[j].Duration()
	return vs1 < vs2
}

type ByDurationDescending []Value

func (vs ByDurationDescending) Len() int {
	return len(vs)
}

func (vs ByDurationDescending) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ByDurationDescending) Less(i, j int) bool {
	vs1, _ := vs[i].Duration()
	vs2, _ := vs[j].Duration()
	return vs1 > vs2
}
