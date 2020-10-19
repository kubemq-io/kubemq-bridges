package source

const (
	promptSourceFirstConnection = "<cyan>Lets add our first connection for kind %s:</>"
	promptShowSource            = "<cyan>Showing Sources %s configuration:</>"
)
const sourceTemplate = `
<red>sources:</>
  <red>name:</> {{.Name}}
  <red>kind:</> {{.Kind}}
  <red>connections:</>
{{ .ConnectionSpec | indent 4}}
`
