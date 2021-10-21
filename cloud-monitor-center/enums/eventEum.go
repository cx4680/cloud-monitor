package enums

type EventEum string

const (
	InsertAlertContact                    EventEum = "insertAlertContact"
	UpdateAlertContact                    EventEum = "updateAlertContact"
	DeleteAlertContact                    EventEum = "deleteAlertContact"
	CertifyAlertContact                   EventEum = "certifyAlertContact"
	InsertAlertContactInformation         EventEum = "insertAlertContactInformation"
	DeleteAlertContactInformation         EventEum = "deleteAlertContactInformation"
	InsertAlertContactGroupRel            EventEum = "insertAlertContactGroupRel"
	DeleteAlertContactGroupRelByContactId EventEum = "DeleteAlertContactGroupRelByContactId"

	InsertAlertContactGroup             EventEum = "insertAlertContactGroup"
	UpdateAlertContactGroup             EventEum = "updateAlertContactGroup"
	DeleteAlertContactGroup             EventEum = "deleteAlertContactGroup"
	DeleteAlertContactGroupRelByGroupId EventEum = "DeleteAlertContactGroupRelByGroupId"

	CreateRule  EventEum = "createRule"
	UpdateRule  EventEum = "updateRule"
	DisableRule EventEum = "disableRule"
	EnableRule  EventEum = "enableRule"
	DeleteRule  EventEum = "deleteRule"
	BindRule    EventEum = "bindRule"
	UnbindRule  EventEum = "unbindRule"
)
