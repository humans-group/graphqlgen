package decorator

const tpl = `
{{ range $d := .Decorators}}
    func {{$d.ConstructorName}}(
        r {{$d.ResolverType}},
        ff ...{{$d.FuncType}},
    ) {{$d.ResolverType}} {
        for _, f := range ff {
            r = f(r)
        }

        return r
    }

    type {{$d.FuncType}} func({{$d.ResolverType}}) {{$d.ResolverType}}
{{ end }}
`
