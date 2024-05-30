# Lab 4: Single-decree Paxos and Multi-Paxos

| Lab 4: | Single-decree Paxos and Multi-Paxos |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **March 14, 2024 23:59** |
| Expected effort:         | 30-40 hours |
| Grading:                 | Pass/fail |
| Submission:              | Group |

## Table of Contents

1. [Introduction](#introduction)
2. [Resources](#resources)

## Introduction

The overall objective of this lab is to implement a single-decree and a multi-decree version of Paxos (also known as Multi-Paxos).
The assignment consist of three parts:

1. A set of theory [questions that you should answer](answers.md).

2. Implementation of the single-decree Paxos algorithm as described [here](singlepaxos/singlepaxos.md).
   This variant of Paxos is only expected to choose a single command.
   It is intended as an exercise and to help you understand the core logic of Paxos before moving on to the more complex Multi-Paxos variant.

   This assignment involve implementing each of the three Paxos roles: Proposer, Acceptor and Learner.
   This implementation has corresponding unit tests and will be automatically verified by QuickFeed.

3. Implementation of the multi-decree Paxos algorithm for each of the three Paxos roles as described [here](multipaxos/multipaxos.md).
   This variant of Paxos is expected to choose multiple commands, and is the version of Paxos that is used in practice.
   This implementation has corresponding unit tests and will be automatically verified by QuickFeed.

## Resources

Several Paxos resources are listed below.
You should use these resources to answer the [questions](answers.md) for this lab.
You are also advised to use them as support literature when working on your implementation now and in future lab assignments.

- [Paxos Explained from Scratch](resources/paxos-scratch-slides.pdf) - slides.
- [Paxos Explained from Scratch](resources/paxos-scratch-paper.pdf) - paper.
- [Paxos Made Insanely Simple](resources/paxos-insanely-simple.pdf) - slides.
  Also contains pseudo code for the Proposer and Acceptor.
- [Paxos Made Simple](resources/paxos-simple.pdf)
- [Paxos Made Moderately Complex](resources/paxos-made-moderately-complex.pdf)
- [Paxos Made Moderately Complex (ACM Computing Surveys)](resources/a42-renesse.pdf)
- [Paxos for System Builders](resources/paxos-system-builders.pdf)
- [The Part-time Parliament](resources/part-time-parliment.pdf)
