package source

const (
	promptSourceFirstConnection = "<cyan>Lets add our first connection for kind %s:</>"
	promptShowSource            = "<cyan>Showing Sources configuration:</>"
)
const sourceTemplate = `
<red>sources:</>
  <red>kind:</> {{.Kind}}
  <red>connections:</>
{{ .ConnectionSpec | indent 4}}
`
