package service

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	for _, x := range []struct {
		args string
		f, g *string
	}{
		{"-f", ptr(""), nil},
		{"--f", ptr(""), nil},
		{"-f foo", ptr("foo"), nil},
		{"--f foo", ptr("foo"), nil},
		{"-f --g", ptr(""), ptr("")},
		{"--f -g", ptr(""), ptr("")},
		{"-f foo -g", ptr("foo"), ptr("")},
		{"--f foo -g", ptr("foo"), ptr("")},
		{"-f --g bar", ptr(""), ptr("bar")},
	} {
		f := &Flag{Key: "--f"}
		g := &Flag{Key: "-g,--gee"}
		m, err := parseArgs([]*Flag{f, g}, strings.Split(x.args, " "))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(x.f, f.Value) {
			t.Fatalf("%s %#v %#v", x.args, x.f, f.Value)
		}
		if !reflect.DeepEqual(x.g, g.Value) {
			t.Fatalf("%s %#v %#v", x.args, x.g, g.Value)
		}
		if !reflect.DeepEqual(m["g"], m["gee"]) {
			t.Fatal(m["g"], m["gee"])
		}
	}
}
