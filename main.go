package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Process struct {
	id         int
	accounts   [3]int
	channel    chan Message
	state      int
	snapshot   [3]int
	snapshotMux sync.Mutex
	markersReceived int
}

type Message struct {
	from    int
	to      int
	amount  int
}

var numProcesses int
var processes []*Process
var wg sync.WaitGroup

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Print("Enter the number of processes: ")
	fmt.Scan(&numProcesses)

	processes = make([]*Process, numProcesses)
	for i := 0; i < numProcesses; i++ {
		processes[i] = &Process{
			id:       i,
			channel:  make(chan Message),
			state:    0,
		}
		processes[i].accounts = [3]int{rand.Intn(1000), rand.Intn(1000), rand.Intn(1000)}
	}

	initialTotal := getTotalAmount()
	fmt.Println("Initial total amount in the system:", initialTotal)

	for _, p := range processes {
		wg.Add(1)
		go transaction(p)
	}

	wg.Wait()

	finalTotal := getTotalAmount()
	fmt.Println("Final total amount in the system:", finalTotal)
}

func transaction(p *Process) {
	defer wg.Done()
	n := numProcesses * 15

	for i := 1; i <= n; i++ {
		from := p.id
		to := rand.Intn(numProcesses)
		amount := rand.Intn(100)

		p.channel <- Message{from, to, amount}
		time.Sleep(time.Millisecond * 100)

		if i%n == 0 {
			p.snapshotMux.Lock()
			takeSnapshot(p)
			p.snapshotMux.Unlock()
		}
	}

	for p.markersReceived < numProcesses-1 {
		time.Sleep(time.Millisecond * 100)
	}

	close(p.channel)
}

func takeSnapshot(p *Process) {
	p.state++
	p.snapshotMux.Lock()
	p.snapshot = p.accounts
	p.snapshotMux.Unlock()

	for _, q := range processes {
		if q.id != p.id {
			p.channel <- Message{from: p.id, to: q.id, amount: -1}
		}
	}

	for {
		msg, more := <-p.channel
		if !more {
			break
		}
		if msg.amount == -1 {
			p.markersReceived++
			q := processes[msg.from]
			q.snapshotMux.Lock()
			q.snapshot = q.accounts
			q.snapshotMux.Unlock()

			if p.markersReceived == numProcesses-1 {
				break
			}
		} else {
			p.snapshotMux.Lock()
			p.snapshot[msg.to] += msg.amount
			p.snapshotMux.Unlock()
		}
	}

	snapshotTotal := 0
	for _, amount := range p.snapshot {
		snapshotTotal += amount
	}
	fmt.Printf("Snapshot for process %d: %v\n", p.id, p.snapshot)
	fmt.Printf("Total amount in snapshot for process %d: %d\n", p.id, snapshotTotal)

	if snapshotTotal != getTotalAmount() {
		fmt.Printf("Inconsistent snapshot for process %d\n", p.id)
	}
}

func getTotalAmount() int {
	total := 0
	for _, p := range processes {
		for _, amount := range p.accounts {
			total += amount
		}
	}
	return total
}
