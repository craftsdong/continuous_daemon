package continuous_daemon

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type CR struct {
	msg []string
	err []string
	num int
}

func (c *CR) Run(msg interface{}) {
	c.num++
	if c.num%2 == 0 {
		panic("error")
	}
	c.msg = append(c.msg, fmt.Sprintf("%s", msg))
}

func (c *CR) Error(err interface{}) {
	c.err = append(c.err, fmt.Sprintf("%s", err))
}

func TestContinuousRunning(t *testing.T) {
	Convey("ContinuousRunning", t, func() {
		cr := &CR{}
		c := NewContinuousRunning(2, 1, cr)
		So(c, ShouldNotBeNil)
		c.Put("test1")
		time.Sleep(100 * time.Millisecond)
		So(len(cr.msg), ShouldEqual, 1)
		So(cr.msg[0], ShouldEqual, "test1")
		So(len(cr.err), ShouldEqual, 0)
		c.Put("test2")
		time.Sleep(100 * time.Millisecond)
		So(len(cr.msg), ShouldEqual, 1)
		So(len(cr.err), ShouldEqual, 1)
		So(cr.err[0], ShouldEqual, "error")
		c.Put("test3")
		time.Sleep(100 * time.Millisecond)
		So(len(cr.msg), ShouldEqual, 2)
		So(len(cr.err), ShouldEqual, 1)
		So(cr.msg[1], ShouldEqual, "test3")
	})
}
