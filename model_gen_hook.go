package main

import (
	"go/types"

	"github.com/99designs/gqlgen/plugin/modelgen"
)

func JsonOmitempty(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			switch field.Type.(type) {
			case *types.Pointer:
				field.Tag = `json:"` + field.Name + `,omitempty"`
			}
		}
	}

	return b
}
