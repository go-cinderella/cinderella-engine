package entitymanager

type HistoricIdentityLinkEntity struct {
	AbstractEntity
	GroupID    *string
	Type       *string
	UserID     *string
	TaskID     *string
	ProcInstID *string
}
