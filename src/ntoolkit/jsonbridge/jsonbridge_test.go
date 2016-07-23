package jsonbridge_test

import (
	"ntoolkit/assert"
	"ntoolkit/jsonbridge"
	"ntoolkit/loopback"
	"testing"
)

type Foo struct {
	Foo int
	Bar int
}

type Bar struct {
	Foo string
	Bar string
}

type Meta struct {
	Type    string
	Channel string
}

type Message1 struct {
	Meta
	Foo
}

type Message2 struct {
	Meta
	Bar
}

func TestNew(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		T.Assert(jsonbridge.New(nil, nil) != nil)
	})
}

func TestReadWrite(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		conn, _ := loopback.New()
		defer conn.Close()

		instance := jsonbridge.New(conn.A, conn.B)

		message1 := Message1{
			Meta{
				Channel: "channel",
				Type:    "Message1",
			},
			Foo{
				Foo: 100,
				Bar: 200,
			},
		}

		message2 := Message2{
			Meta{
				Channel: "channel",
				Type:    "Message2",
			},
			Bar{
				Foo: "100",
				Bar: "200",
			},
		}

		T.Assert(instance.Write(message1) == nil)
		T.Assert(instance.Write(message2) == nil)

		T.Assert(instance.Read() == nil)
		T.Assert(instance.Len() == 2)

		T.Assert(instance.Next() == nil)
		T.Assert(instance.Len() == 1)

		meta := Meta{}
		T.Assert(instance.As(&meta) == nil)
		T.Assert(meta.Channel == "channel")
		T.Assert(meta.Type == "Message1")

		m1 := Message1{}
		T.Assert(instance.As(&m1) == nil)
		T.Assert(m1.Foo.Foo == 100)
		T.Assert(m1.Foo.Bar == 200)

		T.Assert(instance.Next() == nil)
		T.Assert(instance.Len() == 0)

		meta = Meta{}
		T.Assert(instance.As(&meta) == nil)
		T.Assert(meta.Channel == "channel")
		T.Assert(meta.Type == "Message2")

		m2 := Message2{}
		T.Assert(instance.As(&m2) == nil)
		T.Assert(m2.Bar.Foo == "100")
		T.Assert(m2.Bar.Bar == "200")
	})
}
