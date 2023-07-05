Release Notes
=============

## 4.0.0

- Removed `response` package
- Removed `request` package
- Added `webfile` package
- Renamed `path` package to `route`
- Moved `atom`, `rss` and `sitemap` under `feeds/*`

## 3.18.0

- Added security header middleware and response helper.

## 3.17.0

- Another change to the RSS date format

## 3.16.0

- Removed redundant `textType` parameter from `atom.NewText` and `atom.NewHTML`

## 3.15.0

- `view.NewViewHandler` becomes `view.NewHandler`
- Added `atom.NewHtML` alongside `atom.NewText`

## 3.14.0

- Fixed RSS feed date time strings to format correctly according to the spec.
- Added `atom` package for Atom feeds.
- Moved `server.ViewHandler` to `view.Handler`
- Dissolved `server` package

## 3.13.0

- Fixed double encoding in RSS descriptions.

## 3.12.0

- RSS feeds will generate correct RFC-822 date time strings now.

## 3.11.0

- Added `sitemap` package

## 3.10.0

- Added `rss` package

## 3.9.0

- Moved `response` package into `server` package
- Removed the `ignoreBrokenPipeErr` flags
- Removed redundant `*http.Request` parameter from functions
- Simplified error wrapping

## 3.8.0

- Changed `log` to `dlog` after upgrade.

## 3.7.0

- Upgraded dependencies

## 3.6.0

- Upgraded dependencies

## 3.5.0

- Upgraded dependencies

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