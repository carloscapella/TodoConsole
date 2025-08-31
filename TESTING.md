# Running Tests

To run all unit tests for the project, use the following command from the project root:

```
go test ./...
```

This will recursively run all tests in all packages, including subfolders like `internal/usecase`.

You can also run tests with verbose output:

```
go test -v ./...
```

Or run tests for a specific package:

```
go test ./internal/usecase
```

For more options, see the official Go testing documentation: https://golang.org/pkg/testing/
