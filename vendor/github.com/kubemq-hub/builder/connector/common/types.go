package common

const (
	promptSourceStart = `
<cyan>In the next steps, we will configure the Source connection.
We will set:</>
<yellow>Name -</> A unique name for the Source's binding
<yellow>Kind -</> A Source connection type 
<yellow>Properties -</> A list of key/value properties based on the selected kind
<cyan>Lets start binding source configuration:</>`

	promptTargetStart = `
<cyan>In the next steps, we will configure the Target connection.
We will set:</>
<yellow>Name -</> A unique name for the Target's binding
<yellow>Kind -</> A Target connection type 
<yellow>Properties -</> A list of key/value properties based on the selected kind
<cyan>Lets start binding target configuration:</>`

	bindingTemplate = `
<red>name:</> {{.Name}}
{{- .SourceSpec -}}
{{- .TargetSpec -}}
{{- .PropertiesSpec -}}
`

	promptBindingComplete           = "<cyan>We have completed Source and Target binding configuration</>"
	promptBindingDeleteConfirmation = "<cyanBinding %s deleted successfully</>"
	promptBindingEditedConfirmation = "<cyan>Binding %s edited successfully</>"
	promptBindingAddConfirmation    = "<cyan>Binding %s was added successfully</>"
	promptShowBinding               = "<cyan>Showing Binding %s configuration:</>"
	promptShowSource                = "<cyan>Showing Source %s configuration:</>"
	promptShowTarget                = "<cyan>Showing Target %s configuration:</>"

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

	connectorTemplate = `
<red>kind:</> {{.Kind}}
<red>description:</> {{.Description}}
<red>properties:</>
{{ .PropertiesSpec | indent 2}}
`
)
