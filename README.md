# distributed_bank_system

This project is implemented via Go programming language and Python(for the analysis). If you would like to test the code, please make sure you already have Go installed on your computer. We are using Go of version 1.17.1 in this project.

## Building and Running
After you clone this repository to your computer, please follow the steps listed below to test and obtain the graph.

**Build**

Build the code for executable files by following command in the root directory:
```
go build
```
**Node**
The node in this project represent each branch of a bank. For different bide(branch), we will have different operations.
This project tests the functionality under different scenarios. Each scenario requires a different number of nodes and transaction generation frequency.

1. 3 nodes, 0.5 Hz each, running for 100 seconds
2. 8 nodes, 5 Hz each, running for 100 seconds
3. 3 nodes, 0.5 Hz each, runing for 100 seconds, then one node fails, and the rest continue to run for 100 seconds
4. 8 nodes, 5 Hz each, running for 100 seconds, then 3 nodes fail simultaneously, and the rest continue to run for 100 seconds.

Run the following general command to pipe the output of "gentx.py" to each node and start running the program.

```
python3 -u gentx.py [frequency] | ./mp1_node [node name] [port name] [config file name]
(i.e. python3 -u gentx.py 0.5 | ./mp1_node node1 1234 config1.txt)
```
We have tested the functionality of this project by running the shell script in the virtual machine (VM). We put the shell scripts in the *remotesh* directory. The shell script will kill the processes after 100s or 200s, depending on the testing scenario. To test with the shell script, run the following command for 3 nodes, 3 nodes with failure, 8 nodes, and 8 nodes with failure in VM's terminal, respectively. 

```
sh 3nodes1/2/3.sh / sh f3nodes1/2/3.sh / sh 8nodes1/2/3/4/5/6/7/8.sh / sh f8nodes1/2/3/4/5/6/7/8.sh
(i.e. sh f8nodes5.sh for node5 under scenario of 8 nodes with failure)
```
**Analysis Graph Generation**
The graphs for each scenario are in the graph directory under the stat directory. To generate the graphs, first put all the *.txt* files of data collected in the program in the stat directory, and run the *Graphs.ipynb*. In the *Graphs.ipynb*, run the graph generation block under different labels for different graphs. The graphs will be saved in the graph directory automatically. 

**Design Document**

* Scope 

	In this project, we are supposed to design a transaction back system. In this system, we need to keep track of each account activity. Each account with an initial balance of 0 can either deposit to his/her account or transfer money to other accounts. Each account should have a non-negative integer balance after the transaction. Otherwise, the transaction will not be allowed. After each transaction, the balance will be printed with the account name sorted alphabetically. 

* Algorithm Design

  We implemented the ISIS algorithm for this MP. It is an asynchronous and decentralized algorithm that ensures the total ordering of transactions/messages. Every pair of the node is connected and able to communicate. Each node will re-multicast the agreed sequence number (ASN) to ensure reliable ordering. 

  The ISIS algorithm we used can be explained as below:

  1. Node A multicast message m to all other nodes.
  2. Each node X replies to A with a proposed sequence number (PSN) with its node id, avoiding the situation where two nodes propose the same PSN. Each node X will put this message with the PSN to the hold-back queue, which is ordered with the smallest PSN at the front. In the case of the same PSNs, the node ids of the nodes that suggest the PSN is used to break a tie. For example, a message with a PSN of 3 suggested by "node1" is placed closer to the head of the queue than a message with a PSN of 3 suggested by "node3".
  3. Node A will pick the largest PSN among all the PSN received as the agreed sequence number (ASN) along with the node id of the Node, which suggests the largest PSN. If multiple nodes suggest the largest PSN, Node A will pick the largest node id (e.g., "node8" is larger than "node1"). After that, Node A multicasts the PSN and its node id to every Node in the system (including itself). Every Node will update the priority of this message in the hold-back queue to the ASN and update the node id suggesting the priority. It will re-multicast the ASN upon the first receipt, thus satisfying the reliable ordering.
  4. When each node receives and updates the ASN, node id for this message, they will mark the messages in the hold-back queue as deliverable. The node then checks the hold-back queue and delivers the deliverable message at the front of the queue (and removes it) repeatedly until the message at the head of the queue is undeliverable.
  5. If every node agrees on the set of the sequence number and delivers them in order, then the total ordering is satisfied.

* Failure Handling

  1. The node failure is detected with TCP errors in the design. 
  2. When a Node X detects Node A fails, it waits for about 10 seconds, which is two one-way delay times since we assume the maximum message delay between any two nodes is 4-5 seconds. After the timeout, Node X will delete the messages in the hold-back queue sent by node A that are still marked as undeliverable.

<div style="page-break-after: always"></div>
**Graph Measurement**

Graphs are generated in four different scenarios using Jupyter Notebook: small-scale, large-scale, small-scale with failure, and large-scale with failure.

**Process Time**

The processing time is the amount of time until a message is processed at all nodes. In this MP, we calculated the time difference between when a message is first delivered and when a message is finished, where the transaction is over, and the balance is printed.

* 3 nodes, 0.5 Hz each, running for 100 seconds
<img src = "stat/graph/smallScenario/3 Nodes Processing Time.png" alt="drawing" width="500">
* 3 nodes, 0.5 Hz each, runing for 100 seconds, then one node fails, and the rest continue to run for 100 seconds
<img src = "stat/graph/smallScenariowithFailure/3 Nodes Processing Time.png" alt="drawing" width="500">

* 8 nodes, 5 Hz each, running for 100 seconds
<img src = "stat/graph/largeScenario/8 Nodes Processing Time.png" alt="drawing" width="500">

* 8 nodes, 5 Hz each, running for 100 seconds, then 3 nodes fail simultaneously, and the rest continue to run for 100 seconds.
<img src = "stat/graph/largeScenariowithFailure/8 Nodes Processing Time.png" alt="drawing" width="500">
<div style="page-break-after: always"></div>

**Bandwidth**
The bandwidth is calculated by adding up the bytes received/sent at each second for each node. The graphs showing bandwidth under different scenarios are shown below. As shown in the figures, the bandwidth at the first few seconds is large, which is caused by the connection time. When creating the connection, the messages are stacked in the queue. When the connection is done, the messages will be delivered within a short time lead to a large bandwidth at the beginning.

* 3 nodes, 0.5 Hz each, running for 100 seconds
<img src = "stat/graph/smallScenario/node1Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/smallScenario/node2Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/smallScenario/node3Bandwidth.png" alt="drawing" width="500">

* 3 nodes, 0.5 Hz each, runing for 100 seconds, then node 1 fails, and the rest continue to run for 100 seconds
<img src = "stat/graph/smallScenariowithFailure/node1Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/smallScenariowithFailure/node2Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/smallScenariowithFailure/node3Bandwidth.png" alt="drawing" width="500">

* 8 nodes, 5 Hz each, running for 100 seconds
<img src = "stat/graph/largeScenario/node1Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenario/node2Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenario/node3Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenario/node4Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenario/node5Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenario/node6Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenario/node7Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenario/node8Bandwidth.png" alt="drawing" width="500">

* 8 nodes, 5 Hz each, running for 100 seconds, then node1-3 fail simultaneously, and the rest continue to run for 100 seconds.
<img src = "stat/graph/largeScenariowithFailure/node1Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenariowithFailure/node2Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenariowithFailure/node3Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenariowithFailure/node4Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenariowithFailure/node5Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenariowithFailure/node6Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenariowithFailure/node7Bandwidth.png" alt="drawing" width="500">
<img src = "stat/graph/largeScenariowithFailure/node8Bandwidth.png" alt="drawing" width="500">
<div style="page-break-after: always"></div>



