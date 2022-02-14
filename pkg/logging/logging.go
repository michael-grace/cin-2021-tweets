package logging

import "fmt"

func Error(e error) {
	fmt.Printf("ERROR: %s\n", e)
}

type Action string

const (
	ApproveTweet    Action = "APPROVE"
	RejectTweet     Action = "REJECT"
	BlockUser       Action = "BLOCK"
	UnblockUser     Action = "UNBLOCK"
	RemoveFromBoard Action = "REMOVE FROM BOARD"
)

func LogAction(action Action, message string) {
	fmt.Printf("%s: %s\n", action, message)
}
