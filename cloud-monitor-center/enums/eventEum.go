package enums

type EventEum string

const (
	InsertAlertContact            EventEum = "insertAlertContact"
	UpdateAlertContact            EventEum = "updateAlertContact"
	CertifyAlertContact           EventEum = "certifyAlertContact"
	DeleteAlertContact            EventEum = "deleteAlertContact"
	InsertAlertContactInformation EventEum = "insertAlertContactInformation"
	UpdateAlertContactInformation EventEum = "updateAlertContactInformation"
	DeleteAlertContactInformation EventEum = "deleteAlertContactInformation"
	InsertAlertContactGroupRel    EventEum = "insertAlertContactGroupRel"
	DeleteAlertContactGroupRel    EventEum = "deleteAlertContactGroupRel"

	InsertAlertContactGroup EventEum = "insertAlertContactGroup"
	UpdateAlertContactGroup EventEum = "updateAlertContactGroup"
	DeleteAlertContactGroup EventEum = "deleteAlertContactGroup"
	//InsertAlertContactGroupRel EventEum = "insertAlertContactGroupRel"
	//DeleteAlertContactGroupRel EventEum = "deleteAlertContactGroupRel"
)
