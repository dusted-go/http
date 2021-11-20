Release Notes
=============

## 3.0.0

- Added new `server` package
- Deleted `chain` and `verb` middleware packages.
- The `chain` package has been replaced by `server` and `verb` has been discontinued.
- Refactored all other middleware to adhere to the new server.Middleware interface.
- Renamed `proxy.GetRealIP` middleware to `proxy.ForwardedHeaders`.
- Added a new middleware called `proxy.GetRealIP` with new functionality.

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