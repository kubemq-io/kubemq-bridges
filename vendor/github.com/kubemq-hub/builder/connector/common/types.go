package common

const (
	promptBindingConfirm     = "<cyan>Here is the binding configuration:</>%s"
	promptBindingReconfigure = "<cyan>Lets reconfigure the binding:</>"
)

const (
	promptSourceStart = `
<cyan>In the next steps, we will configure the Source connection.
We will set:</>
<yellow>Name -</> A unique name for the Source's binding
<yellow>Kind -</> A Source connection type 
<yellow>Properties -</> A list of key/value properties based on the selected kind
<cyan>Lets start binding source configuration:</>`

	promptSourceConfirm = `<cyan>Here is the binding's Source configuration:</>%s`

	promptSourceReconfigure = `<cyan>Lets reconfigure the binding Source:</>`
)
const (
	promptTargetStart = `
<cyan>In the next steps, we will configure the Target connection.
We will set:</>
<yellow>Name -</> A unique name for the Target's binding
<yellow>Kind -</> A Target connection type 
<yellow>Properties -</> A list of key/value properties based on the selected kind
<cyan>Lets start binding target configuration:</>`

	promptTargetConfirm = `<cyan>Here is the binding's Target configuration:</>%s`

	promptTargetReconfigure = `<cyan>Lets reconfigure the binding Target:</>`
)

const (
	promptPropertiesConfirm     = "<cyan>Here is the binding's middleware properties configuration:</>%s"
	promptPropertiesReconfigure = `<cyan>Lets reconfigure the binding middleware properties:</>`
)

const (
	bindingTemplate = `
<red>name:</> {{.Name}}
{{- .SourceSpec -}}
{{- .TargetSpec -}}
{{- .PropertiesSpec -}}
`
	promptBindingComplete = "<cyan>We have completed Source and Target binding configuration</>"
)

const (
	sourceSpecTemplate = `
<red>source:</>
  <red>name:</> {{.Name}}
  <red>kind:</> {{.Kind}}
  <red>properties:</>
{{ .PropertiesSpec | indent 4}}
`
	targetSpecTemplate = `
<red>target:</>
  <red>name:</> {{.Name}}
  <red>kind:</> {{.Kind}}
  <red>properties:</>
{{ .PropertiesSpec | indent 4}}
`
)
