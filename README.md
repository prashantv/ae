# AppError (ae)
AppError (ae) is a package to wrap Go's `error` type with extra information for larger applications. AppErrors can annotate take an underlying `error` and add:
* Stack trace information
* Optional list of errors with additional information (without losing the original error).

ae can be used as a replacement for the `errors` package and for `fmt.Errorf`.

## Quick Start
Instead of returning `error`, start returning `ae.Err` from your functions.

To convert a standard `error` to  an `ae.Err`, call `ae.Wrap`.

You can add additional information when wrapping the error by using `ae.Wrapf`.

To create a new error, you can either:
 * Use `ae.Errorf` to create an error using a string (similar to fmt.Errorf).
 * Create a standard error struct and wrap it using `ae.Wrap`.

## Benefits

 * When returning errors up the stack in a large application, you can annotate errors without losing the original error by wrapping the error with more information.
 * Get a stack trace of where exactly the original error was caused.

## Testing

In unit tests, you can use aetest.
