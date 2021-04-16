package decorator

import (
	"fmt"
	"path/filepath"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
)

type Plugin struct{}

func (p *Plugin) GenerateCode(cfg *codegen.Data) error {
	var dd []*ObjectDecorator
	for _, obj := range cfg.Objects {
		if !obj.HasResolvers() {
			continue
		}

		resolverName := fmt.Sprintf("%sResolver", obj.Name)
		dd = append(dd, &ObjectDecorator{
			ResolverType:    resolverName,
			ConstructorName: fmt.Sprintf("Decorate%s", resolverName),
			FuncType:        fmt.Sprintf("%sDecoratorFunc", resolverName),
		})
	}

	return templates.Render(templates.Options{
		Template:        tpl,
		PackageName:     cfg.Config.Exec.Package,
		Filename:        filepath.Dir(cfg.Config.Exec.Filename) + "/resolver_decorators.go",
		RegionTags:      true,
		GeneratedHeader: true,
		Data:            data{Decorators: dd},
		Packages:        cfg.Config.Packages,
	})
}

func (p Plugin) Name() string {
	return "resolver_decorators"
}

type data struct {
	Decorators []*ObjectDecorator
}

type ObjectDecorator struct {
	ResolverType    string
	ConstructorName string
	FuncType        string
}
