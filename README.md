# Concurrency in Go

this repository contains concurrency patterns from Katherine Cox-Buday's book "Concurrency in Go". Below you find a detailed explanation of each pattern directly from the book.

NOTE: For the sake of brevity, some explanations were shortened, so there is no guarantee for correct citation. However, I tried to leave out all prose and focus on the relevant details and opinion of Katherine Cox-Buday. 

## Confinement

Confinement is the simple yet powerful idea of ensuring information is only ever available from one concurrent process. When this is achieved, a concurrent program is implicitly safe and no synchronization is needed.
There are two kinds of confinement possible: ad-hoc and lexical.

### Ad-hoc confinement
Ad hoc confinement is when you achieve confinement through a convention – whether it be set by the language is community, the group you work within, or the code base you work within.
Sticking to convention is difficult to achieve in projects of any size.

As the code is touched by many people, and deadlines loom, mistakes might be made, and the confinement and breakdown and cause issues. This is why lexical confinement should be preferred: it wields the compiler to enforce confinement.

### Lexical confinement
Lexical confinement involves using lexical scope to expose only the correct data concurrency primitive for multiple concurrent processes to use. It makes it impossible to do the wrong thing.

You can find a code snippet that enforces lexical confinement in confinement.go.

Why pursue confinement we have synchronization available to us? The answer is improved performance and reduced cognitive load on developers.
Concurrent code that utilizes lexical confinement also has the benefit of usually being simpler to understand and concurrent code without lexically confined variables. This is because within the context of your lexical scope you can write synchronous code.

## The for-select loop

Something you see over and over in Go programs is the for-select loop. It's nothing more than something like this:

```go
for { // either loop infinitely or range over something
    select {
    // do some work with channels
    }
}
```

There are several different variations of the for-select loop such as:
1. Sending iteration variables out on a channel
2. Looping infinitely waiting to be stopped (two variations)

In forselect.go, you find examples of the different variations.