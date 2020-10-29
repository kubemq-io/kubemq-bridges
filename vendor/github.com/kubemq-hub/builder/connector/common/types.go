package common

const (
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
	promptShowSource                = "<cyan>Showing Source configuration:</>"
	promptShowTarget                = "<cyan>Showing Target configuration:</>"

	sourceSpecTemplate = `
<red>source:</>
  <red>kind:</> {{.Kind}}
  <red>properties:</>
{{ .PropertiesSpec | indent 4}}
`
	targetSpecTemplate = `
<red>target:</>
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
