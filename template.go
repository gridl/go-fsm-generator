package main

import "text/template"

var embeddedTemplate = template.Must(template.New("embedded").Parse(`
	package {{.PkgName}}

	// Generated by fsm-generator. DO NOT EDIT.

	//+++ General machine definition +++
	{{$mName := .MachineName}}
	type {{$mName}}State int
	const (
		_ {{$mName}}State = iota
		{{- range $st, $stDef := .States}}
			{{$st}}
		{{- end}}
	)

	var _{{$mName}}StateMap = map[{{$mName}}State]string{
		{{- range $st, $stDef := .States}}
			{{$st}}: "{{$st}}",
		{{- end}}
	}

	func (s {{$mName}}State) String() string {
		return _{{$mName}}StateMap[s]
	}

	
	type {{$mName}}Behaviour interface {
	{{- range $st, $stDef := .States}}
		{{- if ($stDef.IsTerminal)}}
		{{else}}
		{{$mName}}{{$st}}State
		{{- end}}
	{{- end}}
	}

	type {{$mName}} struct {
		state {{$mName}}State
	}

	func (m *{{$mName}}) Current() {{$mName}}State {
		return m.state
	}
	
	func (m *{{$mName}}) Operate(operator {{$mName}}Behaviour) {
		for {
			switch m.state {
			{{- range $st, $stDef := .States}}
				{{- if ($stDef.IsTerminal)}}
				case {{$st}}:
					return
				{{- else}}
				case {{$st}}:
					m.handle{{$st}}Event(operator.Operate{{$st}}())
				{{- end}}
			{{- end}}
		}
		}
	}

	func (m *{{$mName}}) Visualize() string {
		return "//Graphviz format \n" +
				"digraph {{$mName}}{\n" +
				{{- range $st, $stDef := .States}}
				{{- if ($stDef.IsTerminal)}}
					"{{$st}} [shape=Msquare];\n" +
				{{- else}}
					{{- range $ev, $dst := $stDef.Events}}
					"{{$st}} -> {{$dst}} [label={{$ev}}];\n" +
					{{- end}}
				{{- end}}
				{{- end}}
				"}\n"
	}

	

	// Handlers for state transitions
	{{range $st, $stDef := .States}}
	{{- if ($stDef.IsTerminal)}}
	{{- else}}
		func (m *{{$mName}}) handle{{$st}}Event(event {{$mName}}{{$st}}Event) {
			switch event {
			{{- range $ev, $dst := $stDef.Events}}
			case {{$st}}{{$ev}}:
				m.state = {{$dst}}
			{{- end}}
			case {{$st}}Noop:
			}
		}
	{{- end}}
	{{end}}

	//--- Here we will define all events ---
	{{range $st, $stDef := .States}}
		{{if ($stDef.IsTerminal)}}
		{{else}}
		//=== {{$mName}}{{$st}}Event definition ===

		type {{$mName}}{{$st}}Event int
		const (
			_ {{$mName}}{{$st}}Event = iota
			{{- range $ev, $dst := $stDef.Events}}
				{{$st}}{{$ev}}
			{{- end}}
			{{$st}}Noop
		)

		var _{{$mName}}{{$st}}EventMap = map[{{$mName}}{{$st}}Event]string{
			{{- range $ev, $dst := $stDef.Events}}
				{{$st}}{{$ev}}: "{{$st}}{{$ev}}",
			{{- end}}
			{{$st}}Noop: "{{$st}}Noop",
		}

		func (m {{$mName}}{{$st}}Event) String() string {
			return _{{$mName}}{{$st}}EventMap[m]
		}

		type {{$mName}}{{$st}}State interface {
			Operate{{$st}}() {{$mName}}{{$st}}Event
		}
		{{end}}
	{{end}}
`))
