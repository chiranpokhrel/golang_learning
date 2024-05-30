# Lab 1: Getting Started

| Lab 1: | Getting Started |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **January 25, 2024 23:59** |
| Expected effort:         | 10-15 hours |
| Grading:                 | Pass/fail |
| Submission:              | Individually |

## Table of Contents

1. [Introduction](#introduction)
2. [Lab Overview](#lab-overview)
3. [Getting Started](#getting-started)

## Introduction

This first lab will give you an overview of the lab project as a whole.

## Lab Overview

The lab project for this course consist of six assignments (the list may change
during the semester):

1. Getting Started
2. Network Programming in Go
3. Failure Detector and Leader Election
4. Single-Decree Paxos and Multi-Paxos
5. Multi-Paxos with Gorums and Performance Evaluation
6. Bank Application with Reconfiguration

The first two assignments are introductory and individual.
They are intended to setup the working environment for the lab submissions and get you up to speed on the Go programming language and basic network programming.
The remaining assignments form the main lab project and should be done in groups.

The aim of the overall lab project is to implement a fault tolerant application, replicated on several machines, using the Paxos protocol.
The Paxos protocol will be introduced in the course accordingly.
For this we will use the Go programming language.
Go is very well suited to implement the event-driven structure of the Paxos protocol.
You will for each lab construct an independent part of the application, resulting in a complete implementation in Lab 6.

## Getting Started

Enough said, let's get started!
The first lab is composed by two parts:

1. [Setup](setup/README.md) - Setting up the working environment and tools required by all labs.
2. [Go](gointro/README.md) - Basic programming tasks to get you familiar with the Go programming language.
