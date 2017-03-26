package gfwlist

import (
	"testing"
)

func TestFetch(t *testing.T) {
	list, err := fetch()
	if err != nil {
		t.Error(err)
	}
	t.Log(len(list))
}

func TestList(t *testing.T) {
	list, err := BlankList()
	if err != nil {
		t.Error(err)
	}
	for _, one := range list {
		t.Log(one)
	}
}
