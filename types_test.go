package ioc 

import "fmt"

// Interface 1
type I1 interface {
  Get() string
}

type S1 struct {
}

func (self *S1) Get() string {
  return "Hello"
}

// Interface 2, depends on Interface 1
type I2 interface {
  Get2() string
}

type S2 struct {
  Attrib1 I1
}

func (self *S2) Get2() string {
  return self.Attrib1.Get()
}

// Interface 3, depends on Interface 1 and Interface 2
type I3 interface {
  Get3() string
}

type S3 struct {
  Attrib1 I1
  Attrib2 I2
}

func (self *S3) Get3() string {
  return fmt.Sprintf("%s -- %s", self.Attrib1.Get(), self.Attrib2.Get2())
}

// Alternative I3 implementation
type S4 struct {
}

func (self *S4) Get3() string {
  return "Dummy impl"
}

// Container type 
type T1 struct {
  x int
  y int
  S I3 // <-- Notice this is public
}

func NewT1() *T1 {
  return &T1 { x : 1, y : 2 }
}

func (self *T1) test() string {
  return fmt.Sprintf("Output: %s", self.S.Get3())
}
