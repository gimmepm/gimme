# Integration Tests

### Run all tests

From the root of the repo...

```
$ make test
```

### Add an integration test

Add a test to this `integration` dir that matches `test*.sh`. Make the test executable (`chmod 755 <script_file>`)

### Which tests are run?

All tests that match `integration/test*.sh` are run when `make test` is invoked
