package pac

// https://www.barretlee.com/blog/2016/08/25/pac-file/
import (
	"bytes"
	"github.com/gwuhaolin/gfwlist4go/gfwlist"
	"io/ioutil"
	"text/template"
)

type templateParams struct {
	HostMap    map[string]int
	Proxy      string
	DefaultWay string
}

var (
	tmpl *template.Template
)

func init() {
	tmpl = template.New("")
	tmpl = template.Must(tmpl.Parse(`var HOST_MAP = {
{{ range $host, $value := .HostMap }}'{{ $host }}':{{ $value }},{{ end }}
};

var PROXY = '{{ .Proxy }};';
var DIRECT = 'DIRECT;';
var PROXY_DIRECT = PROXY + DIRECT;
var DIRECT_PROXY = DIRECT + PROXY;
function proxyForIndex(val) {
  switch (val) {
    case 0:
      return DIRECT_PROXY;
    case 1:
      return PROXY_DIRECT;
    default:
      return {{ .DefaultWay }};
  }
}

function FindProxyForURL(_, host) {
    if (isPlainHostName(host) ||
        shExpMatch(host, "*.local") ||
        isInNet(dnsResolve(host), "10.0.0.0", "255.0.0.0") ||
        isInNet(dnsResolve(host), "172.16.0.0",  "255.240.0.0") ||
        isInNet(dnsResolve(host), "192.168.0.0",  "255.255.0.0") ||
        isInNet(dnsResolve(host), "127.0.0.0", "255.255.255.0"))
        return "DIRECT";
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
	return {{ .DefaultWay }};
}`))
}

func buildTemplate(proxy, defaultWay string, blankList, whiteList []string) string {
	hostMap := make(map[string]int, len(blankList)+len(whiteList))
	for _, host := range blankList {
		hostMap[host] = 1
	}
	for _, host := range whiteList {
		hostMap[host] = 0
	}
	tmplParams := &templateParams{hostMap, proxy, defaultWay}
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, tmplParams)
	return buf.String()
}

func FetchProxyAutoPac(proxy string, whitelist []string) (string, error) {
	blankList, err := gfwlist.BlankList()
	if err != nil {
		return "", err
	}
	return buildTemplate(proxy, "DIRECT_PROXY", blankList, whitelist), nil
}

func SavePac(proxy string, whitelist []string, filename string) error {
	str, err := FetchProxyAutoPac(proxy, whitelist)
	err = ioutil.WriteFile(filename, []byte(str), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetProxyAllPac(proxy string, whitelist []string) string {
	return buildTemplate(proxy, "PROXY_DIRECT", []string{}, whitelist)
}

func SaveProxyAllPac(proxy string, whitelist []string, filename string) error {
	str := GetProxyAllPac(proxy, whitelist)
	err := ioutil.WriteFile(filename, []byte(str), 0644)
	if err != nil {
		return err
	}
	return nil
}
