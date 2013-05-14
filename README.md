DI & IOC in GO
==

Honestly, is there really a need for this?

Who knows, but if you do want to do it, you probably don't want to use 
traditional DI practices in go, it is really not a great fit.

This is a simple module to provide a 'go like' IOC container, based on the
AutoFac property injection framework.


Usage
--

First, you create a container and bind interfaces to implementation structs.

This is not optional. It cannot be controlled at runtime and must be done
at compile time (see FAQ). You probably want it in some global scope.

    c := ioc.Container{}
    c.Register((*I1)(nil), (*S1)(nil))
    c.Register((*I2)(nil), (*S2)(nil))
    c.Register((*I3)(nil), (*S3)(nil))

Then, you create a type that uses various interfaces, eg:

    type MyType struct {
      X int
      Y int 
      Service I1
    }

Now, if you're testing or using this object manually, you can create an instance
as normal, using:

    x := MyType {
      x : 0,
      y : 0,
      Service : &ServiceImpl{ 
        // ... 
      }
    }

Or you can attempt to resolve any 'nil' properties on your instance using the IOC 
container, like this:

    x := MyType {
      x : 0,
      y : 0,
    }
    err := c.Resolve(x)
    if err != nil {
       // ...
    }

Or, if your type comes from another package with a constructor:

    x := package.NewMyType(10, 10)
    c.Resolve(x)

See the tests; resolution is recursive. Ie. If property X requires a Y, then Y is 
created and resolved automatically.


Direct Usage
--

You can directly resolve interfaces using the .Interface() call on the IOC container,
like this:

    x, err := c.Interface((*ServiceType)(nil))

It's almost always not worth doing this, because it involves having to do additional
type assertions and assign the instance to a location manually.


FAQ
--

1) Can I resolve non-interfaces types?

No.

2) Can I get into a dangerous circular dependency chain with this?

Yes. 

3) Should I even use this at all?

That's a more complicated question; it's useful in some circumstances, for example,
injecting a singleton of a database service into all your web controller types.

Is there a better solution? sometimes. 

That doesn't mean you shouldn't use DI just because you're using go.

Just make sure that if you are using it, you're doing it because there's some
benefit in doing so, not just because 'that's what you do'.
