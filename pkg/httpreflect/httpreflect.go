package httpreflect

type Reflect struct {
	Needle string
}

type HTMLReflect struct {
}

type HeaderReflect struct {
	Header      string
	HeaderValue string
}

func NewReflect(needle string) Reflect {
	return Reflect{}
}
