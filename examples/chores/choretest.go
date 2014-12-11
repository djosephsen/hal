package bothandlers

import (
	"github.com/danryan/hal"
)

// runs once a minute. Prints 'successful test' to the room of your choosing
var ChoreTest = &hal.Chore{
	Name:  `test-chore`,
	Sched: `0 * * * * * *`,
	Room: `<ENTER YOUR GENERAL ROOM NUMBER`,
	//you can use the 'listRooms' handler in choreHandlers.go to find your room number
	Run: func(res *hal.Response) error {
		return res.Send(`successful test!`)
	},
}
