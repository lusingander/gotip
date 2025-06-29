package tip

type TestFunction struct {
	Name string
	Subs []*SubTest
}

type SubTest struct {
	Name                 string
	Subs                 []*SubTest
	IsUnresolvedSubTests bool
}

type Target struct {
	Path         string
	Name         string
	IsUnresolved bool
}
