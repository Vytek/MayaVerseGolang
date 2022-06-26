package main

import (
	"fmt"
	"reflect"

	"github.com/alitto/pond"
	"github.com/lrita/cmap"
	"github.com/rs/xid"
	"gitlab.com/rwxrob/uniq"
)

type Emp struct {
	x int
	y []string
}

var guid xid.ID
var pool *pond.WorkerPool

func main() {
	//Unique ID
	guid = xid.New()
	println(guid.String())
	fmt.Printf("%v\n", guid)

	uid32 := uniq.Base32()
	uuid := uniq.UUID()
	hexid := uniq.Hex(18)
	rgb := uniq.Hex(3)
	b := uniq.Hex(1)

	fmt.Println(uid32) // BCF1KFRJMSAQ1G9HCQ3L25CQOHFIGNQF
	fmt.Println(uuid)  // e98ee42a-d820-4bff-9b0e-67bcff639a17
	fmt.Println(hexid) // 98af788e67de0032b86bb3a3b04f935e72bb
	fmt.Println(rgb)   // 35ba0f
	fmt.Println(b)     // b9

	// Create a buffered (non-blocking) pool that can scale up to 100 workers
	// and has a buffer capacity of 1000 tasks
	pool = pond.New(100, 1000)
	fmt.Printf("t1: %s\n", reflect.TypeOf(pool))

	var list = map[string]*Emp{"e1": {1001, []string{"John", "US"}}}

	e := new(Emp)
	e.x = 1002
	e.y = []string{"Rock", "UK"}

	list["e2"] = e

	fmt.Println(list["e1"])
	fmt.Println(list["e2"])

	var n cmap.Map[string, string]

	// Stores item within map, sets "bar" under key "foo"
	n.Store("foo", "bar")
	n.Store("1", "2")
	n.Store("3", "4")

	tmp2, ok := n.Load("foo")
	if ok {
		fmt.Println(tmp2)
	}

	// Retrieve item from map.
	if tmp, ok := n.Load("foo"); ok {
		bar := tmp
		fmt.Println(bar)
	}

	n.Range(func(key string, value string) bool {
		k, v := key, value
		fmt.Println(k)
		fmt.Println(v)
		return true
	})

	// Deletes item under key "foo"
	n.Delete("foo")
}
