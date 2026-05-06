package baz

import "testing"

type Suite struct{}

func TestValid(t *testing.T) {}

func Testhelper(t *testing.T) {}

func TestNoArg() {}

func TestWrongArg(x int) {}

func TestWrongArity(t *testing.T, x int) {}

func TestWithB(b *testing.B) {}

func (s *Suite) TestMethod(t *testing.T) {}
