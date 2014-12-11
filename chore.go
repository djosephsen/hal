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
		StartChore(c)
		Logger.Debug("appending chore: ",c.Name, " to robot.Chores")
		robot.Chores = append(robot.Chores, *c)
	}
	return nil
}

func KillChore(c *Chore) error{
	c.Timer.Stop()
	return nil
}

func StartChore(c *Chore) error{
	expr := cronexpr.MustParse(c.Sched)
	if expr.Next(time.Now()).IsZero(){
		Logger.Debug("invalid schedule",c.Sched)
		c.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",c.Sched)
	}else{
		c.Next = expr.Next(time.Now())
		dur := c.Next.Sub(time.Now())
			if dur>0{
				if c.Timer == nil{
					c.Timer = time.AfterFunc(dur, c.Trigger) // auto go-routine'd
				}else{
					c.Timer.Reset(dur) // auto go-routine'd
				}
				c.State=fmt.Sprintf("Chore: %s scheduled at: %s",c.Name,c.Next.String())
			}else{
				Logger.Debug("invalid duration",dur)
				c.State=fmt.Sprintf("NOT Scheduled (invalid duration: %s)",dur)
			}
		}
	return nil
}

func GetChoreByName(name string, robot *Robot) *Chore{
	for _, c := range robot.Chores {
		if c.Name == name{
			return &c
		}
	}
	return nil
}
