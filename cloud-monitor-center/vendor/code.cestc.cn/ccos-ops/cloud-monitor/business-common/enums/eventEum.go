package enums

type EventEum string

const (
	InsertAlertContact            EventEum = "insertAlertContact"
	UpdateAlertContact            EventEum = "updateAlertContact"
	DeleteAlertContact            EventEum = "deleteAlertContact"
	CertifyAlertContact           EventEum = "certifyAlertContact"
	InsertAlertContactInformation EventEum = "insertAlertContactInformation"
	UpdateAlertContactInformation EventEum = "updateAlertContactInformation"
	DeleteAlertContactInformation EventEum = "deleteAlertContactInformation"
	InsertAlertContactGroupRel    EventEum = "insertAlertContactGroupRel"
	UpdateAlertContactGroupRel    EventEum = "updateAlertContactGroupRel"
	DeleteAlertContactGroupRel    EventEum = "DeleteAlertContactGroupRel"

	InsertAlertContactGroup EventEum = "insertAlertContactGroup"
	UpdateAlertContactGroup EventEum = "updateAlertContactGroup"
	DeleteAlertContactGroup EventEum = "deleteAlertContactGroup"

	CreateRule   EventEum = "create"
	UpdateRule   EventEum = "update"
	ChangeStatus EventEum = "changeStatus"
	DeleteRule   EventEum = "delete"
	BindRule     EventEum = "bind"
	UnbindRule   EventEum = "unbind"
)
