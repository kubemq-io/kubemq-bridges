package bridges

const (
	promptBindingAddConfirmation = "<cyan>Binding %s was added successfully</>"
	bindingTemplate              = `
<red>name:</> {{.Name}}
{{- .SourcesSpec -}}
{{- .TargetsSpec -}}
{{- .PropertiesSpec -}}
`
	promptBindingComplete           = "<cyan>We have completed Source and Target binding configurations\n</>"
	promptShowBinding               = "<cyan>Showing Binding %s configuration:</>"
	promptBindingDeleteConfirmation = "<cyan>Binding %s deleted successfully\n</>"
	promptBindingEditedConfirmation = "<red>Binding %s edited successfully\n</>"
)
