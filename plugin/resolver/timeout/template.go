package timeout

// TODO: switch to go embed when migrate to go1.16
var (
	wrappersTpl = `
{{ reserveImport "context"  }}
{{ reserveImport "time"  }}

{{ range $w := .Wrappers }}
    func New{{$w.ObjectType}}Decorator(cfg TimeoutsConfig) {{$w.ObjectType}}DecoratorFunc {
        return func(parent {{$w.ObjectType}}) {{$w.ObjectType}} {
            return &{{$w.Type}} {
                cfg: cfg,
                parent: parent,
            }
        }
    }

    type {{$w.Type}} struct {
        cfg TimeoutsConfig
        parent {{$w.ObjectType}}
    }

    {{ range $r := $w.Resolvers }}
        // {{ $r.GoFieldName }} is an adapter method for invoking original method with timeout.
        // config example:
        // {{$w.ObjectName}}:
        //   {{$r.Name}}: 10s
        func (r *{{$w.Type}}) {{ $r.GoFieldName }}{{ $r.ResolverImplementation }} {
            dur, ok := r.cfg.GetTimeout({{$w.ObjectName|quote}}, {{$r.Name|quote}})
            if ok {
                var cancel func()
                ctx, cancel = context.WithTimeout(ctx, dur)
                defer cancel()
            }

            return {{$r.Invocation}}
        }
    {{ end }}
{{ end }}

type TimeoutsConfig interface{
    GetTimeout(object, field string) (time.Duration, bool)
}
	`
	yamlConfigExampleTpl = `
// yaml configuration sample
//
/*
   Objects:
   {{- range $w := .Wrappers }}
	 - Name: {{$w.ObjectName | quote}}
	   Fields:
	   {{- range $r := $w.Resolvers }}
		 - Name: {{ $r.Name | quote }}
		   Timeout: 1s
	   {{- end }}
   {{ end }}
*/
	`
)

func tpl() string {
	return wrappersTpl +
		yamlConfigExampleTpl
}
