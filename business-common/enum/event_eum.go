package enum

type EventEum string

const (
	InsertContact            EventEum = "insertContact"
	UpdateContact            EventEum = "updateContact"
	DeleteContact            EventEum = "deleteContact"
	ActivateContact          EventEum = "activateContact"
	InsertContactInformation EventEum = "insertContactInformation"
	UpdateContactInformation EventEum = "updateContactInformation"
	DeleteContactInformation EventEum = "deleteContactInformation"
	InsertContactGroupRel    EventEum = "insertContactGroupRel"
	UpdateContactGroupRel    EventEum = "updateContactGroupRel"
	DeleteContactGroupRel    EventEum = "DeleteContactGroupRel"

	InsertContactGroup EventEum = "insertContactGroup"
	UpdateContactGroup EventEum = "updateContactGroup"
	DeleteContactGroup EventEum = "deleteContactGroup"

	CreateRule   EventEum = "create"
	UpdateRule   EventEum = "update"
	ChangeStatus EventEum = "changeStatus"
	DeleteRule   EventEum = "delete"
	BindRule     EventEum = "bind"
	UnbindRule   EventEum = "unbind"

	ChangeMonitorProductStatus EventEum = "changeMonitorProductStatus"
	ChangeMonitorItemDisplay   EventEum = "changeMonitorItemDisplay"
)
