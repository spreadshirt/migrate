# io.fs Migration Driver

This migration driver accepts [io/fs.ReadDirFS](https://pkg.go.dev/io/fs) file system interfaces as migration file source.

Please use `NewIOFSDriver` to instantiate the driver, instantiation with `Open` by URI is not supported.