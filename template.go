package main

import "text/template"

var embeddedTemplate = template.Must(template.New("embedded").Parse(`
	package {{.PkgName}}
	
	import "fmt"

	// Generated by go-fsm-generator. DO NOT EDIT.

	//+++ General machine definition +++
	{{$mName := .MachineName}}
	//  {{$mName}} states
	type {{$mName}}State int
	const (
		_ {{$mName}}State = iota
		{{- range $st, $stDef := .States}}
			{{$st}} // {{$st}} state
		{{- end}}
	)

	var _{{$mName}}StateMap = map[{{$mName}}State]string{
		{{- range $st, $stDef := .States}}
			{{$st}}: "{{$st}}",
		{{- end}}
	}

	var _{{$mName}}ParsingStateMap = map[string]{{$mName}}State{
		{{- range $st, $stDef := .States}}
			"{{$st}}": {{$st}},
		{{- end}}
	}

	func (s {{$mName}}State) String() string {
		return _{{$mName}}StateMap[s]
	}

	// {{$mName}} behaviours
	type {{$mName}}Behaviour interface {
	{{- range $st, $stDef := .States}}
		{{- if ($stDef.IsTerminal)}}
		{{else}}
		{{$mName}}{{$st}}State
		{{- end}}
	{{- end}}
	}
	
	// {{$mName}} machine type
	type {{$mName}} struct {
		state {{$mName}}State
	}
	
	// {{$mName}} creates machine with specified initial state
	func New{{$mName}}(state {{$mName}}State) *{{$mName}} {
		return &{{$mName}}{state: state}
	}

	// {{$mName}} can be used to deserialize  machine state
	func New{{$mName}}FromString(stateStr string) (*{{$mName}}, error) {
		state, ok := _{{$mName}}ParsingStateMap[stateStr]
		if !ok {
			return nil, fmt.Errorf("state unknown for {{$mName}}: %s", stateStr)
		}
		return &{{$mName}}{state: state}, nil
	}

	// Current returns current state of {{$mName}}
	func (m *{{$mName}}) Current() {{$mName}}State {
		return m.state
	}
	
	// Operate executes behaviour for the current state {{$mName}}
	func (m *{{$mName}}) Operate(operator {{$mName}}Behaviour) {
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

	// Visualize states and events for {{$mName}} in Graphviz format
	func (m *{{$mName}}) Visualize() string {
		return {{.Description}}
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
		
		// {{$mName}}{{$st}}Event definition
		type {{$mName}}{{$st}}Event int
		const (
			_ {{$mName}}{{$st}}Event = iota
			{{- range $ev, $dst := $stDef.Events}}
				{{$st}}{{$ev}} // {{$st}}{{$ev}} -> {{$dst}}
			{{- end}}
			{{$st}}Noop // remain in {{$st}}
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
		
		// {{$mName}}{{$st}}State behaviour
		type {{$mName}}{{$st}}State interface {
			Operate{{$st}}() {{$mName}}{{$st}}Event
		}
		{{end}}
	{{end}}
`))
