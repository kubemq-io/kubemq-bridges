package source

const (
	promptSourceFirstConnection = "<cyan>Lets add our first connection for kind %s:</>"
)
const sourceTemplate = `
<red>sources:</>
  <red>name:</> {{.Name}}
  <red>kind:</> {{.Kind}}
  <red>connections:</>
{{ .ConnectionSpec | indent 4}}
`
