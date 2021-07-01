package assets

import (
	"crypto/md5" // nolint: gosec // Only used for asset hashing
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dusted-go/diagnostic/log"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

// Asset represents a single asset.
type Asset struct {
	FullPath    string
	ContentType string
}

// Bundle represents multiple assets bundled together.
type Bundle struct {
	FileName string
	Contents []byte
}

// Assets contains all static web assets.
type Assets struct {
	CSS   *Bundle
	JS    *Bundle
	Files map[string]Asset
}

var (
	assets *Assets
)

// Get returns an object of all bundled and minified assets.
func Get() *Assets {
	return assets
}

func bundleAndMinify(log log.Event, dirPath string) *Assets {
	// Bundling:
	// ---
	cssBuilder := strings.Builder{}
	jsBuilder := strings.Builder{}
	files := map[string]Asset{}

	err := filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if !info.IsDir() {
				key := "/" + strings.TrimLeft(path[len(dirPath):], "/")

				switch {
				case strings.HasSuffix(path, ".svg"):
					log.Debug().Fmt("Indexing key %s", key)
					files[key] = Asset{
						FullPath:    path,
						ContentType: "image/svg+xml",
					}
				case strings.HasSuffix(path, ".png"):
					log.Debug().Fmt("Indexing key %s", key)
					files[key] = Asset{
						FullPath:    path,
						ContentType: "image/png",
					}
				case strings.HasSuffix(path, ".jpg"):
					log.Debug().Fmt("Indexing key %s", key)
					files[key] = Asset{
						FullPath:    path,
						ContentType: "image/jpg",
					}
				case strings.HasSuffix(path, ".ico"):
					log.Debug().Fmt("Indexing key %s", key)
					files[key] = Asset{
						FullPath:    path,
						ContentType: "image/x-icon",
					}
				case strings.HasSuffix(path, ".xml"):
					log.Debug().Fmt("Indexing key %s", key)
					files[key] = Asset{
						FullPath:    path,
						ContentType: "application/xml",
					}
				case strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".webmanifest"):
					log.Debug().Fmt("Indexing key %s", key)
					files[key] = Asset{
						FullPath:    path,
						ContentType: "application/json",
					}
				default:
					content, err := os.ReadFile(path)
					if err != nil {
						return fmt.Errorf("failed to read file %s with error: %w", path, err)
					}

					switch {
					case strings.HasSuffix(path, ".css"):
						log.Debug().Fmt("Bundling and minifying %s", path)
						cssBuilder.Write(content)
						cssBuilder.WriteString("\n\n")
					case strings.HasSuffix(path, ".js"):
						log.Debug().Fmt("Bundling and minifying %s", path)
						jsBuilder.Write(content)
						jsBuilder.WriteString("\n\n")
					default:
						fmt.Printf("Unsupported file extension found in %s: %s\n", dirPath, path)
					}
				}
			}
			return nil
		})

	if err != nil {
		panic(err)
	}

	// Minification:
	// ---
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/javascript", js.Minify)

	cssString, err := minifier.String("text/css", cssBuilder.String())
	if err != nil {
		panic(fmt.Errorf("failed to minify CSS: %w", err))
	}

	jsString, err := minifier.String("text/javascript", jsBuilder.String())
	if err != nil {
		panic(fmt.Errorf("failed to minify JavaScript: %w", err))
	}

	// Versioning:
	// ---
	// nolint: gosec // Used for checksums only
	hash := md5.New()
	_, err = io.WriteString(hash, cssString)
	if err != nil {
		panic(fmt.Errorf("failed to compute MD5 hash of CSS content: %w", err))
	}
	cssVersion := hex.EncodeToString(hash.Sum(nil))
	cssFileName := "/" + cssVersion + ".css"

	hash.Reset()
	_, err = io.WriteString(hash, jsString)
	if err != nil {
		panic(fmt.Errorf("failed to compute MD5 hash of JavaScript content: %w", err))
	}
	jsVersion := hex.EncodeToString(hash.Sum(nil))
	jsFileName := "/" + jsVersion + ".js"

	// Return:
	// ---
	return &Assets{
		CSS: &Bundle{
			FileName: cssFileName,
			Contents: []byte(cssString),
		},
		JS: &Bundle{
			FileName: jsFileName,
			Contents: []byte(jsString),
		},
		Files: files,
	}
}

// Handle is a function which automatically discovers,
// bundles and minifies static web assets such as CSS
// and JavaScript files and returns a http.HandlerFunc
// to deal with all web requests to those assets.
func Handle(
	logger log.Event,
	dirPath string,
	cacheDirective string,
	notFoundFn http.HandlerFunc,
	debugMode bool) http.HandlerFunc {
	dirPath = strings.Trim(dirPath, "/") + "/"
	assets = bundleAndMinify(logger, dirPath)

	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.Inherit(r.Context())

		if debugMode {
			assets = bundleAndMinify(logger.SetMinLogLevel(log.Emergency), dirPath)
		}

		key := r.URL.Path

		if key == assets.CSS.FileName {
			if !debugMode {
				// 6 month cache duration
				w.Header().Add("Cache-Control", cacheDirective)
			}
			w.Header().Add("Content-Type", "text/css")
			_, err := w.Write(assets.CSS.Contents)
			if err != nil {
				logger.Error().SetError(err).Msg("Failed to respond with CSS content.")
			}

		} else if key == assets.JS.FileName {
			if !debugMode {
				// 6 month cache duration
				w.Header().Add("Cache-Control", cacheDirective)
			}
			w.Header().Add("Content-Type", "text/javascript")
			_, err := w.Write(assets.JS.Contents)
			if err != nil {
				logger.Error().SetError(err).Msg("Failed to respond with JavaScript content.")
			}

		} else if asset, ok := assets.Files[key]; ok {
			f, err := os.Open(asset.FullPath)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if !debugMode {
				// 6 month cache duration
				w.Header().Add("Cache-Control", cacheDirective)
			}
			w.Header().Add("Content-Type", asset.ContentType)
			_, err = io.Copy(w, f)
			if err != nil {
				panic(err)
			}
		} else {
			notFoundFn(w, r)
		}
	}
}
