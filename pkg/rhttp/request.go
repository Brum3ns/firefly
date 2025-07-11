package rhttp

import (
	"fmt"
	"strings"

	"github.com/projectdiscovery/rawhttp"
	urlutil "github.com/projectdiscovery/utils/url"
)

type Http struct {
	CoreRequests []HttpRaw
	Client       *Client
	config       Config
}

type HttpRaw struct {
	Method             string
	Url                string
	TemplateRawRequest []byte
}

type Config struct {
	Version             string
	Methods             []string
	Headers             map[string][]string
	Url                 string
	URIPath             string
	Body                string
	RespectHSTS         bool
	HeaderPresetBrowser string
}

func NewRequest(config Config, options *rawhttp.Options) (Http, error) {
	return Http{
		config: config,
		Client: NewClient(options),
	}, nil
}

func (req *Http) GetClient() *Client {
	return req.Client
}

func (req *Http) GetURL() string {
	return req.config.Url
}

func (req *Http) GetCoreRawRequests() []HttpRaw {
	return req.CoreRequests
}

func (req *Http) MakeCoreRequests() error {
	var uripath = req.config.URIPath

	u, err := urlutil.ParseURL(req.config.Url, false)
	if req.config.URIPath == "/" {
		if err != nil {
			return fmt.Errorf("could not parse given url, error : %v", err)
		}
		uripath = u.RequestURI()
	}

	if req.config.HeaderPresetBrowser != "" && !validHeaderPreset(req.config.HeaderPresetBrowser) {
		return fmt.Errorf("could not preset headers, invalid value: [%s]", req.config.HeaderPresetBrowser)
	}

	for _, method := range req.config.Methods {
		rawRequest, err := DumpRequestRaw(
			MakeHTTPVersion(req.config.Version),
			method,
			req.config.Url,
			uripath,
			req.config.Headers,
			req.config.HeaderPresetBrowser,
			strings.NewReader(req.config.Body),
			req.Client.Options,
		)
		if err != nil {
			return fmt.Errorf("failed to create core requests, error : %v", err)
		}
		// Append raw HTTP requests
		httpRequest := HttpRaw{
			Method:             method,
			Url:                req.config.Url,
			TemplateRawRequest: rawRequest,
		}
		req.CoreRequests = append(req.CoreRequests, httpRequest)
	}
	return nil
}
