package handlers

import (
	"fmt"
	"github.com/danryan/hal"
	"time"
	"strings"
)

// Ask your bot for the ID of the current Room
var ListRooms = &hal.Handler{
	Method:  hal.RESPOND,
	Pattern: `(what room)|(list *room)|(room *list)`,
	Run: func(res *hal.Response) error {
		room := res.Message.Room
		reply := fmt.Sprintf("Current room is: %s",room)
		return res.Send(reply)
	},
}

// List all initialized chores (running or not)
var ListChores = &hal.Handler{
	Method:  hal.RESPOND,
	Pattern: `(list chores)|(chore list)`,
	Run: func(res *hal.Response) error {
		var reply string
		if len(res.Robot.Chores) == 0{
			reply=`No chores have been registered (sorry?)`
		}else{
			reply=`Name  :small_blue_diamond:  Schedule  :small_blue_diamond:  Firing in  :small_blue_diamond: Current State`
			for _,c := range res.Robot.Chores{
				reply = fmt.Sprintf("%s\n%s:small_blue_diamond:%s:small_blue_diamond:%v:small_blue_diamond:%s",reply,c.Name, c.Sched, c.Next.Sub(time.Now()), c.State)
			}
		}
		return res.Reply(reply)
	},
}

// Stop or start individual chores by name
var ManageChores = &hal.Handler{
	Method:  hal.RESPOND,
	Pattern: `(start|stop) chore .*`,
	Run: func(res *hal.Response) error {
		var reply string
		cname:=strings.SplitAfterN(res.Match[0],` `,4)
		c:=hal.GetChoreByName(cname[3],res.Robot)
		if c == nil{ 
			reply = fmt.Sprintf("Chore not found: %s",(cname[3]))
		}else{
			if cname[1]==`stop `{
				hal.KillChore(c)
			}else{
				hal.StartChore(c)
			}
			reply = fmt.Sprintf("%s\n%s:small_blue_diamond:%s:small_blue_diamond:%v:small_blue_diamond:%s",reply,c.Name, c.Sched, c.Next.Sub(time.Now()), c.State)
		}
		return res.Reply(reply)
	},
}
