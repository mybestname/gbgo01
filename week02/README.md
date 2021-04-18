#  第2周： 异常（错误）处理

## Error vs. Exception

- 对于真正意外的情况，那些表示不可恢复的程序错误，例如索引越界、不可恢复的环境问题、栈溢出，我们才使用 panic。对于其他的错误情况，我们应该是期望使用 error 来进行判定。
- You only need to check the error value if you care about the result.  -- Dave
- no that exceptions are bad. exceptions are too hard and not smart enough to handle them.

### Error  
- 简单
- 考虑失败，而不是成功（plan for failure, not success）
- 没有隐藏的控制流
- 完全交给你来控制 error
- [Error are values](https://blog.golang.org/errors-are-values) 

## Error Type

### Sentinel Error (避免使用哨兵错误)

- 即 `if err == ErrSomething` , Go使用特定的值来表示错误。
- 我们把预定义的，特定的某种错误叫做Sentinel Error。
- 哨兵错误是最不灵活的错误处理策略。
  - 因为调用方必须使用 `==` 将结果与预先声明的值进行比较。
  - 当想要提供更多的上下文时，因为返回一个不同的错误将破坏相等性检查。
  - 携带上下文的`fmt.Errorf()`，会破坏调用者的`==` ，调用者将被迫查看`error.Error()`方法的输出，以查看它是否与特定的字符串匹配。
- 不应该依赖`error.Error()`的输出，该结果是给人看的（用于log和stdout），而不是给程序处理的。
- Sentinel errors在API中的问题
  - Sentinel Error必须是公共的，然必要求文档，增加 API 的表面积。 
  - 返回Sentinel Error的接口，则该接口的所有实现都将被限制为仅返回该错误，即使它们可以提供更具描述性的错误。
    - `io.Reader`/`io.Copy` 的实现者比如返回 io.EOF 来告诉调用者没有更多数据了，但这又不是错误。
  - Sentinel errors在两包之间建立源码依赖。那么形成循环依赖的情况概率加大。
  
### Error Type (避免使用Error type)
- Error type 指实现了 `error` 接口的自定义类型
- 例子[`os.PathError`](https://golang.org/src/io/fs/fs.go?s=8967:9030#L233)

```golang
type PathError struct {
	Op   string
	Path string
	Err  error
}
func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }
```
#### 问题
- 调用者要使用类型断言和类型switch，就要让自定义的 error 变为 public。这种模型会导致和调用者产生强耦合，从而导致 API 变得脆弱。
- 虽然error type 比 sentinel errors 更好，但是 error types 共享 error values 许多相同的问题。
- 避免将它们作为公共API的一部分。

### Opaque Error （更好的方法）
- 不透明的错误处理
- 代码和调用者的耦合是最少的。
- 例子[net.Error](https://golang.org/pkg/net/#Error) ([src](https://golang.org/src/net/net.go?s=13516:13637#L387))
```golang
type Error interface {
	error
	Timeout() bool   // Is the error a timeout?
	Temporary() bool // Is the error temporary?
}
```
  - 不暴露数据（具体的error类型或者error值）
  - 而只暴露行为

## Handling Error

## Go 1.13 errors

## Go 2 Error Inspection

## References

- https://www.infoq.cn/news/2012/11/go-error-handle/
- https://golang.org/doc/faq#exceptions
- https://www.ardanlabs.com/blog/2014/10/error-handling-in-go-part-i.html
- https://www.ardanlabs.com/blog/2014/11/error-handling-in-go-part-ii.html
- https://www.ardanlabs.com/blog/2017/05/design-philosophy-on-logging.html
- https://medium.com/gett-engineering/error-handling-in-go-53b8a7112d04
- https://medium.com/gett-engineering/error-handling-in-go-1-13-5ee6d1e0a55c
- https://rauljordan.com/2020/07/06/why-go-error-handling-is-awesome.html
- https://morsmachine.dk/error-handling
- https://crawshaw.io/blog/xerrors
- https://dave.cheney.net/2012/01/18/why-go-gets-exceptions-right
- https://dave.cheney.net/2015/01/26/errors-and-exceptions-redux
- https://dave.cheney.net/2014/11/04/error-handling-vs-exceptions-redux
- https://dave.cheney.net/2014/12/24/inspecting-errors
- https://dave.cheney.net/2016/04/07/constant-errors
- https://dave.cheney.net/2019/01/27/eliminate-error-handling-by-eliminating-errors
- https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
- https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
- https://blog.golang.org/errors-are-values
- https://blog.golang.org/error-handling-and-go
- https://blog.golang.org/go1.13-errors
- https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html

