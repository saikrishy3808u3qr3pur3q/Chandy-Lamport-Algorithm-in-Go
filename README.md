# Chandy-Lamport Algorithm in Go

This repository contains an implementation of the Chandy-Lamport distributed snapshot algorithm in Go. The Chandy-Lamport algorithm is used to capture a consistent global snapshot of a distributed system, even when processes are concurrently sending messages.

## Usage

To run the program, simply compile the Go code and execute the resulting binary: 

### go build -o chandy-lamport
### ./chandy-lamport

The program will prompt you to enter the number of processes in the system. It will then simulate transactions between processes and take snapshots periodically. Finally, it will output the initial and final total amounts in the system, along with the snapshots taken by each process.

## Implementation Details

1)The program creates a number of Process structs, each representing a process in the distributed system.  <br />
2)Each process maintains a set of accounts and a communication channel for receiving messages. <br />
3)Processes randomly send messages to each other, simulating transactions. <br />
4)The takeSnapshot function is called periodically to capture a snapshot of the system's state. <br />
5)The snapshots are used to calculate the total amount of money in the system and ensure consistency. <br />

## Requirements

### Go 1.13 or higher

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
