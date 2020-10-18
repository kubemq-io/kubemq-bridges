package target

const (
	promptTargetFirstConnection = "<cyan>Lets add our first connection for kind %s:</>"
)
const targetTemplate = `
<red>targets:</>
  <red>name:</> {{.Name}}
  <red>kind:</> {{.Kind}}
  <red>connections:</>
{{ .ConnectionSpec | indent 2}}
`
