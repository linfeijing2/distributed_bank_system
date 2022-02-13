package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const timeout = 10
const msgSize = "120"

var (
	myID          string
	numOtherNodes int
	configContent [][]string
	msgCounter    int
	psn           int
	mu            sync.Mutex
	bandWidth     *os.File
	processTime   *os.File

	nodeMap     = make(map[string]*nodeInfo)
	msgHistory  = make(map[string]bool)
	msgMap      = make(map[string]*TransactionWrapper)
	hbQueue     = NewPriorityQueue()
	msgPSNQueue = make(map[string]map[string]int)
	accounts    = make(map[string]int)
)

func main() {
	myID = os.Args[1]
	listenReady := make(chan bool)
	numOtherNodes, configContent = parsingConfigs(os.Args[3])
	setupNodeMap()
	bwpath, processpath := "", ""

	if numOtherNodes > 2 {
		bwpath = os.Args[1] + "_Bandwidth.txt"
		processpath = os.Args[1] + "_Processing.txt"
	} else {
		bwpath = os.Args[1] + "_small_Bandwidth.txt"
		processpath = os.Args[1] + "_small_Processing.txt"
	}

	tmp, err := os.Create(bwpath)
	if err != nil {
		log.Fatal("Fail to create " + bwpath)
	}
	bandWidth = tmp
	tmp, err = os.Create(processpath)
	if err != nil {
		log.Fatal("Fail to create " + processpath)
	}
	processTime = tmp
	defer bandWidth.Close()
	defer processTime.Close()

	go listenNodes(listenReady) // will only receive connection request from nodes with smaller ID

	dialNodesWithBiggerID() // dial to nodes with bigger ID

	<-listenReady
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.Split(scanner.Text(), " ")
		trw := TransactionWrapper{Deliverable: false, OriginalNodeID: myID, ProposedNodeID: myID}

		if input[0][0] == 'D' {
			tmp, err := strconv.Atoi(input[2])
			if err != nil {
				log.Fatal(err)
			}
			trw.Tr = Transaction{TransactionType: 'D', Money: tmp, Account1: input[1]}
		} else {
			tmp, err := strconv.Atoi(input[4])
			if err != nil {
				log.Fatal(err)
			}
			trw.Tr = Transaction{TransactionType: 'R', Money: tmp, Account1: input[3], Account2: input[1]}
		}

		msgCounter++
		msgIdentifier := myID + "_" + strconv.Itoa(msgCounter)
		mu.Lock()
		msgHistory[msgIdentifier] = true
		msgMap[msgIdentifier] = &trw
		psn++
		hbQueue.Insert(msgIdentifier, psn)
		// log.Println("Proposing " + strconv.Itoa(psn) + " for " + msgIdentifier)
		msgPSNQueue[msgIdentifier] = make(map[string]int)
		msgPSNQueue[msgIdentifier][myID] = psn
		mu.Unlock()
		multicastTransaction(msgIdentifier, trw.Tr)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func listenNodes(listenReady chan bool) {
	numOfSmallerID := 0
	for i := 0; i < numOtherNodes; i++ {
		if configContent[i][0] < myID {
			numOfSmallerID++
		}
	}

	listener, err := net.Listen("tcp", ":"+os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for ; numOfSmallerID > 0; numOfSmallerID-- {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		nodeID := setupNodeInfo(conn)
		go readConn(nodeID)
	}
	listenReady <- true
}

func dialNodesWithBiggerID() {
	for i := 0; i < numOtherNodes; i++ {
		if configContent[i][0] > myID {
			for !nodeMap[configContent[i][0]].isActive {
				conn, err := net.Dial("tcp", configContent[i][1]+":"+configContent[i][2])
				if err != nil {
					// log.Println(err)
					continue
				}
				nodeMap[configContent[i][0]].conn = conn
				nodeMap[configContent[i][0]].isActive = true
				nodeMap[configContent[i][0]].enc = gob.NewEncoder(conn)
				nodeMap[configContent[i][0]].dec = gob.NewDecoder(conn)
				err = nodeMap[configContent[i][0]].enc.Encode(Message{NodeID: myID})
				if err != nil {
					log.Fatal(err)
				}
				go readConn(configContent[i][0])
			}
		}
	}
}

func readConn(nodeID string) {
	enc := nodeMap[nodeID].enc
	dec := nodeMap[nodeID].dec
	for {
		msg := Message{}
		err := dec.Decode(&msg)
		if err != nil {
			// log.Println("Error reading "+nodeID+" :", err)
			setNodeToInactive(nodeID)
			return
		}
		bandWidth.WriteString(msgSize + " " + getTimeStr() + "\n")
		bandWidth.Sync()

		if msg.MsgType == 'T' {
			mu.Lock()
			msgHistory[msg.Identifier] = true
			msgMap[msg.Identifier] = &TransactionWrapper{Tr: msg.Tr, Deliverable: false, OriginalNodeID: nodeID, ProposedNodeID: myID}
			psn++
			hbQueue.Insert(msg.Identifier, psn)
			msgPSN := psn
			mu.Unlock()

			err := enc.Encode(Message{NodeID: myID, Identifier: msg.Identifier, MsgType: 'P', ProposedSeqNum: msgPSN})
			if err != nil {
				// log.Println("Error proposing to "+nodeID+" :", err)
				setNodeToInactive(nodeID)
				return
			}
			bandWidth.WriteString(msgSize + " " + getTimeStr() + "\n")
			bandWidth.Sync()
			// log.Println("Proposing " + strconv.Itoa(msgPSN) + " for " + msg.Identifier)
		} else if msg.MsgType == 'A' {
			mu.Lock()
			if msgMap[msg.Identifier] == nil || msgMap[msg.Identifier].Deliverable {
				mu.Unlock()
				continue
			}
			if msg.AgreedSeqNum > psn {
				psn = msg.AgreedSeqNum
			}
			msgMap[msg.Identifier].Deliverable = true
			msgMap[msg.Identifier].ProposedNodeID = msg.AgreedNodeID
			hbQueue.UpdatePriority(msg.Identifier, msg.AgreedSeqNum)
			checkAndDeliver()
			mu.Unlock()
			multicastAgreed(msg.Identifier, msg.AgreedSeqNum, msg.AgreedNodeID)
		} else { // 'P'
			mu.Lock()
			_, ok := msgPSNQueue[msg.Identifier]
			if ok {
				msgPSNQueue[msg.Identifier][nodeID] = msg.ProposedSeqNum
				checkAndMulticastASN(msg.Identifier)
			}
			mu.Unlock()
		}
	}
}

func checkAndDeliver() {
	for {
		msgIdentifier, err := hbQueue.Peek()
		if err != nil {
			return
		} // queue might be empty
		if msgMap[msgIdentifier].Deliverable {
			// tmp, _ := hbQueue.PeekPriority()
			// log.Println("Now delivering " + msgIdentifier + " with priority " + strconv.Itoa(tmp))
			// hbQueue.PrintQueue()
			deliver(msgIdentifier)
			processTime.WriteString(msgIdentifier + " " + getTimeStr() + "\n")
			processTime.Sync()
			hbQueue.Pop()
			delete(msgMap, msgIdentifier)
		} else {
			return
		}
	}
}

func deliver(msgIdentifier string) {
	tr := msgMap[msgIdentifier].Tr
	if tr.TransactionType == 'R' {
		balance2, ok := accounts[tr.Account2]
		if !ok {
			// log.Println(msgIdentifier + " is illegal transaction. The transfer account doesn't exist.")
			printBalance()
			return
		}
		if balance2 < tr.Money {
			// log.Println(msgIdentifier + " is illegal transaction. The transfer account doesn't have enough money.")
			printBalance()
			return
		}
		accounts[tr.Account2] = balance2 - tr.Money
	}
	accounts[tr.Account1] += tr.Money
	printBalance()
}

func checkAndMulticastASN(msgIdentifier string) {
	maxProposed := msgPSNQueue[msgIdentifier][myID]
	maxProposedNodeID := myID
	for i := range nodeMap {
		if nodeMap[i].isActive {
			if msgPSNQueue[msgIdentifier][i] == 0 {
				return
			}
			if msgPSNQueue[msgIdentifier][i] == maxProposed && i > maxProposedNodeID {
				maxProposedNodeID = i
				continue
			}
			if msgPSNQueue[msgIdentifier][i] > maxProposed {
				maxProposed = msgPSNQueue[msgIdentifier][i]
				maxProposedNodeID = i
			}
		}
	}
	delete(msgPSNQueue, msgIdentifier)
	if maxProposed > psn {
		psn = maxProposed
	}
	msgMap[msgIdentifier].Deliverable = true
	msgMap[msgIdentifier].ProposedNodeID = maxProposedNodeID
	hbQueue.UpdatePriority(msgIdentifier, maxProposed)
	checkAndDeliver()
	multicastAgreed(msgIdentifier, maxProposed, maxProposedNodeID)
}

func multicastTransaction(msgIdentifier string, tr Transaction) {
	multicast(&Message{NodeID: myID, Identifier: msgIdentifier, MsgType: 'T', Tr: tr})
}

func multicastAgreed(msgIdentifier string, agreedSeqNum int, agreedNodeID string) {
	multicast(&Message{NodeID: myID, Identifier: msgIdentifier, MsgType: 'A', AgreedSeqNum: agreedSeqNum, AgreedNodeID: agreedNodeID})
}

func multicast(msg *Message) {
	for i := range nodeMap {
		if nodeMap[i].isActive {
			enc := nodeMap[i].enc
			err := enc.Encode(*msg)
			if err != nil {
				// log.Println("Error multicasting to "+i+" :", err)
				go setNodeToInactive(i)
				continue
			}
			bandWidth.WriteString(msgSize + " " + getTimeStr() + "\n")
			bandWidth.Sync()
		}
	}
}

func setNodeToInactive(nodeID string) {
	mu.Lock()
	if !nodeMap[nodeID].isActive {
		mu.Unlock()
		return
	}
	defer nodeMap[nodeID].conn.Close()
	nodeMap[nodeID].isActive = false

	for msgIdentifier := range msgPSNQueue {
		checkAndMulticastASN(msgIdentifier)
	}
	mu.Unlock()

	cleanHBQueue(nodeID)
}

func cleanHBQueue(nodeID string) {
	time.Sleep(timeout * time.Second)
	mu.Lock()
	for msgIdentifier, trw := range msgMap {
		if trw.OriginalNodeID == nodeID && !trw.Deliverable {
			delete(msgMap, msgIdentifier)
			hbQueue.Delete(msgIdentifier)
		}
	}
	mu.Unlock()
}

// func main1() {
// 	msgMap["node1_101"] = &TransactionWrapper{ProposedNodeID: "node1"}
// 	msgMap["node7_33"] = &TransactionWrapper{ProposedNodeID: "node7"}
// 	hbQueue.Insert("node1_101", 22)
// 	hbQueue.Insert("node7_33", 22)
// 	hbQueue.Insert("node3_33", 25)
// 	x, _ := hbQueue.Peek()
// 	println(x)
// 	hbQueue.PrintQueue()
// 	hbQueue.Pop()
// 	x, _ = hbQueue.Peek()
// 	println(x)
// 	hbQueue.PrintQueue()
// 	hbQueue.Delete("node7_33")
// 	x, _ = hbQueue.Peek()
// 	println(x)
// 	hbQueue.PrintQueue()
// }
