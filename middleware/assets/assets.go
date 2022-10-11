package assets

import (
	"context"
	"crypto/md5" // nolint: gosec // Only used for asset hashing
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dusted-go/diagnostic/v3/dlog"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

type file struct {
	PhysicalFileName string
	ContentType      string
}

type Bundle struct {
	VirtualFileName string
	Contents        []byte
}

type Middleware struct {
	CSS *Bundle
	JS  *Bundle

	files          map[string]file
	dirPath        string
	cacheDirective string
	devMode        bool
}

func NewMiddleware(
	ctx context.Context,
	dirPath string,
	cacheDirective string,
	devMode bool,
) (*Middleware, error) {
	mw := &Middleware{
		CSS: &Bundle{
			VirtualFileName: "",
			Contents:        []byte{},
		},
		JS: &Bundle{
			VirtualFileName: "",
			Contents:        []byte{},
		},
		files:          map[string]file{},
		dirPath:        dirPath,
		cacheDirective: cacheDirective,
		devMode:        devMode,
	}
	if err := mw.initAssets(ctx, dirPath, devMode); err != nil {
		return nil, err
	}
	return mw, nil
}

func (m *Middleware) Next(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// During development hot reload assets:
		// ---
		ctx := r.Context()
		if m.devMode {
			err := m.initAssets(ctx, m.dirPath, m.devMode)
			if err != nil {
				panic(err)
			}
		}

		path := r.URL.Path

		// CSS request:
		// ---
		if path == m.CSS.VirtualFileName {
			if !m.devMode {
				w.Header().Add("Cache-Control", m.cacheDirective)
			}
			w.Header().Add("Content-Type", "text/css")
			_, err := w.Write(m.CSS.Contents)
			if err != nil {
				dlog.New(ctx).Error().Err(err).Msg("Failed to respond with CSS content.")
			}
			return
		}

		// JS request:
		// ---
		if path == m.JS.VirtualFileName {
			if !m.devMode {
				w.Header().Add("Cache-Control", m.cacheDirective)
			}
			w.Header().Add("Content-Type", "text/javascript")
			_, err := w.Write(m.JS.Contents)
			if err != nil {
				dlog.New(ctx).Error().Err(err).Msg("Failed to respond with JavaScript content.")
			}
			return
		}

		// All other files (images, icons, webmanifests, etc.):
		// ---
		if asset, ok := m.files[path]; ok {
			f, err := os.Open(asset.PhysicalFileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if !m.devMode {
				w.Header().Add("Cache-Control", m.cacheDirective)
			}
			w.Header().Add("Content-Type", asset.ContentType)
			_, err = io.Copy(w, f)
			if err != nil {
				panic(err)
			}
			return
		}

		// If no match then proceed with other middleware:
		// ---
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) initAssets(
	ctx context.Context,
	dirPath string,
	devMode bool,
) error {
	// Bundling:
	// ---
	cssBuilder := strings.Builder{}
	jsBuilder := strings.Builder{}
	files := map[string]file{}

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
					dlog.New(ctx).Debug().Fmt("Indexing key %s", key)
					files[key] = file{
						PhysicalFileName: path,
						ContentType:      "image/svg+xml",
					}
				case strings.HasSuffix(path, ".png"):
					dlog.New(ctx).Debug().Fmt("Indexing key %s", key)
					files[key] = file{
						PhysicalFileName: path,
						ContentType:      "image/png",
					}
				case strings.HasSuffix(path, ".jpg"):
					dlog.New(ctx).Debug().Fmt("Indexing key %s", key)
					files[key] = file{
						PhysicalFileName: path,
						ContentType:      "image/jpg",
					}
				case strings.HasSuffix(path, ".ico"):
					dlog.New(ctx).Debug().Fmt("Indexing key %s", key)
					files[key] = file{
						PhysicalFileName: path,
						ContentType:      "image/x-icon",
					}
				case strings.HasSuffix(path, ".xml"):
					dlog.New(ctx).Debug().Fmt("Indexing key %s", key)
					files[key] = file{
						PhysicalFileName: path,
						ContentType:      "application/xml",
					}
				case strings.HasSuffix(path, ".json") || strings.HasSuffix(path, ".webmanifest"):
					dlog.New(ctx).Debug().Fmt("Indexing key %s", key)
					files[key] = file{
						PhysicalFileName: path,
						ContentType:      "application/json",
					}
				default:
					content, err := os.ReadFile(path)
					if err != nil {
						return fmt.Errorf("error reading asset file '%s': %w", path, err)
					}

					switch {
					case strings.HasSuffix(path, ".css"):
						dlog.New(ctx).Debug().Fmt("Bundling and minifying %s", path)
						cssBuilder.Write(content)
						cssBuilder.WriteString("\n\n")
					case strings.HasSuffix(path, ".js"):
						dlog.New(ctx).Debug().Fmt("Bundling and minifying %s", path)
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
		return fmt.Errorf("error walking filepath '%s': %w", dirPath, err)
	}

	// Minification:
	// ---
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("text/javascript", js.Minify)

	cssString, err := minifier.String("text/css", cssBuilder.String())
	if err != nil {
		return fmt.Errorf("error minifying CSS: %w", err)
	}

	jsString, err := minifier.String("text/javascript", jsBuilder.String())
	if err != nil {
		return fmt.Errorf("error minifying JavaScript: %w", err)
	}

	// Versioning:
	// ---
	cssVersion := "output.dev"
	jsVersion := "output.dev"

	if !devMode {
		// nolint: gosec // Used for checksums only
		hash := md5.New()
		_, err = io.WriteString(hash, cssString)
		if err != nil {
			return fmt.Errorf("error computing MD5 hash of CSS content: %w", err)
		}
		cssVersion = hex.EncodeToString(hash.Sum(nil))

		hash.Reset()
		_, err = io.WriteString(hash, jsString)
		if err != nil {
			return fmt.Errorf("error computing MD5 hash of JavaScript content: %w", err)
		}
		jsVersion = hex.EncodeToString(hash.Sum(nil))
	}

	cssFileName := "/" + cssVersion + ".css"
	jsFileName := "/" + jsVersion + ".js"

	// Return:
	// ---
	m.CSS = &Bundle{
		VirtualFileName: cssFileName,
		Contents:        []byte(cssString),
	}
	m.JS = &Bundle{
		VirtualFileName: jsFileName,
		Contents:        []byte(jsString),
	}
	m.files = files

	return nil
}

// LogFilter is a dlog.Filter which filters logs during asset initialisation.
var LogFilter = dlog.FilterFunc(func(msg string) bool {
	return !(strings.HasPrefix(msg, "HTTP/1.1 ") &&
		(strings.HasSuffix(msg, ".css") ||
			strings.HasSuffix(msg, ".js") ||
			strings.HasSuffix(msg, ".svg") ||
			strings.HasSuffix(msg, ".jpg") ||
			strings.HasSuffix(msg, ".png")))
})
