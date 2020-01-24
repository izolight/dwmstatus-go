package dwmstatus

type Refresher interface {
	Refresh() string
}
