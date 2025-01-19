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

# Preventing Goroutine leaks

> 
> If a goroutine is responsible for creating a goroutine, it is also responsible for ensuring it can stop the goroutine.
> 

Goroutines are not garbage collected by the Go runtime. So how do we go about ensuring their cleaned up?

The goroutine has a few paths to termination:
1. When it has completed its work
2. When it cannot continue its work due to an unrecoverable error
3. When it's told to stop working

We get the first two paths for free - these paths are your algorithm - but what about work cancellation?
This turns out to be the most important bit because of the network effect: 
if you've begun a goroutine, it's most likely cooperating with several other goroutines in some sort of organized fashion.

The parent goroutine (often the main goroutine) with this full contextual knowledge should be able to tell its child goroutines to terminate.<br>
To successfully mitigate go routine leaks is to establish a signal between the parent goroutine and its children that allows the parent to signal cancellation to its children.<br>
By convention, this signal is usually a read only channel named done.
The parent goroutine passes this channel to the child goroutine and then closes to channel when it wants to cancel the child goroutine.

See leak_prevention.go for some examples.

## The or-channel

At times, you may find yourself wanting to combine one or more channels into a single channel that closes if any of its component channels close.<br>
Sometimes you cannot know the number of the channels you're working with at runtime. In this case, or if you just prefer a one-liner, you can combine these channels to get her using the or-channel pattern.<br>
This pattern creates a composite done channel through recursion and goroutines (see or_channel.go). This is a fairly concise function that enables you to combine any number of channels together into a single channel that will close as soon as any of its component channels are closed, or written to.

<b>When to use this pattern?</b><br>
This pattern is useful to employ at the intersection of modules in your system. At these intersections, you tend to have multiple conditions for canceling trees of courts through your cold stack.
Using the or function, you can simply combine these together and pass it down the stack.

We'll take a look at another way of doing this with the context package that is also very nice, and perhaps a bit more descriptive.

## Error handling

We should give our paths the same attention we give our algorithms. The most fundamental question when thinking about error handling is, "who should be responsible for handling the error?"

Separate your concerns: in general, your current processes should send their errors to another part of your program that is complete information about the state of your program, and can make a more informed decision about what to do (see error_handling.go).

The key thing to note here is how we have coupled the potential result with the potential error. This represents the complete set of possible outcomes created from the goroutine checkStatus, and allows our main goroutine to make decisions about what to do when errors occur.
In broader terms, we've successfully separated the concerns of error handling from our producer goroutine. This is desirable because the goroutine that spawned the producer goroutine - in this case our main goroutine - has more context about the running program, and can make more intelligent decisions about what to do with errors.

Because errors are returned from checkStatus and not handled internally within the goroutine, arrow handling follows to familiar Go pattern.

Again, the main takeaway here is that errors should be considered first-class citizens when constructing values to return from goroutines. If your goroutine can produce errors, those errors should be tightly, coupled with your result type, and pass along through the same line of communication - just like regular synchronous functions. 

## Pipelines

A pipeline is just under tool you can use to form and abstraction in your system. In particular, it is a very powerful tool to use when your program needs to process streams or batches of data.

In computer science, the word pipeline is used for something that transports data from one place to another. Pipeline is nothing more than a series of things to take, perform an operation on it, and pass the data back out.

We call each of these operations a stage of the pipeline.

By using a pipeline, to separate the concern of each stage, which provides numerous benefits:
1. You can modify stages independent of one another
2. You can mix and match how stages are combined independent of modifying the stages
3. You can process each stage concurrent to upstream or downstream stages
4. You can fan-out or
5. rate-limit portions of your pipeline

See pipeline.go for more information.

<b>What are the properties of a pipeline stage?</b>
1. A stage consumes and returns to same type
2. A stage must be reified by the language so that it may be passed around. Functions and go are reified and fit this purpose nicely.

reified ... within the context of languages, reification means that the language exposes a concept to the developers so that they can work with it directly. Functions and goal are set to be reified because you can define variables that have a type of a function signature. This also means you can pass functions around your program.

batch processing ... operating on chunks of data all at once (instead of one discrete value at a time)
stream processing ... the stage receives and emits one element at a time

### Best practice for constructing pipelines

Channels are uniquely, suited to constructing pipelines and go because they fulfill all of our basic requirements.
1. They can receive and emit values
2. They can safely be used concurrently
3. They can be ranged over
4. They are reified by the language

In bestpractice_pipeline.go, the examples from pipeline.go have been converted to work with channels.