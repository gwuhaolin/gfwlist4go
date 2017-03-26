package pac

import (
	"text/template"
	"bytes"
)

type Pac struct {
	BlankList []string
	WhiteList []string
	Proxy     string
}

type templateParams struct {
	HostMap map[string]int
	Proxy   string
}

var (
	tmpl *template.Template
)

func init() {
	tmpl = template.New("")
	tmpl = template.Must(tmpl.Parse(`
var HOST_MAP = {
    {{ range $host, $value := .HostMap }}'{{ $host }}': {{ $value }},
    {{ end }}
};

var PROXY = '{{ .Proxy }}';
var DIRECT = 'DIRECT;';
var PROXY_DIRECT = PROXY + DIRECT;
var DIRECT_PROXY = DIRECT + PROXY;
function proxyForIndex(index) {
    switch (index) {
        case 0:
            return DIRECT;
        case 1:
            return PROXY_DIRECT;
        default:
            return DIRECT_PROXY
    }
}

function FindProxyForURL(_, host) {
    var pos = host.lastIndexOf('.');
    while (true) {
        pos = host.lastIndexOf('.', pos - 1);
        if (pos <= 0) {
            return proxyForIndex(HOST_MAP[host]);
        } else {
            var suffix = host.substring(pos + 1);
            var index = HOST_MAP[suffix];
            if (index !== undefined) {
                return proxyForIndex(index);
            }
        }
    }
}`))
}

func (pac *Pac) String() string {
	hostMap := make(map[string]int, len(pac.BlankList)+len(pac.WhiteList))
	for _, host := range pac.BlankList {
		hostMap[host] = 1
	}
	for _, host := range pac.WhiteList {
		hostMap[host] = 0
	}
	tmplParams := &templateParams{hostMap, pac.Proxy + ";", }
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, tmplParams)
	return buf.String()
}
