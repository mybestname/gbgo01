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

