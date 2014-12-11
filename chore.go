package hal

import (
	"fmt"
	"time"
	"github.com/gorhill/cronexpr"
)

type Chore struct {
	Name 		string
	Usage 	string
	Schedule string
	Room		string
	State 	string
	Resp 		*Response
	Run 		func(*Response) error
	Next 		time.Time
	Timer 	*time.Timer
}

func (c *Chore) Trigger(){
	c.State="running"
	go c.Run(c.Resp)
	expr := cronexpr.MustParse(c.Schedule)
	if expr.Next(time.Now()).IsZero(){
		Logger.Debug("invalid schedule",c.Schedule)
		c.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",c.Schedule)
	}else{
		c.Next = expr.Next(time.Now())
		dur := time.Now().Sub(c.Next)
		c.Timer.Reset(dur)
		c.State=fmt.Sprintf("Scheduled: %s",c.Next.String())
	}
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
		expr := cronexpr.MustParse(c.Schedule)
		if expr.Next(time.Now()).IsZero(){
			Logger.Debug("invalid schedule",c.Schedule)
			c.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",c.Schedule)
			Logger.Debug("appending chore: ",c.Name, " to robot.Chores")
			robot.Chores = append(robot.Chores, *c)
	    	return fmt.Errorf("Chore.go: invalid schedule: %v", c.Schedule)
		}else{
			c.Resp = NewResponseFromThinAir(robot, c.Room)
			c.Next = expr.Next(time.Now())
			dur := time.Now().Sub(c.Next)
			c.Timer = time.AfterFunc(dur, c.Trigger) // auto go-routine'd
			c.State=fmt.Sprintf("Scheduled: %s",c.Next.String())
			Logger.Debug("appending chore: ",c.Name, " to robot.Chores")
			robot.Chores = append(robot.Chores, *c)
		}
	}
	return nil
}

func KillChore(c *Chore) error{
	c.Timer.Stop()
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
