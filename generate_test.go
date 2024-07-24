package main

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestNewMethods(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   RepositoryImpl
		expect []*Method
	}{
		{
			"no methods",
			RepositoryImpl{
				Repository: Repository{
					Methods: []*Method{{Ident: "A"}},
				},
				ImplMethods: []string{"A"},
			},
			[]*Method{},
		},
		{
			"one method",
			RepositoryImpl{
				Repository: Repository{
					Methods: []*Method{{Ident: "A"}},
				},
				ImplMethods: []string{"B"},
			},
			[]*Method{{Ident: "A"}},
		},
		{
			"method args/returns are qualified with package name",
			RepositoryImpl{
				Repository: Repository{
					Package: "api",
					Methods: []*Method{
						{
							Ident: "A",
							Params: Params{
								{Ident: "a", Type: "Zoowee"},
								{Ident: "ctx", Type: "context.Context"},
							},
							Returns: Params{
								{Ident: "b", Type: "Mama"},
							},
						},
					},
				},
				ImplMethods: []string{"ctx"},
			},
			[]*Method{
				{
					Ident: "A",
					Params: Params{
						{Ident: "a", Type: "api.Zoowee"},
						{Ident: "ctx", Type: "context.Context"},
					},
					Returns: Params{
						{Ident: "b", Type: "api.Mama"},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			got := test.have.NewMethods()
			require.Len(got, len(test.expect))
			require.ElementsMatch(got, test.expect)
		})
	}
}

func TestParamsHas(t *testing.T) {
	require := require.New(t)
	noCtxErr := Params{
		{Type: "int"},
	}
	ctxErr := Params{
		{Type: "context.Context"},
		{Type: "error"},
	}
	require.False(noCtxErr.HasCtx())
	require.False(noCtxErr.HasError())
	require.True(ctxErr.HasCtx())
	require.True(ctxErr.HasError())
}

func TestParamsSrc(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   Params
		expect string
	}{
		{
			"empty params",
			Params{},
			"",
		},
		{
			"single param with ident and type",
			Params{
				{"a", "int"},
			},
			"a int",
		},
		{
			"multiple params with ident and type",
			Params{
				{"a", "int"},
				{"ok", "bool"},
			},
			"a int, ok bool",
		},
		{
			"ctx is qualified",
			Params{
				{"", "context.Context"},
				{"", "int"},
			},
			"ctx context.Context, _ int",
		},
		{
			"ctx is qualified and other parmas retain name",
			Params{
				{"_", "context.Context"},
				{"yoyo", "int"},
			},
			"ctx context.Context, yoyo int",
		},
		{
			"err is qualified",
			Params{
				{Type: "error"},
				{Type: "bool"},
			},
			"err error, _ bool",
		},
		{
			"adjacent types are grouped",
			Params{
				{"yep", "bool"},
				{"nope", "bool"},
				{"one", "int"},
				{"two", "int"},
				{"three", "int"},
			},
			"yep, nope bool, one, two, three int",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			got := test.have.ParamsSrc()
			require.Equal(test.expect, got)
		})
	}
}

func TestReturnsSrc(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   Params
		expect string
	}{
		{
			"bracketed if named",
			Params{
				{"a", "int"},
			},
			"(a int)",
		},
		{
			"no brackets if not named",
			Params{
				{Type: "int"},
			},
			"int",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			got := test.have.ReturnsSrc()
			require.Equal(test.expect, got)
		})
	}
}

func TestName(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   Repository
		expect string
	}{
		{
			"repository suffix is trimmed",
			Repository{Ident: "FooRepository"},
			"Foo",
		},
		{
			"ident if no Repository suffix",
			Repository{Ident: "Foo"},
			"Foo",
		},
		{
			"Repository if Ident is Repository",
			Repository{Ident: "Repository"},
			"Repository",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			implName := test.have.Name()
			require.Equal(test.expect, implName)
		})
	}
}

func TestDependencies(t *testing.T) {
	for _, test := range []struct {
		name    string
		have    Repository
		qualify string
		expect  string
	}{
		{
			"S if Repository",
			Repository{Ident: "Repository"},
			"S",
			"S",
		},
		{
			"S is suffixed to name",
			Repository{Ident: "FooRepository"},
			"S",
			"FooS",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			implName := test.have.QualifyString(test.qualify)
			require.Equal(test.expect, implName)
		})
	}
}

func TestQualifiedName(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   Repository
		expect string
	}{
		{
			"package is prefixed",
			Repository{
				Package: "foo",
				Ident:   "Repository",
			},
			"foo.Repository",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			implName := test.have.QualifiedName()
			require.Equal(test.expect, implName)
		})
	}
}

func TestImplName(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   Repository
		expect string
	}{
		{
			"first letter of Ident is lowercased",
			Repository{Ident: "FooRepository"},
			"fooRepositoryImpl",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			implName := test.have.ImplName()
			require.Equal(test.expect, implName)
		})
	}
}

func TestGenerateMethodImpl(t *testing.T) {
	type Input struct {
		Repository
		Method
	}
	for _, test := range []struct {
		name   string
		have   Input
		expect string
	}{
		{
			"no context and no error",
			Input{
				Repository{
					Ident: "Repository",
				},
				Method{
					Ident: "A",
				},
			},
			`
  func (r *repositoryImpl) A() {
    panic("TODO: implement Repository.A")
  }
`,
		},
		{
			"err is wrapped with method metadata",
			Input{
				Repository{
					Package: "foo",
					Ident:   "Repository",
				},
				Method{
					Ident: "A",
					Returns: Params{
						{Type: "bool"},
						{Type: "error"},
					},
				},
			},
			`
  func (r *repositoryImpl) A() (_ bool, err error) {
    defer func() {
      if err != nil {
        err = eris.Wrap(err, "foo.Repository.A")
      }
    }()
    panic("TODO: implement foo.Repository.A")
  }
`,
		},
		{
			"otel span is created and ended if ctx is provided",
			Input{
				Repository{
					Package: "foo",
					Ident:   "Repository",
				},
				Method{
					Ident: "A",
					Params: Params{
						{Type: "context.Context"},
						{Type: "bool"},
					},
				},
			},
			`
  func (r *repositoryImpl) A(ctx context.Context, _ bool) {
    ctx, span := otel.GetTracerProvider().Tracer("foo").Start(ctx, "Repository.A") //nolint:all
    defer span.End()
    panic("TODO: implement foo.Repository.A")
  }
`,
		},
		{
			"otel span is created and err is tracked if ctx and err are provided",
			Input{
				Repository{
					Package: "foo",
					Ident:   "Repository",
				},
				Method{
					Ident: "A",
					Params: Params{
						{Type: "context.Context"},
					},
					Returns: Params{
						{Type: "error"},
					},
				},
			},
			`
  func (r *repositoryImpl) A(ctx context.Context) (err error) {
    ctx, span := otel.GetTracerProvider().Tracer("foo").Start(ctx, "Repository.A") //nolint:all
    defer func() {
      if err != nil {
        err = eris.Wrap(err, "foo.Repository.A")
        span.SetStatus(codes.Error, "")
        span.RecordError(err)
      }
      span.End()
    }()
    panic("TODO: implement foo.Repository.A")
  }
`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			impl, err := generateMethodImpl(test.have.Repository, test.have.Method)
			require.NoError(err)
			require.Equal(test.expect, impl)
		})
	}
}

func TestGenerateRepositoryImpl(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   Repository
		expect string
	}{
		{
			"generates repository without qualification",
			Repository{
				Package: "foo",
				Ident:   "Repository",
			},
			`
type Dependencies struct {
  fx.In
	// Add dependencies here
}

var Options = fx.Options(
	fx.Provide(
		NewRepository,
	),
)

func NewRepository(deps Dependencies) foo.Repository {
	return &repositoryImpl{
    Dependencies: deps,
	}
}

type repositoryImpl struct {
  Dependencies
}
`,
		},
		{
			"generates repository without qualification",
			Repository{
				Package: "foo",
				Ident:   "BarRepository",
			},
			`
type BarDependencies struct {
  fx.In
	// Add dependencies here
}

var BarOptions = fx.Options(
	fx.Provide(
		NewBarRepository,
	),
)

func NewBarRepository(deps BarDependencies) foo.BarRepository {
	return &barRepositoryImpl{
    BarDependencies: deps,
	}
}

type barRepositoryImpl struct {
  BarDependencies
}
`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			impl, err := generateRepositoryImpl(test.have)
			require.NoError(err)
			require.Equal(test.expect, impl)
		})
	}
}

func TestGenerateRepositoryImplsForFile(t *testing.T) {
	for _, test := range []struct {
		name     string
		fsys     map[string]string
		filepath string
		have     []*RepositoryImpl
		expect   string
	}{
		{
			"existing src retained",
			map[string]string{
				"go.mod": `
        module example
        `,
				"internal/one.go": `// Some comment
package internal

import (
  "errors"
)

type repositoryImpl struct{}

var x errors.Something
`,
			},
			"internal/one.go",
			[]*RepositoryImpl{
				{
					Repository:  Repository{Package: "api"},
					ImplPackage: "internal",
				},
			},
			`// Some comment
package internal

import (
	"errors"
)

type repositoryImpl struct{}

var x errors.Something
`,
		},
		{
			"methods appended to end of file and api import added",
			map[string]string{
				"go.mod": `
        module example
        `,
				"internal/one.go": `package internal

import (
  "errors"
)

type repositoryImpl struct{}

var x errors.Something
`,
			},
			"internal/one.go",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "example/api",
						Ident:       "Repository",
						Methods: []*Method{
							{Ident: "A"},
						},
					},
					ImplPackage: "internal",
					ImplMethods: []string{},
				},
			},
			`package internal

import (
	"errors"
)

type repositoryImpl struct{}

var x errors.Something

func (r *repositoryImpl) A() {
	panic("TODO: implement api.Repository.A")
}
`,
		},
		{
			"new repositories added to end of file",
			map[string]string{
				"go.mod": `
        module example
        `,
				"internal/one.go": `package internal

import (
  "errors"
)

type repositoryImpl struct{}

var x errors.Something
`,
			},
			"internal/one.go",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "example/api",
						Ident:       "Repository",
						Methods: []*Method{
							{Ident: "A"},
						},
					},
					ImplPackage: "internal",
					ImplMethods: []string{},
				},
			},
			`package internal

import (
	"errors"
)

type repositoryImpl struct{}

var x errors.Something

func (r *repositoryImpl) A() {
	panic("TODO: implement api.Repository.A")
}
`,
		},
		{
			"repository imports are propagated",
			map[string]string{
				"go.mod": `
        module example
        `,
				"internal/one.go": `package internal

type repositoryImpl struct{}

var x something.Something
`,
			},
			"internal/one.go",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "example/api",
						Ident:       "Repository",
						Methods:     []*Method{},
						Imports: []Import{
							{Path: "somewhere/something"},
						},
					},
					ImplPackage: "internal",
					ImplMethods: []string{},
				},
			},
			`package internal

import "somewhere/something"

type repositoryImpl struct{}

var x something.Something
`,
		},
		{
			"new repositories are added to end of file",
			map[string]string{
				"go.mod": `
        module example
        `,
				"internal/one.go": `package internal

var hoyminoy string
`,
			},
			"internal/one.go",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "api",
						Ident:       "Repository",
					},
					IsNew:       true,
					ImplPackage: "internal",
					ImplMethods: []string{},
				},
			},
			`package internal

import (
	"example/api"

	"go.uber.org/fx"
)

var hoyminoy string

type Dependencies struct {
	fx.In
	// Add dependencies here
}

var Options = fx.Options(
	fx.Provide(
		NewRepository,
	),
)

func NewRepository(deps Dependencies) api.Repository {
	return &repositoryImpl{
		Dependencies: deps,
	}
}

type repositoryImpl struct {
	Dependencies
}
`,
		},
		{
			"named imports are preserved",
			map[string]string{
				"go.mod": `
        module example
        `,
				"internal/one.go": `package internal

import (
	hahafucku "example/api"

	"go.uber.org/fx"
)

func NewRepository(deps Dependencies) hahafucku.Repository {
	return &repositoryImpl{
		Dependencies: deps,
	}
}
`,
			},
			"internal/one.go",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "api",
						Ident:       "Repository",
					},
					IsNew:       false,
					ImplPackage: "internal",
					ImplMethods: []string{},
				},
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "api",
						Ident:       "BRepository",
					},
					IsNew:       true,
					ImplPackage: "internal",
					ImplMethods: []string{},
				},
			},
			`package internal

import (
	hahafucku "example/api"

	"go.uber.org/fx"
)

func NewRepository(deps Dependencies) hahafucku.Repository {
	return &repositoryImpl{
		Dependencies: deps,
	}
}

type BDependencies struct {
	fx.In
	// Add dependencies here
}

var BOptions = fx.Options(
	fx.Provide(
		NewBRepository,
	),
)

func NewBRepository(deps BDependencies) hahafucku.BRepository {
	return &bRepositoryImpl{
		BDependencies: deps,
	}
}

type bRepositoryImpl struct {
	BDependencies
}
`,
		},
		{
			"package created from scratch",
			map[string]string{
				"go.mod": `
        module example
        `,
			},
			"internal/one.go",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "api",
						Ident:       "Repository",
					},
					IsNew:       true,
					ImplPackage: "internal",
					ImplMethods: []string{},
				},
			},
			`// This file will be automatically regenerated based on the API. Any repository implementations
// will be copied through when generating and new methods will be added to the end.
package internal

import (
	"example/api"

	"go.uber.org/fx"
)

type Dependencies struct {
	fx.In
	// Add dependencies here
}

var Options = fx.Options(
	fx.Provide(
		NewRepository,
	),
)

func NewRepository(deps Dependencies) api.Repository {
	return &repositoryImpl{
		Dependencies: deps,
	}
}

type repositoryImpl struct {
	Dependencies
}
`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			fsys := make(fstest.MapFS, len(test.fsys))
			for path, content := range test.fsys {
				fsys[path] = &fstest.MapFile{Data: []byte(content), Mode: 0644}
			}
			got, err := generateRepositoryImplsForFile(
				fsys,
				test.filepath,
				test.have,
			)
			require.NoError(err)
			t.Log(got)
			require.Equal(test.expect, got)
		})
	}
}
