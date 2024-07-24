package main

import (
	"go/ast"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestCrawlAPI(t *testing.T) {
	for _, test := range []struct {
		name    string
		apiRoot string
		have    []string
		expect  map[string][]string
	}{
		{
			"empty apiRoot returns **/*.go",
			"",
			[]string{
				"one.go",
				"two.go",
				"miss.js",
				"hit/three.go",
			},
			map[string][]string{
				".":   {"one.go", "two.go"},
				"hit": {"three.go"},
			},
		},
		{
			"apiRoot returns apiRoot/*.go",
			"api",
			[]string{
				"api/hit.go",
				"api/miss.js",
				"notapi/miss.go",
			},
			map[string][]string{
				"api": {"hit.go"},
			},
		},
		{
			"apiRoot returns apiRoot/**/*.go",
			"api",
			[]string{
				"api/entity1/hit.go",
				"api/entity1/ya.go",
				"api/entity1/na.goe",
				"api/entity2/hit.go",
				"api/miss.js",
				"notapi/miss.go",
			},
			map[string][]string{
				"api/entity1": {"hit.go", "ya.go"},
				"api/entity2": {"hit.go"},
			},
		},
		{
			"skips test files",
			"api",
			[]string{
				"api/hit.go",
				"api/miss_test.go",
			},
			map[string][]string{
				"api": {"hit.go"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			fs := make(fstest.MapFS, len(test.have))
			for _, file := range test.have {
				fs[file] = &fstest.MapFile{Data: []byte{}}
			}
			got, err := crawlAPI(fs, test.apiRoot)
			require.NoError(err)
			require.EqualValues(test.expect, got)
		})
	}
}

func TestComputeImplPackagePath(t *testing.T) {
	for _, test := range []struct {
		name                              string
		apiRoot, implRoot, apiPackagePath string
		expect                            string
	}{
		{
			"canonicalizes paths; empty apiRoot",
			"",
			"internal",
			"./movies/v1",
			"internal/movies/v1",
		},
		{
			"non-empty apiRoot and implRoot",
			"api",
			"internal",
			"api/movies/v1",
			"internal/movies/v1",
		},
		{
			"arbitrarily nested paths",
			"api",
			"internal",
			"api/movies/cinema/tv",
			"internal/movies/cinema/tv",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := computeImplPackagePath(test.apiRoot, test.implRoot, test.apiPackagePath)
			require.Equal(t, test.expect, got)
		})
	}
}

func TestGetModule(t *testing.T) {
	for _, test := range []struct {
		name   string
		fsys   map[string]string
		root   string
		expect string
	}{
		{
			"go.mod in root directory",
			map[string]string{
				"go.mod": `
        module example

        go 1.22.1
        `,
			},
			".",
			"example",
		},
		{
			"go.mod in parent directory",
			map[string]string{
				"go.mod": `
        module example

        go 1.22.1
        `,
				"api/v1/movies/test.go": ``,
				"api/test.go":           ``,
			},
			"api/v1/movies",
			"example",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			fsys := make(fstest.MapFS, len(test.fsys))
			for path, content := range test.fsys {
				fsys[path] = &fstest.MapFile{Data: []byte(content), Mode: 0644}
			}
			got, err := getModule(
				fsys,
				test.root,
			)
			require.NoError(err)
			require.Equal(test.expect, got)
		})
	}
}

func TestLoadLocalPackage(t *testing.T) {
	for _, test := range []struct {
		name                    string
		fsys                    map[string]string
		astFile                 *ast.File
		have                    string
		expectPath, expectAlias string
	}{
		{
			"correctly prefixes module",
			map[string]string{
				"go.mod": `
        module example

        go 1.22.1
        `,
			},
			&ast.File{},
			"api/v1/movies",
			"example/api/v1/movies",
			"",
		},
		{
			"removes _ in alias",
			map[string]string{
				"go.mod": `
        module example

        go 1.22.1
        `,
			},
			&ast.File{},
			"api/v1/movies_list",
			"example/api/v1/movies_list",
			"movieslist",
		},
		{
			"uses existing alias if set",
			map[string]string{
				"go.mod": `
        module example

        go 1.22.1
        `,
			},
			&ast.File{
				Imports: []*ast.ImportSpec{
					{
						Name: &ast.Ident{Name: "moovies"},
						Path: &ast.BasicLit{Value: `"example/api/v1/movies_list"`},
					},
				},
			},
			"api/v1/movies_list",
			"example/api/v1/movies_list",
			"moovies",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			fsys := make(fstest.MapFS, len(test.fsys))
			for path, content := range test.fsys {
				fsys[path] = &fstest.MapFile{Data: []byte(content), Mode: 0644}
			}
			path, alias, err := loadLocalPackage(
				fsys,
				test.astFile,
				test.have,
			)
			require.NoError(err)
			require.Equal(test.expectPath, path)
			require.Equal(test.expectAlias, alias)
		})
	}
}
