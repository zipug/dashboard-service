package models

type ChatMember struct {
	ID        int
	FirstName string
	LastName  string
	Username  string
	Bio       string
	ApiToken  string
	PhotoUrl  string
	Photo     struct {
		SmallFileID       string
		SmallFileUniqueID string
		BigFileID         string
		BigFileUniqueID   string
	}
}
