package hal

import (
	"fmt"
	"time"
	"github.com/gorhill/cronexpr"
)

type Chore struct {
	Name 		string
	Usage 	string
	Sched 	string
	Room		string
	State 	string
	Resp 		*Response
	Run 		func(*Response) error
	Next 		time.Time
	Timer 	*time.Timer
}

func (c *Chore) Trigger(){
	Logger.Debug("Triggered: ",c.Name)
	c.State="running"
	go c.Run(c.Resp)
	StartChore(c)
}

// NewResponseFromThinAir returns a new Response object pointing at the general room
func NewResponseFromThinAir(robot *Robot, room string) *Response {
   return &Response{
      Robot: robot,
      Envelope: &Envelope{
         Room: room,
			User: &User{
				ID: `0`,
				Name: Config.Name,
			},
      },
      Message: &Message{
			Room: room,
			Text: `thin air!`,
			Type: `chore`,
   	},
	}
}

// initialize and schedule the chores
func (robot *Robot) Schedule(chores ...*Chore) error{
	for _, c := range chores {
		c.Resp = NewResponseFromThinAir(robot, c.Room)
		StartChore(c)
		Logger.Debug("appending chore: ",c.Name, " to robot.Chores")
		robot.Chores = append(robot.Chores, c)
	}
	return nil
}

func KillChore(c *Chore) error{
	Logger.Debug(`Stopping: `,c.Name)
	c.State=`Halted by request`
	c.Timer.Stop()
	return nil
}

func StartChore(c *Chore) error{
	Logger.Debug("Re-Starting: ",c.Name)
	expr := cronexpr.MustParse(c.Sched)
	if expr.Next(time.Now()).IsZero(){
		Logger.Debug("invalid schedule",c.Sched)
		c.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",c.Sched)
	}else{
		Logger.Debug("valid Schedule: ",c.Sched)
		c.Next = expr.Next(time.Now())
		dur := c.Next.Sub(time.Now())
			if dur>0{
				Logger.Debug("valid duration: ",dur)
				Logger.Debug("testing timer.. ")
				if c.Timer == nil{
					Logger.Debug("creating new timer ")
					c.Timer = time.AfterFunc(dur, c.Trigger) // auto go-routine'd
				}else{
					Logger.Debug("pre-existing timer found, resetting to: ",dur)
					c.Timer.Reset(dur) // auto go-routine'd
				}
			c.State=fmt.Sprintf("Scheduled: %s",c.Next.String())
			}else{
				Logger.Debug("invalid duration",dur)
				c.State=fmt.Sprintf("Halted. (invalid duration: %s)",dur)
			}
		}
	Logger.Debug("all set! Chore: ",c.Name, "scheduled at: ",c.Next)
	return nil
}

func GetChoreByName(name string, robot *Robot) *Chore{
	for _, c := range robot.Chores {
		if c.Name == name{
			return c
		}else{
			Logger.Debug("chore not found: ",name)
		}
	}
	return nil
}
