package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type nodeInfo struct {
	isActive bool
	conn     net.Conn
	enc      *gob.Encoder
	dec      *gob.Decoder
}

type TransactionWrapper struct {
	Tr             Transaction
	Deliverable    bool
	OriginalNodeID string
	ProposedNodeID string
}

type Transaction struct {
	TransactionType byte   // either 'R' (transfer) or 'D' (deposit)
	Money           int    // number of money to be add into Account1
	Account1        string // e.g. "abcde"
	Account2        string // if type is 'R', money will also be subtracted from this account
}

type Message struct {
	NodeID     string // id of the node at another end of the connection
	Identifier string // e.g. "2_36" is the 36th message sent by node2
	MsgType    byte   // either 'T' (transaction) or 'A' (agreed#) or 'P' (proposed#), indicating the type of the message

	/* A transaction msg will have following additional fields */
	Tr Transaction

	/* A agreed sequence number msg will have following additional fields */
	AgreedSeqNum int
	AgreedNodeID string

	/* A proposed sequence number msg will have following additional fields */
	ProposedSeqNum int
}

func parsingConfigs(fileName string) (int, [][]string) {
	f, e := os.Open(fileName)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()

	var configContent_ [][]string
	var numOtherNodes_ int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if numOtherNodes_ == 0 {
			numOtherNodes_, e = strconv.Atoi(scanner.Text())
			if e != nil {
				log.Fatal(e)
			}
			continue
		}
		nodeinfo := strings.Split(scanner.Text(), " ")
		configContent_ = append(configContent_, nodeinfo)
	}
	return numOtherNodes_, configContent_
}

func setupNodeMap() {
	for i := 0; i < numOtherNodes; i++ {
		nodeMap[configContent[i][0]] = &nodeInfo{}
	}
}

func setupNodeInfo(conn net.Conn) string {
	dec := gob.NewDecoder(conn)
	msg := Message{}
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("Error getting nodeID: ", err)
	}
	nodeID := msg.NodeID
	// log.Println("The nodeID we get is " + nodeID)
	nodeMap[nodeID].conn = conn
	nodeMap[nodeID].isActive = true
	nodeMap[nodeID].enc = gob.NewEncoder(conn)
	nodeMap[nodeID].dec = dec
	// log.Println("Now listening to " + nodeID)
	return nodeID
}

func printBalance() {
	keys := make([]string, 0, len(accounts))
	for k := range accounts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Print("BALANCES")
	for _, k := range keys {
		fmt.Print(" " + k + ":" + strconv.Itoa(accounts[k]))
	}
	fmt.Println()
}

// func printBalance() {
// 	keys := make([]string, 0, len(accounts))
// 	for k := range accounts {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
// 	fmt.Print("BALANCES")
// 	fmt.Fprint(os.Stderr, "BALANCES")
// 	for _, k := range keys {
// 		fmt.Print(" " + k + ":" + strconv.Itoa(accounts[k]))
// 		fmt.Fprint(os.Stderr, " "+k+":"+strconv.Itoa(accounts[k]))
// 	}
// 	fmt.Println()
// 	fmt.Fprintln(os.Stderr)
// }

func getTimeStr() string {
	var secs float64 = float64(time.Now().UnixNano()) / 1000000000
	return strconv.FormatFloat(secs, 'f', -1, 64)
}
