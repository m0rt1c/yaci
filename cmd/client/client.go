package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"

	"github.com/yaci/pkg/chord"
)

var (
	join             = flag.Bool("join", false, "Connect to a ring. example: -join -name <name>")
	leave            = flag.Bool("leave", false, "Leave a ring. example: -leave -name <name>")
	new              = flag.Bool("new", false, "Create a ring. example: -new -name <name>")
	base             = flag.Int("base", 2, "Base for ring modulo. Modulo = Base^Exponent - 1")
	exponent         = flag.Int("exponent", 64, "Exponent for ring modulo.")
	port             = flag.Int("port", 6368, "Port for ring's p2p communications. Default 6368. If 0 will be random.")
	timeout          = flag.Int("timeout", 2000, "Time in millisecons for chord ring stabilize.")
	fingerlength     = flag.Int("fingerlength", 5, "Finger table dimensions.")
	nextlength       = flag.Int("nextlength", 4, "Lenght of successors' buffer.")
	lookup           = flag.Bool("lookup", false, "Lookup key in a ring. example: -lookup -name <name> -key <key>")
	list             = flag.Bool("list", false, "List local nodes and rings.")
	finger           = flag.Bool("finger", false, "List also local nodes finger tables.")
	simple           = flag.Bool("simple", false, "Use a simpler and less efficient lookup alghoritm. Included only for completeness.")
	name             = flag.String("name", "homering.ga", "Hostname of a ring.")
	remoteport       = flag.Int("remoteport", 6368, "Port of the host when joining.")
	key              = flag.String("key", "00000", "Key of an item.")
	chordService     = flag.String("csname", "localhost", "Address of the chord service.")
	chordServicePort = flag.Int("csport", 6367, "Port of the chord service.")
)

func printNodeInfo(i chord.NodeInfo) {
	fmt.Println("Node:")
	fmt.Printf("\tID: %v\n", i.ID)
	fmt.Printf("\tAddress: %s\n", i.Address)
	fmt.Printf("\tPort: %d\n", i.Port)
}

func printRingInfo(r chord.RingInfo) {
	fmt.Println("Ring:")
	fmt.Printf("\tName: %s\n", r.Name)
	fmt.Printf("\tModulo: %d\n", r.Modulo)
	fmt.Printf("\tExponent: %d\n", r.ModuloExponent)
	fmt.Printf("\tBase: %v\n", r.ModuloBase)
	fmt.Printf("\tNextLength: %v\n", r.NextBufferLength)
	fmt.Printf("\tTimeout: %d ms\n", r.Timeout)
}

// Dummy Chord Client
func main() {
	flag.Parse()

	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", *chordService, *chordServicePort))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var args chord.ServiceArgs
	var reply chord.ServiceReply

	if *join {
		args.Name = *name
		args.Port = *remoteport
		args.LocalPort = *port

		err = client.Call("Service.JoinRing", args, &reply)
		if err != nil {
			log.Fatal("error:", err)
		}
		printNodeInfo(reply.Node)
		printRingInfo(reply.Ring)
	} else if *new {
		args.Name = *name
		args.Port = *port
		args.Base = *base
		args.Exponent = *exponent
		args.Timeout = *timeout
		args.FingerTableLength = *fingerlength
		args.NextBufferLength = *nextlength

		err = client.Call("Service.CreateRing", args, &reply)
		if err != nil {
			log.Fatal("error:", err)
		}
		printNodeInfo(reply.Node)
		printRingInfo(reply.Ring)
	} else if *leave {
		args.Name = *name

		err = client.Call("Service.Leave", args, &reply)
		if err != nil {
			log.Fatal("error:", err)
		}
		fmt.Println(reply.Message)
	} else if *lookup {
		var method string
		if *simple {
			method = "Service.SimpleLookup"
		} else {
			method = "Service.Lookup"
		}
		args.Name = *name
		args.Key = *key
		err = client.Call(method, args, &reply)
		if err != nil {
			log.Fatal("error:", err)
		}
		fmt.Println(reply.Message)
		printNodeInfo(reply.Node)
	} else if *list {
		err = client.Call("Service.List", args, &reply)
		if err != nil {
			log.Fatal("error:", err)
		}
		for _, n := range reply.List {
			fmt.Println("***********")
			printRingInfo(n.Ring)
			printNodeInfo(n.NodeInfo)
			fmt.Println("Successors:", len(n.Successors))
			for _, next := range n.Successors {
				printNodeInfo(next)
			}
			fmt.Println("Pred")
			printNodeInfo(n.Pred)
			if *finger {

				fmt.Println("Finger table:", len(n.FingerTable))
				for k := range n.FingerTable {
					fmt.Println(k)
					printNodeInfo(n.FingerTable[k])
				}
			}
		}
	}

}
