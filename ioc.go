package ioc 

import "reflect"
import "fmt"

type Binding struct {
  Type reflect.Type
  Impl reflect.Type
  Instance reflect.Value
  Ready bool
}

type Container struct {
  items map[reflect.Type]*Binding
}

func (self *Container) resolveType(i interface{}) (t reflect.Type) {
  defer func() {
    if recover() != nil {
      t = nil
    }
  }()
  raw := reflect.TypeOf(i)
  t = raw.Elem()
  return t
}

func (self *Container) createInstance(t reflect.Type) (reflect.Value, error) {
  i := reflect.New(t)
  r := i.Interface()
  err := self.Resolve(r)
  return i, err
}

// Invoke this using pointers to types, like this:
// Register((*IType)(nil), (*TType)(nil))
func (self *Container) Register(i interface{}, t interface{}) error {
  
  if self.items == nil {
    self.items = map[reflect.Type]*Binding {}
  }

  key := self.resolveType(i)
  if key == nil {
    return fail("Invalid interface type (use (*T)(nil))")
  }

  impl := self.resolveType(t)
  if impl == nil {
    return fail("Invalid concrete type (use (*T)(nil))")
  }

  item, found := self.items[key]
  if found {
    return fail("Duplicate binding for type '%s'", key.Name())
  }

  item = &Binding{
    Type : key,
    Impl : impl,
    Ready : false,
  }
  
  self.items[key] = item
  return nil
}

// Attempt to resolve any interface hooks on the given instance
// that are currently nil, using the registered service. If the
// instance of the given service does not exist yet, create it.
//
// This is a recursive function that resolves the new service
// after creating it; cyclic loops are entirely possible.
func (self *Container) Resolve(target interface{}) error {

  if self.items == nil {
    self.items = map[reflect.Type]*Binding {}
  }

  t := reflect.TypeOf(target)
  v := reflect.ValueOf(target)
  for t.Kind() == reflect.Ptr {
    t = t.Elem()
    v = v.Elem()
  }

  pc := t.NumField()
  for i := 0; i < pc; i++ {
    field := t.Field(i)
    if field.Type.Kind() == reflect.Interface {

      fmt.Printf("Located field %s\n", field.Name)
      value := v.Field(i)
      if value.IsNil() {

        fmt.Printf("Located nil field value! %s\n", field.Name)
        record, found := self.items[field.Type]
        if !found {
          return fail("Unable to resolve property, no binding for '%s' on type '%s'", field.Type.Name(), t.Name())
        } else {

          if !record.Ready {
            inst, err := self.createInstance(record.Impl)
            if err != nil {
              return err
            }
            record.Instance = inst
            record.Ready = true
          }

          // Bind the instance to the record, if we can~
          fmt.Printf("Attempting to bind instance %+v\n", record.Instance)
          fmt.Printf("... to object %s property %s of type %+v\n", t.Name(), field.Name, field.Type)
          value.Set(record.Instance)
        }
      }
    }
  }
  return nil
}

// Resolve an interface directly, note this is not the preferred way to
// obtain an interface as the result must be type-asserted back to the
// required type.
func (self *Container) Interface(target interface{}) (interface{}, error) {

  if self.items == nil {
    self.items = map[reflect.Type]*Binding {}
  }

  t := reflect.TypeOf(target)
  tt := t.Elem()

  record, found := self.items[tt]
  if !found {
    return nil, fail("Unable to resolve interface, no binding for '%s'", tt.Name())
  } 

  if !record.Ready {
    inst, err := self.createInstance(record.Impl)
    if err != nil {
      return nil, err
    }
    record.Instance = inst
    record.Ready = true
  }

  return record.Instance.Interface(), nil
}
