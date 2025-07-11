package knowledge

import (
	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/httpreflect"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

type Storage struct {
	Payload        string
	HTMLNode       httpnode.HTMLNode
	HeaderNode     httpnode.HeaderMergeNode
	JSONNode       httpnode.JSONNode
	Response       rhttp.Response
	Request        rhttp.Http
	ExtractBody    extract.Result
	ExtractHeaders extract.Result
	Reflect        httpreflect.Surrounding
	//Request                 output.Request
	HttpReflectSurroundings []httpreflect.Surrounding
}

func NewStorage(result output.Result) Storage {
	return Storage{
		Payload:    result.Payload,
		HTMLNode:   httpnode.GetHTMLNode(result.Response.Body),
		HeaderNode: httpnode.GetHeaderNode(result.Response.Headers),
		Response:   result.Response,

		//Request:  result.Request,
		//HttpReflectSurroundings: ,
	}
}
