package main

import (
	"flag"
	"os"
	"testing"
)

var fail bool

func init() {
	flag.BoolVar(&fail, "fail", false, "Specify to make the test fail")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestWithFlags(t *testing.T) {
	if fail {
		t.Errorf("Fail. flag to fail was provided")
	}
}
