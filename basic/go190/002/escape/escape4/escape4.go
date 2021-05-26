package escape4

// original code from with modification
// https://www.ardanlabs.com/blog/2017/06/design-philosophy-on-data-and-semantics.html
// https://github.com/ardanlabs/gotraining/tree/master/topics/go/profiling
// Design Philosophy On Data And Semantics
//
// > “Value semantics keep values on the stack, which reduces pressure on the Garbage Collector (GC).
// > However, value semantics require various copies of any given value to be stored, tracked and
// > maintained. Pointer semantics place values on the heap, which can put pressure on the GC.
// > However, pointer semantics are efficient because only one value needs to be stored, tracked
// > and maintained.”  - Bill Kennedy (Go in Action 作者 & ardanlabs go培训）
// >
// > "If the logs are not working for you during dev, they are certainly not going to work for bug."
// >  -- William Kennedy  (Go in Action 作者 & ardanlabs go培训）
// >
// > “The hardest bugs are those where your mental model of the situation is just wrong, so you
// > can’t see the problem at all” - Brian Kernighan
// >
// > “C is the best balance I’ve ever seen between power and expressiveness. You can do almost
// >  anything you want to do by programming fairly straightforwardly and you will have a very
// >  good mental model of what’s going to happen on the machine; you can predict reasonably
// >  well how quickly it’s going to run, you understand what’s going on ….” - Brian Kernighan
// >
// > “If you don’t understand the data, you don’t understand the problem. This is because all problems
// >  are unique and specific to the data you are working with. When the data is changing, your
// >  problems are changing. When your problems are changing, the algorithms (data transformations)
// >  needs to change with it.” - Bill Kennedy
// >
// > "Integrity means that every allocation, every read of memory and every write of memory is
// >  accurate, consistent and efficient. The type system is critical to making sure we have
// >  this micro level of integrity.” - William Kennedy
// >
// > “Methods are valid when it is practical or reasonable for a piece of data to have a capability.”
// >  - William Kennedy
// >
// > “Polymorphism means that you write a certain program and it behaves differently depending on the
// >  data that it operates on.” - Tom Kurtz (inventor of BASIC)
// >
//
// - semantics/Mental model -> 使用值语义还是指针语义需要仔细斟酌 语义和心智模型的重要性（对机器<->心智模型的真正理解）
// - Debugging -> 有用的日志才有用，log什么？
// - Data Oriented Design 而非ODD -> 还是对机器的理解。机器<->心智
// - 类型的关键性（type is life） -> 值语义还是指针语义? 机器<->心智
// - The purpose of methods is to give a piece of data capability.
//   - should I use a value receiver or pointer receiver? Once I hear this question, I know that the developer
//     doesn’t have a good grasp of these semantics.
//   - receiver type for a method 需要了然于胸。
// -
// usage of godoc analysis
// 1. install
//   $ go install golang.org/x/tools/cmd/godoc@latest
// 2. run (must in GOPATH mode)
//   $ GOPATH=$HOME/go GO111MODULE=off godoc -analysis="pointer,type" -http=:8080
// 3. need to copy/link source to GOPATH/src to do the analysis
// 4. more on https://golang.org/lib/godoc/analysis/help.html
//
