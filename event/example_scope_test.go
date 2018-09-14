// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package event_test

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/event"
	"time"
)

//一次关闭多个订阅
func ExampleSubscriptionScope_Track() {
	var feed1 event.Feed
	var feed2 event.Feed
	var feed3 event.Feed
	var feedscope event.SubscriptionScope
	ch := make(chan int)
	//跟踪订阅的事件
	fs1 := feedscope.Track(feed1.Subscribe(ch))
	fs2 := feedscope.Track(feed2.Subscribe(ch))
	fs3 := feedscope.Track(feed3.Subscribe(ch))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case v := <-ch:
				fmt.Println(v)
			case <-fs1.Err():
				fmt.Println("feed1 close")
				return
			case <-fs2.Err():
				fmt.Println("feed2 close")
				return
			case <-fs3.Err():
				fmt.Println("feed3 close")
				return
			}
		}
	}()

	fmt.Println("subscribe feed count", feedscope.Count())

	feed1.Send(1)
	feed2.Send(2)
	feed3.Send(3)
	//等待一段时间发送处理完后再关闭
	time.Sleep(time.Second)
	fs1.Unsubscribe()
	fmt.Println("subscribe feed count", feedscope.Count())

	feed1.Send(4)
	time.Sleep(time.Second)

	//一次关闭所有订阅事件
	feedscope.Close()
	wg.Wait()

	// Output:
	// subscribe feed count 3
	// 1
	// 2
	// 3
	// subscribe feed count 2
	// feed1 close
}

// This example demonstrates how SubscriptionScope can be used to control the lifetime of
// subscriptions.
//此示例演示了如何使用SubscriptionScope控制订阅的生命周期。
//
// Our example program consists of two servers, each of which performs a calculation when
// requested. The servers also allow subscribing to results of all computations.
// 包含两个计算服务，通过订阅时间获取计算结果
type divServer struct{ results event.Feed }
type mulServer struct{ results event.Feed }

func (s *divServer) do(a, b int) int {
	r := a / b
	s.results.Send(r)
	return r
}

func (s *mulServer) do(a, b int) int {
	r := a * b
	s.results.Send(r)
	return r
}

// The servers are contained in an App. The app controls the servers and exposes them
// through its API.
// 服务包含在APP中。app负责控制服务和暴漏API
type App struct {
	divServer
	mulServer
	scope event.SubscriptionScope
}

func (s *App) Calc(op byte, a, b int) int {
	switch op {
	case '/':
		return s.divServer.do(a, b)
	case '*':
		return s.mulServer.do(a, b)
	default:
		panic("invalid op")
	}
}

// The app's SubscribeResults method starts sending calculation results to the given
// channel. Subscriptions created through this method are tied to the lifetime of the App
// because they are registered in the scope.
func (s *App) SubscribeResults(op byte, ch chan<- int) event.Subscription {
	switch op {
	case '/':
		return s.scope.Track(s.divServer.results.Subscribe(ch))
	case '*':
		return s.scope.Track(s.mulServer.results.Subscribe(ch))
	default:
		panic("invalid op")
	}
}

// Stop stops the App, closing all subscriptions created through SubscribeResults.
func (s *App) Stop() {
	s.scope.Close()
}

func ExampleSubscriptionScope() {
	// Create the app.
	var (
		app  App
		wg   sync.WaitGroup
		divs = make(chan int)
		muls = make(chan int)
	)

	// Run a subscriber in the background.
	divsub := app.SubscribeResults('/', divs)
	mulsub := app.SubscribeResults('*', muls)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("subscriber exited")
		defer divsub.Unsubscribe()
		defer mulsub.Unsubscribe()
		for {
			select {
			case result := <-divs:
				fmt.Println("division happened:", result)
			case result := <-muls:
				fmt.Println("multiplication happened:", result)
			case <-divsub.Err():
				return
			case <-mulsub.Err():
				return
			}
		}
	}()

	// Interact with the app.
	app.Calc('/', 22, 11)
	app.Calc('*', 3, 4)

	// Stop the app. This shuts down the subscriptions, causing the subscriber to exit.
	app.Stop()
	wg.Wait()

	// Output:
	// division happened: 2
	// multiplication happened: 12
	// subscriber exited
}
