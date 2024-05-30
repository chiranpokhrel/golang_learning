# Lab 5: Multi-Paxos with Gorums and Performance Evaluation

| Lab 5: | Multi-Paxos with Gorums and Performance Evaluation |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **April 25, 2024 23:59** |
| Expected effort:         | 40-50 hours |
| Grading:                 | Pass/fail |
| Submission:              | Group |

## Table of Contents

1. [Introduction](#introduction)
2. [Gorums-based Multi-Paxos](#gorums-based-multi-paxos)
3. [Performance Evaluation](#performance-evaluation)

## Introduction

The overall objective of this lab is to implement Multi-Paxos using Gorums and do a performance evaluation of your implementation.
Thus this lab is divided in two parts:

## Gorums-based Multi-Paxos

Implement multi-decree Paxos algorithm using Gorums as described [here](gorumspaxos/gorums-paxos.md).
A TA will verify your implementation during lab hours.

## Performance Evaluation

You are required to conduct a performance evaluation of your implementation to see how it performs
under different settings.

You should deploy a set of replicas representing a service in which clients will send requests.
Each request is converted to a proposal by the replica that receives it, and if the replica is the leader,
it is proposed in the consensus protocol. If the replica is not the leader, the replica could forward
the message to the current leader.

You are free to propose different designs to deploy your service and how to handle the client requests.
Use the [paxosclient](gorumspaxos/cmd/paxosclient/main.go) as a starting point.
Clients should be able to query for the decided values and display them in the console.

You should focus on two main metrics:

- **Throughput:** How many requests per second can the service execute on
  average when under full load. Throughput should be measured at the replicas.

- **Latency:** What is the average round trip time (RTT) for client requests
  when the service is under full load? The clients should measure latency.

The client should send one request at a time and wait for a reply before sending a new one.
Clients should be automated to send a certain number of requests.
During the experiments, you should keep increasing the number of clients and requests to saturate your system and find the maximum throughput.

You should calculate statistics for both metrics, e.g., mean, standard deviation, median, minimum, maximum, 90th percentile.
You should perform experiments using different values for at least the following variables:

- Cluster size (e.g., three or five replicas)
- Number of clients
- Number of client requests

Results should be presented in a report using graphs and/or tables.
You should commit the report as a pdf file to your repository.
_You should be able to present and explain the results from the report during the lab approval_.
