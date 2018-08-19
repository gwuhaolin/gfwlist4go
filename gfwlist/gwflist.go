package gfwlist

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var (
	gfwlistUrlList = []string{
		"https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt",
		"https://pagure.io/gfwlist/raw/master/f/gfwlist.txt",
		"http://repo.or.cz/gfwlist.git/blob_plain/HEAD:/gfwlist.txt",
		"https://bitbucket.org/gfwlist/gfwlist/raw/HEAD/gfwlist.txt",
		"https://gitlab.com/gfwlist/gfwlist/raw/master/gfwlist.txt",
		"https://git.tuxfamily.org/gfwlist/gfwlist.git/plain/gfwlist.txt",
	}
)

func fetch() ([]string, error) {
	var res *http.Response
	var err error
	for _, gfwlistUrl := range gfwlistUrlList {
		res, err = http.Get(gfwlistUrl)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	decoder := base64.NewDecoder(base64.StdEncoding, res.Body)
	reader := bufio.NewReader(decoder)
	var list []string
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		list = append(list, string(line))
	}
	return list, err
}

func parse(line string) string {

	/* remove space */
	line = strings.Trim(line, " ")

	if line == "" {
		return ""
	}

	/* ignore ip address */
	if net.ParseIP(line) != nil {
		return ""
	}

	/* ignore pattern */
	if strings.Index(line, ".") == -1 {
		return ""
	}

	/* ignore comment, whitelist, regex */
	if line[0] == '[' ||
		line[0] == '!' ||
		line[0] == '/' ||
		line[0] == '@' {
		return ""
	}

	return gethostname(line)
}

func gethostname(line string) string {
	c := line[0]
	ss := line

	/* replace '*' */
	if strings.Index(ss, "/") == -1 {
		if strings.Index(ss, "*") != -1 && ss[:2] != "||" {
			ss = strings.Replace(ss, "*", "/", -1)
		}
	}

	switch c {
	case '.':
		ss = fmt.Sprintf("http://%s", ss[1:])
	case '|':
		switch ss[1] {
		case '|':
			ss = fmt.Sprintf("http://%s", ss[2:])
		default:
			ss = ss[1:]
		}
	default:
		if strings.HasPrefix(ss, "http") {
			ss = ss
		} else {
			ss = fmt.Sprintf("http://%s", ss)
		}
	}
	ss = strings.Replace(ss, "%2F", "/", -1)

	/* process */
	u, err := url.Parse(ss)
	if err != nil {
		return ""
	}
	host := u.Host
	if n := strings.Index(host, "*"); n != -1 {
		for i := n; i < len(host); i++ {
			if host[i] == '.' {
				host = host[i:]
				break
			}
		}
	}
	return strings.TrimLeft(host, ".")
}

func BlankList() ([]string, error) {
	lineList, err := fetch()
	if err != nil {
		return nil, err
	}
	hostList := []string{}
	for _, line := range lineList {
		host := parse(line)
		if host != "" {
			hostList = append(hostList, host)
		}
	}
	return hostList, nil
}
