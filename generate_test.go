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
    ctx, span := otel.GetTracerProvider().Tracer("foo").Start(ctx, "Repository.A")
    defer span.End()
    _ = ctx
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
    ctx, span := otel.GetTracerProvider().Tracer("foo").Start(ctx, "Repository.A")
    defer func() {
      if err != nil {
        err = eris.Wrap(err, "foo.Repository.A")
        span.SetStatus(codes.Error, "")
        span.RecordError(err)
      }
      span.End()
    }()
    _ = ctx
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

func TestGenerateRepositoryTest(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   RepositoryImpl
		expect string
	}{
		{
			"no methods",
			RepositoryImpl{
				Repository: Repository{
					Ident:   "Repository",
					Package: "api",
				},
				ImplPackage: "internal",
			},
			`
func TestRepository(t *testing.T) {
	t.Parallel()
	type Deps struct {
		fx.In
		Repository api.Repository
	}
	ctx := context.Background()
	_ = ctx
	app := fxtest.New(t,
		fx.Provide(
			// Add dependencies here
		),
    internal.Options,

		fxutils.Test(t, func(t *testing.T, d internal.Dependencies) {
      r := d.Repository
      _ = r
			// Initialization code here
		}),
	)
	app.Run()
}
`,
		},
		{
			"multiple methods",
			RepositoryImpl{
				Repository: Repository{
					Ident:   "Repository",
					Package: "api",
					Methods: []*Method{
						{
							Ident: "A",
						},
						{
							Ident: "B",
						},
					},
				},
				ImplPackage: "internal",
			},
			`
func TestRepository(t *testing.T) {
	t.Parallel()
	type Deps struct {
		fx.In
		Repository api.Repository
	}
	ctx := context.Background()
	_ = ctx
	app := fxtest.New(t,
		fx.Provide(
			// Add dependencies here
		),
    internal.Options,

		fxutils.Test(t, func(t *testing.T, d internal.Dependencies) {
      r := d.Repository
      _ = r
			// Initialization code here

			t.Run("A", func(t *testing.T) {
				t.Skip()
			})

			t.Run("B", func(t *testing.T) {
				t.Skip()
			})
		}),
	)
	app.Run()
}
`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			got, err := generateRepositoryTest(test.have)
			require.NoError(err)
			require.Equal(test.expect, got)
		})
	}
}

func TestGenerateRepositoryTestsForFile(t *testing.T) {
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
				"internal/one_test.go": `// Some comment
package internal_test

import (
  "errors"
)

type repositoryImpl struct{}

var x errors.Something
`,
			},
			"internal/one_test.go",
			[]*RepositoryImpl{
				{
					Repository:  Repository{Package: "api"},
					ImplPackage: "internal",
				},
			},
			`// Some comment
package internal_test

import (
	"errors"
)

type repositoryImpl struct{}

var x errors.Something
`,
		},
		{
			"tests for new repositories added to end of file",
			map[string]string{
				"go.mod": `
        module example
        `,
				"internal/one_test.go": `package internal_test

import (
  "errors"
)

type repositoryImpl struct{}

var x errors.Something
`,
			},
			"internal/one_test.go",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "api",
						Ident:       "Repository",
						Methods: []*Method{
							{Ident: "A"},
							{Ident: "B"},
						},
					},
					ImplPackage:     "internal",
					ImplPackagePath: "internal",
					IsNew:           true,
				},
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "api",
						Ident:       "Another",
						Methods: []*Method{
							{Ident: "A"},
							{Ident: "B"},
						},
					},
					ImplPackage:     "internal",
					ImplPackagePath: "internal",
					IsNew:           true,
				},
				{
					Repository: Repository{
						Package:     "api",
						PackagePath: "api",
						Ident:       "Nope",
					},
					ImplPackage:     "internal",
					ImplPackagePath: "internal",
				},
			},
			`package internal_test

import (
	"context"
	"errors"
	"example/api"
	"example/internal"
	"testing"

	fxutils "github.com/MathGaps/core/pkg/utils/fx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type repositoryImpl struct{}

var x errors.Something

func TestRepository(t *testing.T) {
	t.Parallel()
	type Deps struct {
		fx.In
		Repository api.Repository
	}
	ctx := context.Background()
	_ = ctx
	app := fxtest.New(t,
		fx.Provide(
		// Add dependencies here
		),
		internal.Options,

		fxutils.Test(t, func(t *testing.T, d internal.Dependencies) {
			r := d.Repository
			_ = r
			// Initialization code here

			t.Run("A", func(t *testing.T) {
				t.Skip()
			})

			t.Run("B", func(t *testing.T) {
				t.Skip()
			})
		}),
	)
	app.Run()
}

func TestAnother(t *testing.T) {
	t.Parallel()
	type Deps struct {
		fx.In
		Repository api.Another
	}
	ctx := context.Background()
	_ = ctx
	app := fxtest.New(t,
		fx.Provide(
		// Add dependencies here
		),
		internal.AnotherOptions,

		fxutils.Test(t, func(t *testing.T, d internal.AnotherDependencies) {
			r := d.Repository
			_ = r
			// Initialization code here

			t.Run("A", func(t *testing.T) {
				t.Skip()
			})

			t.Run("B", func(t *testing.T) {
				t.Skip()
			})
		}),
	)
	app.Run()
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
			got, err := generateRepositoryTestsForFile(
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

func TestGenerateRepositoryStubFile(t *testing.T) {
	for _, test := range []struct {
		name   string
		have   []*RepositoryImpl
		expect string
	}{
		{
			"no repositories",
			[]*RepositoryImpl{},
			`// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

import "go.uber.org/fx"

var Repositories = fx.Options()
`,
		},
		{
			"repositories from same package",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Ident:       "B",
						PackagePath: "api/waltuh",
						Filename:    "b.go",
					},
					ImplFilename:    "b_impl.go",
					ImplPackage:     "waltuh",
					ImplPackagePath: "internal/waltuh",
				},
				{
					Repository: Repository{
						Ident:       "Repository",
						PackagePath: "api/waltuh",
						Filename:    "repository.go",
					},
					ImplFilename:    "repository_impl.go",
					ImplPackage:     "waltuh",
					ImplPackagePath: "internal/waltuh",
				},
				{
					Repository: Repository{
						Ident:       "A",
						PackagePath: "api/waltuh",
						Filename:    "a.go",
					},
					ImplFilename:    "a_impl.go",
					ImplPackage:     "waltuh",
					ImplPackagePath: "internal/waltuh",
				},
			},
			`// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=waltuh/mocks.go -pkg=waltuh -rm -skip-ensure ../api/waltuh A B Repository

import (
	"example/internal/waltuh"

	"go.uber.org/fx"
)

var Repositories = fx.Options(
	waltuh.Options,
	waltuh.AOptions,
	waltuh.BOptions,
)
`,
		},
		{
			"repositories from different packages",
			[]*RepositoryImpl{
				{
					Repository: Repository{
						Ident:       "Repository",
						Package:     "waltuh",
						PackagePath: "api/waltuh",
						Filename:    "repository.go",
					},
					ImplFilename:    "repository_impl.go",
					ImplPackage:     "waltuhimpl",
					ImplPackagePath: "internal/waltuhimpl",
				},
				{
					Repository: Repository{
						Ident:       "Repository",
						Package:     "jesse",
						PackagePath: "api/jesse",
						Filename:    "repository.go",
					},
					ImplFilename:    "repository_impl.go",
					ImplPackage:     "jesseimpl",
					ImplPackagePath: "internal/jesseimpl",
				},
			},
			`// DO NOT MODIFY
// This file will be automatically regenerated based on the API.
package internal

//go:generate moq -out=jesseimpl/mocks.go -pkg=jesseimpl -rm -skip-ensure ../api/jesse Repository
//go:generate moq -out=waltuhimpl/mocks.go -pkg=waltuhimpl -rm -skip-ensure ../api/waltuh Repository

import (
	"example/internal/jesseimpl"
	"example/internal/waltuhimpl"

	"go.uber.org/fx"
)

var Repositories = fx.Options(
	jesseimpl.Options,
	waltuhimpl.Options,
)
`,
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			cli.Impl = "internal"
			require := require.New(t)
			fsys := make(fstest.MapFS)
			fsys["go.mod"] = &fstest.MapFile{Data: []byte(`
        module example

        go 1.22.1`,
			), Mode: 0644}
			got, err := generateRepositoryStubFile(
				fsys,
				"internal",
				test.have...,
			)
			require.NoError(err)
			require.Equal(test.expect, got)
		})
	}
}
