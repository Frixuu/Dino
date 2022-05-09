package dino

// must panics with the passed error, if it is not nil.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustAdd registers a service of type T as a singleton in the provided container.
//
// If the operation fails, this method will panic.
func MustAdd[T any, TImpl any](c *Container) {
	must(Add[T, TImpl](c))
}

// MustAddNamed registers a service of type T as a singleton in the provided container.
//
// If the operation fails, this method will panic.
func MustAddNamed[T any, TImpl any](c *Container, name string) {
	must(AddNamed[T, TImpl](c, name))
}

// MustAddTransient registers a service of type T as a transient in the provided container.
//
// If the operation fails, this method will panic.
func MustAddTransient[T any, TImpl any](c *Container) {
	must(AddTransient[T, TImpl](c))
}

// MustAddTransientNamed registers a service of type T as a transient in the provided container.
//
// If the operation fails, this method will panic.
func MustAddTransientNamed[T any, TImpl any](c *Container, name string) {
	must(AddTransientNamed[T, TImpl](c, name))
}

// AddInstance registers an object of type TImpl as a service of type T
// in the container under a global namespace.
//
// If the operation fails, this method will panic.
func MustAddInstance[T any, TImpl any](c *Container, instance TImpl) {
	must(AddInstance[T](c, instance))
}

// AddInstanceNamed registers an object of type TImpl as a service of type T
// in the container under a provided namespace.
//
// If the operation fails, this method will panic.
func MustAddInstanceNamed[T any, TImpl any](c *Container, name string, instance TImpl) {
	must(AddInstanceNamed[T](c, name, instance))
}

// MustGet tries to create, retrieve or inject an object of type T.
//
// If the operation fails, this method will panic.
func MustGet[T any](c *Container) T {
	svc, err := Get[T](c)
	must(err)
	return svc
}

// MustGetNamed tries to create, retrieve or inject an object of type T.
//
// If the operation fails, this method will panic.
func MustGetNamed[T any](c *Container, name string) T {
	svc, err := GetNamed[T](c, name)
	must(err)
	return svc
}
