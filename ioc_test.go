package ioc 

import "testing"

func Test_GivenInstanceFromConstructor_CanRecursivelyResolveProperties(t *testing.T) {
  c := Container{}
  c.Register((*I1)(nil), (*S1)(nil))
  c.Register((*I2)(nil), (*S2)(nil))
  c.Register((*I3)(nil), (*S3)(nil))

  service := NewT1()
  err := c.Resolve(service)
  if err != nil {
    t.Error(err.Error())
  }

  expected := "Output: Hello -- Hello"
  output := service.test()
  if output != expected {
    t.Errorf("Invalid output stream: '%s' != '%s'", output, expected)
  }
}

func Test_GivenInstanceFromImplicitConstructor_CanRecursivelyResolveProperties(t *testing.T) {
  c := Container{}
  c.Register((*I1)(nil), (*S1)(nil))
  c.Register((*I2)(nil), (*S2)(nil))
  c.Register((*I3)(nil), (*S3)(nil))

  service := &T1{}
  err := c.Resolve(service)
  if err != nil {
    t.Error(err.Error())
  }

  expected := "Output: Hello -- Hello"
  output := service.test()
  if output != expected {
    t.Errorf("Invalid output stream: '%s' != '%s'", output, expected)
  }
}

func Test_GivenInstanceFromImplicitConstructorWithBinding_NoPropertiesAreSet(t *testing.T) {
  c := Container{}
  c.Register((*I1)(nil), (*S1)(nil))
  c.Register((*I2)(nil), (*S2)(nil))
  c.Register((*I3)(nil), (*S3)(nil))

  service := &T1{
    S : &S4{},
  }
  err := c.Resolve(service)
  if err != nil {
    t.Error(err.Error())
  }

  expected := "Output: Dummy impl"
  output := service.test()
  if output != expected {
    t.Errorf("Invalid output stream: '%s' != '%s'", output, expected)
  }
}

func Test_GivenInstance_WithInvalidRegister_ResolveFails(t *testing.T) {
  c := Container{}
  c.Register((*I2)(nil), (*S2)(nil))
  c.Register((*I3)(nil), (*S3)(nil))

  service := &T1{}
  err := c.Resolve(service)
  if err == nil {
    t.Error("Did not fail to resolve service as required")
  }
  print(err.Error())
}

func Test_CanDirectlyResolveInterface(t *testing.T) {
  c := Container{}
  c.Register((*I1)(nil), (*S1)(nil))
  c.Register((*I2)(nil), (*S2)(nil))
  c.Register((*I3)(nil), (*S3)(nil))

  i, err := c.Interface((*I3)(nil))
  i3 := i.(I3)
  if err != nil {
    t.Error(err)
  }
  if i3 == nil {
    t.Error("Failed to resolve interface directly")
  }
}

func Test_CantResolveNonInterfaceTypes(t *testing.T) {
  c := Container{}
  c.Register((*I1)(nil), (*S1)(nil))
  c.Register((*I2)(nil), (*S2)(nil))
  c.Register((*I3)(nil), (*S3)(nil))

  _, err := c.Interface((*S3)(nil))
  if err == nil {
    t.Error(err)
  }
}

func Test_CantResolveArbitraryPointerDepths(t *testing.T) {
  c := Container{}
  c.Register((*I1)(nil), (*S1)(nil))
  c.Register((*I2)(nil), (*S2)(nil))
  c.Register((*I3)(nil), (*S3)(nil))

  _, err := c.Interface((**I3)(nil))
  if err == nil {
    t.Error(err)
  }
}

func Test_CantResolveUnboundInterfaces(t *testing.T) {
  c := Container{}
  c.Register((*I1)(nil), (*S1)(nil))
  c.Register((*I2)(nil), (*S2)(nil))

  _, err := c.Interface((*I1)(nil))
  if err != nil { t.Error(err) }

  _, err = c.Interface((*I2)(nil))
  if err != nil { t.Error(err) }

  _, err = c.Interface((*I3)(nil))
  if err == nil { t.Error("Resolved binding that did not exist") }
}
