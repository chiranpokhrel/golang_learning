# Lab 1: Getting Started

## Table of Contents

1. [Introduction](#introduction)
2. [Go Resources](#go-resources)
3. [Primary Resources](#primary-resources)
4. [Additional Resources](#additional-resources)
5. [Tasks](#tasks)
6. [Task 1: Exercises from Tour of Go](#task-1-exercises-from-tour-of-go)
7. [Task 2: Go Language Questions](#task-2-go-language-questions)
8. [Task 3: Go Exercises](#task-3-go-exercises)

## Introduction

In this part of the lab you will be learning the Go programming language, which we will be using throughout the course.
Go is a very nice language in many ways.
It provides built in primitives to design concurrent programs using lightweight threads, called goroutines.
In addition, these goroutines can communicate using channels, instead of by means of shared memory which is the most common approach in other languages such as C and Java.
Usage of goroutines and channels will be covered in depth in later labs.
For now you just need to know that they exist.
Go is also a very small language in the sense that it has very few keywords, and as such it does not take too long to learn the fundamentals of the language.

## Go Resources

Below we provide some resources which will be useful when you start programming in Go.

### Primary Resources

[The Go Language web page](https://go.dev) provides a variety of helpful documentation, including:

* Learning the basics with ["A Tour of Go"](https://go.dev/tour/)
* [The Go documentation](https://go.dev/doc/)
* [Standard library package documentation](https://pkg.go.dev/std)

### Additional Resources

* [Effective Go](https://go.dev/doc/effective_go) gives tips for writing clear, idiomatic Go code
* [Frequently Asked Questions](https://go.dev/doc/faq)
* [The Go Blog](https://go.dev/blog/all) contains many great articles that describe important idioms and intricacies of the language.
* [Collection of Videos about Go](https://go.dev/wiki/GoTalks)
* [The Go Programming Language (book)](http://www.gopl.io)
* [Golang Tutorial Series](https://golangbot.com/learn-golang-series/)

**Troubleshooting tip:**
When searching for information about Go, use the term *golang* instead of go.

## Tasks

### Task 1: Exercises from Tour of Go

Start learning the basics of Go by completing ["A Tour of Go"](http://go.dev/tour/).
You should do at least the following exercises.

* [Exercise: Loops and Functions](https://go.dev/tour/flowcontrol/8)
* [Exercise: Slices](https://go.dev/tour/moretypes/18)
* [Exercise: Maps](https://go.dev/tour/moretypes/23)
* [Exercise: Errors](https://go.dev/tour/methods/20)
* [Exercise: rot13Reader](https://go.dev/tour/methods/23)

Note that you can change the code inline in the browser and run the code to see the results.

### Task 2: Go Language Questions

Answer these multiple choice questions about [Go programming](go_questions.md).

### Task 3: Go Exercises

Before you start working on the assignments below, make sure that your local working copy has all the latest changes from the course [assignments](https://github.com/dat520-2024/assignments) repository.
Instructions for fetching the latest changes are [here](https://github.com/dat520-2024/info/blob/main/lab-submission.md#update-local-working-copy-from-course-assignments).

1. In the following, we will use `sequence/fibonacci.go` exercise as an example.
   The file contains the following skeleton code and task description:

    ```golang
    package sequence

    // Task: Fibonacci numbers
    //
    // fibonacci(n) returns nth Fibonacci number, and is defined by the
    // recurrence relation F_n = F_n-1 + F_n-2, with seed values F_0=0 and F_1=1.
    func fibonacci(n uint) uint {
    0
    }
    ```

2. Implement the function body according to the specification so that all the tests in `sequence/fibonacci_test.go` passes.
   The test file looks like this:

    ```golang
    package sequence

    import "testing"

    var fibonacciTests = []struct {
        in, want uint
    }{
        {0, 0},
        {1, 1},
        {2, 1},
        {3, 2},
        {4, 3},
        {5, 5},
        {6, 8},
        {7, 13},
        {8, 21},
        {9, 34},
        {10, 55},
        {20, 6765},
    }

    func TestFibonacci(t *testing.T) {
        for _, ft := range fibonacciTests {
            got := fibonacci(ft.in)
            if got != ft.want {
                t.Errorf("fibonacci(%q) = %q, want %q", ft.in, got, ft.want)
            }
        }
    }
    ```

3. There are several ways to run the tests. If you run:

   ```console
   go test
   ```

   the Go tool will run all tests found in files whose file name ends with `_test.go` (in the current directory).
   Similarly, you can also run a specific test as follows:

   ```console
   go test -run TestFibonacci
   ```

4. You should ***not*** edit files or code that are marked with a `// DO NOT EDIT` comment.
   Please make separate `filename_test.go` files if you wish to write and run your own tests.

5. When you have completed a task and sufficiently many local tests pass, you may push your code to GitHub.
   This will trigger QuickFeed which will then run a separate test suite on your code.

   Using `sequence/fibonacci.go` as an example, use the following procedure to commit and push your changes to GitHub and QuickFeed:

    ```console
    $ git add fibonacci.go
    $ git commit
    // This will open an editor for you to write a commit message
    // Use for example "Implemented Assignment Fibonacci"
    $ git push
    ```

6. QuickFeed will now build and run a test suite on the code you submitted.
   You can check the output by going to the [QuickFeed web interface](https://uis.itest.run).
   The results (build log) is available from the Labs menu.
   Note that the results shows output for all the tests in current lab assignment.
   You will want to focus on the output for the specific test results related to the task you're working on.

7. Follow the same process for the other tasks included in this lab assignment.
   Each task contains a single `.go` template file, along with a task description and a `_test.go` file with tests.

8. When you are done with all assignments and want to submit the final version, please follow these [instructions](https://github.com/dat520-2024/info/blob/main/lab-submission.md#final-submission-of-labx).
