package service

import (
	"log"
	"testing"
)

// test load

type OuterModule struct {
	Inner *InnerModule
}

func (o *OuterModule) Init(c *Config) {}

type InnerModule struct {
	TestOption  bool
	DidRunStart bool
}

func (i *InnerModule) Init(c *Config) {
	c.Start = func() {
		i.DidRunStart = true
	}
}

func TestServiceLoad(t *testing.T) {
	BootPrintln = log.Println
	tm, stop := New(&OuterModule{
		Inner: &InnerModule{TestOption: true},
	}).StartForTest()
	defer stop()
	if tm.Inner.TestOption != true {
		t.Fatal("Expected Inner.TestOption to be true")
	}
	if tm.Inner.DidRunStart != true {
		t.Fatal("Expected Inner.DidRunStart to be true")
	}
}
