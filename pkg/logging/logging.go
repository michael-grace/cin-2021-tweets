/**
URY Tweet Board
Copyright (C) 2022 Michael Grace <michael.grace@ury.org.uk>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
