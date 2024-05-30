# Lab 2: Network Programming in Go

| Lab 2: | Network Programming in Go |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **February 8, 2024 23:59** |
| Expected effort:         | 20-25 hours |
| Grading:                 | Pass/fail |
| Submission:              | Individually |

## Table of Contents

1. [Introduction](#introduction)
2. [Lab 2 Assignments](#lab-2-assignments)

## Introduction

The goal of this lab assignment is to get you started with network programming in Go.
The overall aim of the lab project is to implement a fault tolerant distributed application.
Knowledge of network programming in Go is naturally a prerequisite for accomplishing this.

This lab assignment consist of four parts.
In the first part you are expected to implement a simple echo server that is able to respond to different commands specified as text.
In the second part, you will be implementing an in-memory key-value storage using the gRPC framework.
In the third part, you will be implementing a fault tolerant storage server using the Gorums framework.
Finally, the last part of this lab consists of answering some simple network programming related questions.

The most important package in the Go standard library that you will use in the first part of this assignment is the [`net`](http://go.dev/pkg/net) package.
It is recommended that you actively use the documentation available for this package during your work on this lab assignment.

## Lab 2 Assignments

- [UDP Echo Server](uecho/echo-server.md)
- [GRPC Server](grpc/grpc.md)
- [Gorums](gorums/gorums.md)
- [Poll](poll/poll.md)
- [Networking Questions](networking/network_questions.md)
