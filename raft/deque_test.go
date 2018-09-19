package raft

import (
	"gopkg.in/oleiade/lane.v1"
	"testing"
)
// 循环队列
func TestDeque(t *testing.T) {
	mdeque := lane.NewDeque()
	mdeque.Append("1")
	mdeque.Prepend("2")
	mdeque.Append("3")
	mdeque.Prepend("4")
	t.Log("first ", mdeque.First())
	t.Log("last", mdeque.Last())
	t.Log("capacity", mdeque.Capacity())
	t.Log("size", mdeque.Size())
}
