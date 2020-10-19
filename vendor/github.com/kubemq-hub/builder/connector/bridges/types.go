package bridges

const (
	promptBindingConfirm = "<cyan>Here is the binding configuration:</>%s"
)

const (
	promptSourceStart = `
<cyan>In the next steps, we will configure the Source connections.
We will set:</>
<yellow>Name -</> A unique name for the Source's binding
<yellow>Kind -</> A Source connection type 
<yellow>Connections -</> A list of connections properties based on the selected kind

<cyan>Lets start binding source configuration:</>`

	promptSourceConfirm = `<cyan>Here is the binding's Source configuration:</>%s`

	promptSourceReconfigure = `<cyan>Lets reconfigure the binding Source:</>`
)
const (
	promptTargetStart = `
<cyan>In the next steps, we will configure the Target connections.
We will set:</>
<yellow>Name -</> A unique name for the Target's binding
<yellow>Kind -</> A Target connection type 
<yellow>Connections -</> A list of connections properties based on the selected kind

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
{{- .SourcesSpec -}}
{{- .TargetsSpec -}}
{{- .PropertiesSpec -}}
`

	promptBindingStartMenu = "<cyan>Lets configure the binding list:</>"
	promptBindingComplete  = "<cyan>We have completed Source and Target binding configurations\n</>"
	promptShowList         = "<cyan>Current Bindings list:</>"
	promptShowBinding      = "<cyan>Showing Binding %s configuration:</>"

	promptBindingDeleteConfirmation = "<red>Binding %s deleted successfully</>"
	promptBindingDeleteCanceled     = "<red>Delete binding operation, cancelled</>"
	promptBindingShowCanceled       = "<red>Show binding operation, cancelled</>"
	promptBindingEditedConfirmation = "<red>Binding %s edited successfully</>"
	promptBindingEditedNoSave       = "<red>Binding %s was not edited</>"
	promptBindingEditCanceled       = "<red>Edit binding operation, cancelled</>"
)
