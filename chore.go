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
	State 	string
	Run 		func() error
	Next 		time.Time
	Timer 	*time.Timer
}

func (c *Chore) Trigger(){
	c.State="running"
	go c.Run()
	expr := cronexpr.MustParse(c.Schedule)
	if expr.Next(time.Now()).IsZero(){
		c.Next = expr.Next(time.Now())
		dur := time.Now().Sub(c.Next)
		c.Timer.Reset(dur)
		c.State=fmt.Sprintf("Scheduled: %s",c.Next.String())
	}else{
		Logger.Debug("invalid schedule",c.Schedule)
		c.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",c.Schedule)
	}
}

// NewResponseFromThinAir returns a new Response object pointing at the general room
func NewResponseFromThinAir(robot *Robot, msg *Message) *Response {
   return &Response{
      Robot: robot,
      Envelope: &hal.Envelope{
         Room: msg.Room,
			User: Config.Name
      },
      Message: &hal.Message{
			Room: 
			Text: `thin air!`
			Type: `chore`
   }
}


// initialize and schedule the chores
func (robot *Robot) Schedule(chores ...*Chore) error{
	for _, c := range chores {
		expr := cronexpr.MustParse(c.Schedule)
		if expr.Next(time.Now()).IsZero(){
			c.Resp = &hal.Response{ 
				Robot: robot,
				Envelope: &hal.Envelope{
					User:robot.Adapter.botname


					
			c.Next = expr.Next(time.Now())
			dur := time.Now().Sub(c.Next)
			c.Timer = time.AfterFunc(dur, c.Trigger) // auto go-routine'd
			c.State=fmt.Sprintf("Scheduled: %s",c.Next.String())
			robot.chores = append(robot.chores, *c)
		}else{
			Logger.Debug("invalid schedule",c.Schedule)
			c.State=fmt.Sprintf("NOT Scheduled (invalid Schedule: %s)",c.Schedule)
	    	return fmt.Errorf("Chore.go: invalid schedule: %v", c.Schedule)
		}
	}
	return nil
}

func KillChore(c *Chore) error{
	c.Timer.Stop()
	return nil
}

func GetChoreByName(name string, robot *Robot) *Chore{
	for _, c := range robot.chores {
		if c.Name == name{
			return &c
		}
	}
	return nil
}


