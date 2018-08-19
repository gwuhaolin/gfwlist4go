package pac

import (
	"bytes"
	"github.com/gwuhaolin/gfwlist4go/gfwlist"
	"io/ioutil"
	"text/template"
)

type Pac struct {
	BlankList []string
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
	tmpl = template.Must(tmpl.Parse(`var HOST_MAP = {
{{ range $host, $value := .HostMap }}'{{ $host }}':{{ $value }},{{ end }}
};

var PROXY = '{{ .Proxy }}';
var PROXY_DIRECT = PROXY + ';DIRECT';
var DIRECT_PROXY = 'DIRECT;' + PROXY;
function proxyForIndex(val) {
	if(val){
		return PROXY_DIRECT;
	}else{
		return DIRECT_PROXY;
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
	return DIRECT_PROXY;
}`))
}

func FetchPac(proxy string) (string, error) {
	blankList, err := gfwlist.BlankList()
	if err != nil {
		return "", err
	}
	doc := Pac{
		BlankList: blankList,
		Proxy:     proxy,
	}
	hostMap := make(map[string]int, len(doc.BlankList))
	for _, host := range doc.BlankList {
		hostMap[host] = 1
	}
	tmplParams := &templateParams{hostMap, doc.Proxy}
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, tmplParams)
	return buf.String(), nil
}

func SavePac(proxy string, filename string) error {
	str, err := FetchPac(proxy)
	err = ioutil.WriteFile(filename, []byte(str), 0644)
	if err != nil {
		return err
	}
	return nil
}
