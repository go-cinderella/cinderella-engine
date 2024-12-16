package entitymanager

type IdentityLinkEntity struct {
	AbstractEntity
	Rev        *int32
	GroupID    *string
	Type       *string
	UserID     *string
	TaskID     *string
	ProcInstID *string
	ProcDefID  *string
}
