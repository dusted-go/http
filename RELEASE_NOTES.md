Release Notes
=============

## 3.4.0

- Added flag to control the handling of broken pipe errors

## 3.3.0

- Broken pipe errors won't surface any longer

## 3.2.0

- Updated dependencies
- Updated Go 1.19
- Recoverer func takes now a stack of type `stack.Trace` instead of a `[]byte`

## 3.1.0

- Improved the `assets` middleware to set the CSS and JS output files to a static name when `devMode` is switched on.

## 3.0.0

- Stable v3 release

## 3.0.0-alpha-2

- Added `Plaintext` method to `response` package.
- Added `ViewHandler` to `server` package.

## 3.0.0-alpha-1

- Added new `server` package
- Deleted `chain` and `verb` middleware packages.
- The `chain` package has been replaced by `server` and `verb` has been discontinued.
- Refactored all other middleware to adhere to the new server.Middleware interface.
- Renamed `proxy.GetRealIP` middleware to `proxy.ForwardedHeaders`.
- Added a new middleware called `proxy.GetRealIP` with new functionality.
- Removed the `route` package and moved the `ShiftPath` function to `request`.
- Refactored `assets` package to adhere to the `server.Middleware` interface.
- Updated the `diagnostic/log` package to the `2.0.0-alpha-*` release.
- Added `assets.LogFilter` to filter asset logs.

## 2.1.0

`httptrace` package changes the request object so that previous middleware also has access to tracing data.

## 2.0.0

Modified `ForceHTTPS` function to accept a slice of host strings.

## 1.2.0

Added `route.ShiftPath`

## 1.1.0

Added `redirect.Hosts` middleware to redirect from one host to another (e.g. `www.foo.bar` --> `foo.bar`)

## 1.0.1

Fixed error logging bug in `assets`.

## 1.0.0

Initial release of various packages.