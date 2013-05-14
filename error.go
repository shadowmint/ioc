package ioc

import "fmt"

type Error struct {
  Message string
}

func fail(msg string, args ...interface{}) *Error {
  return &Error {
    Message : fmt.Sprintf(msg, args...),
  }
}

func (self *Error) Error() string {
  return self.Message
}
