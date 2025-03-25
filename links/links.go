package links

type SpbstuLink int

const (
	FacultiesGet = iota
)

var spbStuLinks = map[int]string{
	FacultiesGet: "https://ruz.spbstu.ru/",
}

func (s SpbstuLink) Link() string {
	return spbStuLinks[int(s)]
}
