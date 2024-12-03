package timeout

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
)

type Plugin struct {
	outFile string
}

func NewPlugin(outFile string) *Plugin {
	return &Plugin{
		outFile: outFile,
	}
}

func (p *Plugin) GenerateCode(cfg *codegen.Data) error {
	var ww []*ObjectWrapper
	for _, obj := range cfg.Objects {
		w := ObjectWrapper{
			ObjectName: obj.Name,
			ObjectType: fmt.Sprintf("%sResolver", obj.Name),
			Type:       fmt.Sprintf("%sResolverWithTimeout", obj.Name),
		}

		for _, f := range obj.Fields {
			if !f.IsResolver {
				continue
			}

			w.Resolvers = append(w.Resolvers, &FieldResolver{
				Field: f,
			})
		}

		if len(w.Resolvers) > 0 {
			ww = append(ww, &w)
		}
	}

	return templates.Render(templates.Options{
		Template:        tpl(),
		PackageName:     cfg.Config.Exec.Package,
		Filename:        filepath.Dir(cfg.Config.Exec.Filename) + "/" + p.outFile,
		Packages:        cfg.Config.Packages,
		GeneratedHeader: true,
		Data:            data{Wrappers: ww},
	})
}

func (p Plugin) Name() string {
	return "resolver_timeouts"
}

type data struct {
	Wrappers []*ObjectWrapper
}

type ObjectWrapper struct {
	Type       string
	ObjectName string
	ObjectType string
	Resolvers  []*FieldResolver
}

type FieldResolver struct {
	*codegen.Field
}

func (r *FieldResolver) Invocation() string {
	args := make([]string, 0, len(r.Args))
	if !r.Object.Root {
		args = append(args, "obj")
	}

	for _, arg := range r.Args {
		args = append(args, arg.Name)
	}

	return fmt.Sprintf("r.parent.%s(ctx, %s)",
		r.GoFieldName, strings.Join(args, ", "))
}

func (r *FieldResolver) ResolverImplementation() string {
	return r.ShortResolverDeclaration()
}
