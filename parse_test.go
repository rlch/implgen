package main

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestParseRepositoriesForPackage(t *testing.T) {
	for _, test := range []struct {
		name        string
		fsys        map[string]string
		packagePath string
		files       []string
		expect      []*Repository
	}{
		{
			"empty fsys returns empty slice",
			map[string]string{},
			"api",
			[]string{},
			[]*Repository{},
		},
		{
			"multiple repositories across multiple files",
			map[string]string{
				"api/one.go": `
        package api

        type ARepository interface { A() }
        `,
				"api/two.go": `
        package api

        type BRepository interface { B() }
        `,
			},
			"api",
			[]string{"one.go", "two.go"},
			[]*Repository{
				{
					Package:  "api",
					Filename: "one.go",
					Ident:    "ARepository",
					Methods:  []*Method{{Ident: "A"}},
				},
				{
					Package:  "api",
					Filename: "two.go",
					Ident:    "BRepository",
					Methods:  []*Method{{Ident: "B"}},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			ctx := context.Background()
			fsys := make(fstest.MapFS, len(test.fsys))
			for path, content := range test.fsys {
				fsys[path] = &fstest.MapFile{Data: []byte(content), Mode: 0644}
			}
			got, err := parseRepositoriesForPackage(
				ctx,
				fsys,
				test.packagePath,
				test.files,
			)
			require.NoError(err)
			testRepositories(t, test.expect, got)
		})
	}
}

func TestParseRepositories(t *testing.T) {
	for _, test := range []struct {
		name   string
		src    string
		expect []*Repository
	}{
		{
			"Repository returns empty slice",
			`
      package main

      type Foo struct {}
      `,
			[]*Repository{},
		},
		{
			"Repository returns slice with Repository",
			`
      package main

      type ARepository interface { }
      `,
			[]*Repository{
				{Package: "main", Ident: "ARepository"},
			},
		},
		{
			"multiple Repository returns slice with Repository",
			`
      package main

      type (
        ARepository interface { }
        BRepository interface { }
      )
      type CRepository interface { }
      `,
			[]*Repository{
				{Package: "main", Ident: "ARepository"},
				{Package: "main", Ident: "BRepository"},
				{Package: "main", Ident: "CRepository"},
			},
		},
		{
			"Methods with no args and returns",
			`
      package main

      type Repository interface {
        A()
        B()
        C()
      }
      `,
			[]*Repository{
				{
					Package: "main",
					Ident:   "Repository",
					Methods: []*Method{
						{Ident: "A"},
						{Ident: "B"},
						{Ident: "C"},
					},
				},
			},
		},
		{
			"Methods with args and returns",
			`
      package main

      type Repository interface {
        A(a int) int
        B(a int, b string) (int, error)
        C(a int, b, c bool) (_ int, err error, x string)
      }
      `,
			[]*Repository{
				{
					Package: "main",
					Ident:   "Repository",
					Methods: []*Method{
						{
							Ident:   "A",
							Params:  []*Param{{Ident: "a", Type: "int"}},
							Returns: []*Param{{Type: "int"}},
						},
						{
							Ident: "B",
							Params: []*Param{
								{Ident: "a", Type: "int"},
								{Ident: "b", Type: "string"},
							},
							Returns: []*Param{
								{Type: "int"},
								{Type: "error"},
							},
						},
						{
							Ident: "C",
							Params: []*Param{
								{Ident: "a", Type: "int"},
								{Ident: "b", Type: "bool"},
								{Ident: "c", Type: "bool"},
							},
							Returns: []*Param{
								{Ident: "_", Type: "int"},
								{Ident: "err", Type: "error"},
								{Ident: "x", Type: "string"},
							},
						},
					},
				},
			},
		},
		{
			"Ignores whitespace",
			`
      package main

      type Repository interface {
        A      (a   int)       int
        B(a   int,   b             string   ) (   int,   error   )
        C(a      int, b    , c  bool) (_    int, err    error,     x string    )
      }
      `,
			[]*Repository{
				{
					Package: "main",
					Ident:   "Repository",
					Methods: []*Method{
						{
							Ident:   "A",
							Params:  []*Param{{Ident: "a", Type: "int"}},
							Returns: []*Param{{Type: "int"}},
						},
						{
							Ident: "B",
							Params: []*Param{
								{Ident: "a", Type: "int"},
								{Ident: "b", Type: "string"},
							},
							Returns: []*Param{
								{Type: "int"},
								{Type: "error"},
							},
						},
						{
							Ident: "C",
							Params: []*Param{
								{Ident: "a", Type: "int"},
								{Ident: "b", Type: "bool"},
								{Ident: "c", Type: "bool"},
							},
							Returns: []*Param{
								{Ident: "_", Type: "int"},
								{Ident: "err", Type: "error"},
								{Ident: "x", Type: "string"},
							},
						},
					},
				},
			},
		},
		{
			"multiple Repository with methods",
			`
      package main

      type ARepository interface { A() }
      type BRepository interface { B(); C() }
      `,
			[]*Repository{
				{
					Package: "main",
					Ident:   "ARepository",
					Methods: []*Method{
						{Ident: "A"},
					},
				},
				{
					Package: "main",
					Ident:   "BRepository",
					Methods: []*Method{
						{Ident: "B"},
						{Ident: "C"},
					},
				},
			},
		},
		{
			"imports are parsed",
			`
      package main
      
      import (
        "context"
        "fmt"
        err "errors"
      )

      type ARepository interface {}
      type BRepository interface {}
      `,
			[]*Repository{
				{
					Package: "main",
					Ident:   "ARepository",
					Imports: []Import{
						{Path: "context"},
						{Path: "fmt"},
						{Path: "errors", Name: "err"},
					},
				},
				{
					Package: "main",
					Ident:   "BRepository",
					Imports: []Import{
						{Path: "context"},
						{Path: "fmt"},
						{Path: "errors", Name: "err"},
					},
				},
			},
		},
	} {
		// slog.SetLogLoggerLevel(slog.LevelDebug)
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			ctx := context.Background()
			tree, err := tsparser.ParseCtx(ctx, nil, []byte(test.src))
			require.NoError(err)
			repos, err := parseRepositories([]byte(test.src), tree)
			require.NoError(err)
			testRepositories(t, test.expect, repos)
		})
	}
}

func TestParseRepositoryImplFile(t *testing.T) {
	type expect struct {
		packageName string
		repImpls    []string
		methods     map[string][]string
	}
	for _, test := range []struct {
		name   string
		have   string
		expect expect
	}{
		{
			"empty file",
			"package main",
			expect{
				"main",
				[]string{},
				map[string][]string{},
			},
		},
		{
			"parses repository impl declarations",
			`
      package main

      type (
        repositoryImpl    struct {}
        fooRepositoryImpl struct {}
      )

      type bazRepositoryImpl struct {}
      `,
			expect{
				"main",
				[]string{
					"repositoryImpl",
					"fooRepositoryImpl",
					"bazRepositoryImpl",
				},
				map[string][]string{},
			},
		},
		{
			"parses method declarations",
			`
      package main

      type (
        fooRepositoryImpl struct {}
        bazRepositoryImpl struct {}
      )

      func (i fooRepositoryImpl)  A() {}
      func (i *fooRepositoryImpl) B() {}
      func (i bazRepositoryImpl)  C() {}
      `,
			expect{
				"main",
				[]string{
					"fooRepositoryImpl",
					"bazRepositoryImpl",
				},
				map[string][]string{
					"fooRepositoryImpl": {"A", "B"},
					"bazRepositoryImpl": {"C"},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			ctx := context.Background()
			pkg, impls, methods, err := parseRepositoryImplFile(ctx, []byte(test.have))
			require.NoError(err)
			require.Equal(test.expect.packageName, pkg)
			require.ElementsMatch(test.expect.repImpls, impls)
			require.Equal(test.expect.methods, methods)
		})
	}
}

func TestParseRepositoryImpls(t *testing.T) {
	for _, test := range []struct {
		name        string
		fsys        map[string]string
		packagePath string
		files       []string
		have        []*Repository
		expect      []*RepositoryImpl
	}{
		{
			"empty fsys returns empty slice",
			map[string]string{},
			"internal",
			[]string{},
			[]*Repository{},
			[]*RepositoryImpl{},
		},
		{
			"suffixes impl as default impl package name",
			map[string]string{},
			"api",
			[]string{},
			[]*Repository{
				{
					Package:  "api",
					Filename: "one.go",
					Ident:    "ARepository",
				},
			},
			[]*RepositoryImpl{
				{
					IsNew:        true,
					ImplPackage:  "apiimpl",
					ImplFilename: "a_impl.go",
					ImplMethods:  []string{},
				},
			},
		},
		{
			"multiple repositories across multiple files",
			map[string]string{
				"internal/one.go": `
        package internal

        type aRepositoryImpl struct {}

        func (a *aRepositoryImpl) A() {}
        `,
				"internal/two.go": `
        package internal

        type bRepositoryImpl struct {}

        func (b bRepositoryImpl) B() {}
        `,
			},
			"internal",
			[]string{"one.go", "two.go"},
			[]*Repository{
				{
					Package:  "api",
					Filename: "one.go",
					Ident:    "ARepository",
					Methods: []*Method{
						{Ident: "A"},
						{Ident: "B"},
					},
				},
				{
					Package:  "api",
					Filename: "two.go",
					Ident:    "BRepository",
					Methods:  []*Method{{Ident: "B"}},
				},
			},
			[]*RepositoryImpl{
				{
					Repository{},
					false,
					"internal",
					"one.go",
					[]string{"A"},
				},
				{
					Repository{},
					false,
					"internal",
					"two.go",
					[]string{"B"},
				},
			},
		},
		{
			"single repositories methods dispersed across multiple files",
			map[string]string{
				"internal/one.go": `
        package internal

        type repositoryImpl struct {}

        func (a *repositoryImpl) A() {}
        `,
				"internal/two.go": `
        package internal

        func (b repositoryImpl) B() {}
        `,
			},
			"internal",
			[]string{"one.go", "two.go"},
			[]*Repository{
				{
					Package:  "api",
					Filename: "one.go",
					Ident:    "Repository",
					Methods: []*Method{
						{Ident: "A"},
						{Ident: "B"},
					},
				},
			},
			[]*RepositoryImpl{
				{
					Repository{},
					false,
					"internal",
					"one.go",
					[]string{"A", "B"},
				},
			},
		},
		{
			"single repositories methods dispersed across multiple files",
			map[string]string{
				"internal/one.go": `
        package internal

        type repositoryImpl struct {}

        func (a *repositoryImpl) A() {}
        `,
				"internal/two.go": `
        package internal

        func (b repositoryImpl) B() {}
      j `,
			},
			"internal",
			[]string{"one.go", "two.go"},
			[]*Repository{
				{
					Package:  "api",
					Filename: "one.go",
					Ident:    "Repository",
					Methods: []*Method{
						{Ident: "A"},
						{Ident: "B"},
					},
				},
			},
			[]*RepositoryImpl{
				{
					Repository{},
					false,
					"internal",
					"one.go",
					[]string{"A", "B"},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			ctx := context.Background()
			fsys := make(fstest.MapFS, len(test.fsys))
			for path, content := range test.fsys {
				fsys[path] = &fstest.MapFile{Data: []byte(content), Mode: 0644}
			}
			got, err := parseRepositoryImpls(
				ctx,
				fsys,
				test.packagePath,
				test.have,
			)
			require.NoError(err)
			testRepositoryImpls(t, test.expect, got)
		})
	}
}

func testRepositories(t *testing.T, expected, actual []*Repository) {
	t.Helper()
	require := require.New(t)
	require.Len(actual, len(expected))
	for i, repo := range actual {
		expect := expected[i]
		require.Equal(expect.Package, repo.Package)
		require.Equal(expect.Ident, repo.Ident)
		require.Len(repo.Methods, len(expect.Methods))
		for j, method := range repo.Methods {
			require.Equal(expect.Methods[j], method)
		}
		require.ElementsMatch(expect.Imports, repo.Imports)
	}
}

func testRepositoryImpls(t *testing.T, expected, actual []*RepositoryImpl) {
	t.Helper()
	require := require.New(t)
	require.Len(actual, len(expected))
	for i, repo := range actual {
		expect := expected[i]
		require.Equal(expect.ImplPackage, repo.ImplPackage)
		require.Equal(expect.ImplFilename, repo.ImplFilename)
		require.Equal(expect.IsNew, repo.IsNew)
		require.Len(repo.ImplMethods, len(expect.ImplMethods))
		for j, method := range repo.ImplMethods {
			require.Equal(expect.ImplMethods[j], method)
		}
	}
}
