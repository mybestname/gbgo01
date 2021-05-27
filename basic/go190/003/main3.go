package main

func m2() {
	s1 := make([]int,8191)                     // MaxImplicitStackVarSize = int64(64*1024) /8 = 8192
	for i := range s1 {
		s1[i] = i+1
	}
	//println(len(s1),cap(s1), s1[8190])

	s2 := []int{1,2,3,4,5}
	//println(len(s2),cap(s2), s2[4] )

	var a  = [...]int{0:100, 2:7, 1310719:-1}  // MaxStackVarSize = int64(10 * 1024 * 1024) / 8 = 1310720
	//println(len(a),cap(a), a[0], a[2], a[1310719])
	s3 := a[:]
	for i := range s3 {
		s3[i] = i+1
	}
	//println(len(s3),cap(s3), s3[0], s3[2], s3[1310719])
	_ = s1
	_ = s2
	_ = s3

}
func main() {
	m2();
	// output :
	// 8191 8191 8191
	// 1310720 1310720 100 7 -1
	// 1310720 1310720 1 3 1310720
}

//go:generate go tool compile -W -N -l -m=2 main3.go
//go:generate go tool objdump main3.o
//go:generate rm main3.o
//GOSSAFUNC=m2 go build -gcflags "-N -l" main3.go

// 对于go来说，如果想初始化slice，同时指定默认值，只能靠slice literal
// 因为从make的实现来说，参数只是指定backend的初始数组长度（cap）以及slice的长度。
// 默认填充的是零值。
// 而ssa阶段的OpSliceMake操作，输入为：slice的元素类型、backend的对应数组的指针、slice大小，slice容量。也不可能填充默认值。
//
// 那么literal方式的slice是唯一可以指定slice默认值的。
// 或者使用literal方式的数组初始化，并获得一个指向该数组的slice。
//
// literal模式的slice的初始化
//
// 1. 语法扫描和分析器三架马车 (cmd/compile/internal/syntax包中的：token，scanner和parser）
// - https://github.com/golang/go/blob/master/src/cmd/compile/internal/syntax/tokens.go
// - https://github.com/golang/go/blob/master/src/cmd/compile/internal/syntax/scanner.go
//   - https://github.com/golang/go/blob/master/src/cmd/compile/internal/syntax/source.go
// - https://github.com/golang/go/blob/master/src/cmd/compile/internal/syntax/parser.go
//
// 1.0 语法扫描和分析的入口点： (gc -> noder -> parser -> scanner)
//  - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/gc/main.go#L192  LoadPackage( filenames )
//  - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/noder/noder.go#L41-L64
//    - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/syntax/syntax.go#L68 syntax.Parse
//
// 1.1 scanner一个token一个token的读取原文件
// https://github.com/golang/go/blob/master/src/cmd/compile/internal/syntax/scanner.go#L88
//
// ```
//  func (s *scanner) next() {
//  	...
//  redo:
//  	// skip white space
//  	s.stop()
//  	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
//  		s.nextch()
//  	}
//
//  	// token start
//  	...
//  	switch s.ch {
//  	...
//  	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
//  		s.number(false)
//  	...
//  	case '[':
//  		s.nextch()
//  		s.tok = _Lbrack
//  	case '{':
//  		s.nextch()
//  		s.tok = _Lbrace
//  	...
//  	case ')':
//  		s.nextch()
//  		s.nlsemi = true
//  		s.tok = _Rparen
//  	...
//  	case '=':
//  		s.nextch()
//  		if s.ch == '=' {
//  			s.nextch()
//  			s.op, s.prec = Eql, precCmp
//  			s.tok = _Operator
//  			break
//  		}
//  		s.tok = _Assign
//  	...
//  	default:
//  		s.errorf("invalid character %#U", s.ch)
//  		s.nextch()
//  		goto redo
//  	}
//  	return
//  	...
//  }
// ```
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/syntax/source.go#L113
// ```
// func (s *source) nextch() {
//  	...
//  	// fast common case: at least one ASCII character
//  	if s.ch = rune(s.buf[s.r]); s.ch < sentinel {
//  	...
// ```
// 每一个token 存入 source这个结构体的buffer中。
// ```
//  type source struct {
//  	in   io.Reader
//  	errh func(line, col uint, msg string)
//
//  	buf       []byte // source buffer
//  	ioerr     error  // pending I/O error, or nil
//  	b, r, e   int    // buffer indices (see comment above)
//  	line, col uint   // source position of ch (0-based)
//  	ch        rune   // most recently read character
//  	chw       int    // width of ch
//  }
// ```
// 而 scanner 包含这个source 结构体
// ```
//  type scanner struct {
//  	source
//  	mode   uint
//  	nlsemi bool // if set '\n' and EOF translate to ';'
//
//  	// current token, valid after calling next()
//  	line, col uint
//  	blank     bool // line is blank up to col
//  	tok       token
//  	lit       string   // valid if tok is _Name, _Literal, or _Semi ("semicolon", "newline", or "EOF"); may be malformed if bad is true
//  	bad       bool     // valid if tok is _Literal, true if a syntax error occurred, lit may be malformed
//  	kind      LitKind  // valid if tok is _Literal
//  	op        Operator // valid if tok is _Operator, _AssignOp, or _IncOp
//  	prec      int      // valid if tok is _Operator, _AssignOp, or _IncOp
//  }
// ```
// 1.2 parser根据语法规则去处理token，生成AST
// 而 parser的结构体又包含了scanner
// ```
//  type parser struct {
//  	file  *PosBase
//  	errh  ErrorHandler
//  	mode  Mode
//  	pragh PragmaHandler
//  	scanner              //包含扫描后的token化的source
//
//  	base   *PosBase // current position base
//  	first  error    // first error encountered
//  	errcnt int      // number of errors encountered
//  	pragma Pragma   // pragmas
//
//  	fnest  int    // function nesting level (for error handling)
//  	xnest  int    // expression nesting level (for complit ambiguity resolution)
//  	indent []byte // tracing support
//  }
// ```
//
// 1.3 noder 生成 syntax.File 结构（AST），然后转化为 IR Node 树。
//
// 1.3.1 syntax.File的定义
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/syntax/nodes.go#L32-L111
//
// 1.3.2 生成 *syntax.File 结构
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/syntax/parser.go#L376 fileOrNil()
//   - 注意：这里是和scanner的关联点，因为parser结构体包含了scanner结构体，所以当parser调用next()时候，调用了scanner.next()获取下一个token
//     - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/syntax/parser.go#L405-L423
// ```
//   for p.tok != _EOF {
//  		switch p.tok {
//  		case _Const:
//  			p.next()   // 这个调用的是 scanner，来获得下一个token
//  			f.DeclList = p.appendGroup(f.DeclList, p.constDecl)
//  		....
// ```
// 1.3.3 ir节点树的结构的定义
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/ir/package.go#L10
//
// 1.3.4 noder 把AST变成一个节点树。
// https://github.com/golang/go/blob/master/src/cmd/compile/internal/noder/noder.go#L78-L82
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/noder/irgen.go#L23 check2
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/noder/irgen.go#L99 generate
//
// 例如：decls(decls []syntax.Decl) []ir.Node，输入是[]sntax.Decl 结构，输出变为ir.Node
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/noder/decl.go#L21
// ```
//
//
//
// 1.4 Ir节点的定义
// https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/ir/node.go#L20
//
//
//
//
//
//
//
//
//
// ==================================================================================================
// 附注：1.17(trunk)的重大变化 type2/noder2/irgren/
// ==================================================================================================
// 关于范型提案
//   - https://github.com/golang/go/issues/43651 (lanlancetaylor  Jan 13)
//      spec: add generic programming using type parameters
//   - https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md （March 19, 2021）
//     Type Parameters Proposal
//
// types和types2的区别 (1.17 trunk 目的 go的generic支持)
//   - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/types/type.go
//   - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/types2/type.go
//
// 关于noder2 :
//  parse a simple generic function and print noder IR via 'go tool compile -G=2 -W=2 func.go'
//
// noder2中的types/types2的桥
//   - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/noder/types.go
//
// 改使用irgen，
// - https://github.com/golang/go/commit/ef5285fbd0636965d916c81dbf87834731f337b2  mdempsky Jan-14
//    [dev.typeparams] cmd/compile: add types2-based noder
//    This CL adds "irgen", a new noding implementation that utilizes types2
//    to guide IR construction. Notably, it completely skips dealing with
//    constant and type expressions (aside from using ir.TypeNode to
//    interoperate with the types1 typechecker), because types2 already
//    handled those. It also omits any syntax checking, trusting that types2
//    already rejected any errors.
//
//    It currently still utilizes the types1 typechecker for the desugaring
//    operations it handles (e.g., turning OAS2 into OAS2FUNC/etc, inserting
//    implicit conversions, rewriting f(g()) functions, and so on). However,
//    the IR is constructed in a fully incremental fashion, so it should be
//    easy to now piecemeal replace those dependencies as needed.
//
//    Nearly all of "go test std cmd" passes with -G=3 enabled by
//    default. The main remaining blocker is the number of test/run.go
//    failures. There also appear to be cases where types2 does not provide
//    us with position information. These will be iterated upon.
//
// 如何在编译器打开范型 : (1.17 trunk)
//    run -gcflags=-G=3
//    例子：
//    - https://github.com/golang/go/blob/193d5141318d65cea310d995258288bd000d734c/test/typeparam/graph.go#L14-L30
//
// 关于-G参数的含义：
//   1. https://github.com/golang/go/issues/43931  (mdempsky Jan 27)
//      all: merge dev.typeparams to master during Go 1.17
//
//      mdempsky:(Feb 23)
//       - Currently you still need to use -G=3. We haven't yet renumbered the flag values like I suggested
//         in the proposal.
//       - Beware that the generics-specific support is still very early. Non-generic Go code is expected
//         to fully work, so please file issues if you find anything that works without -G=3 but breaks
//         with -G=3. But don't be surprised if code using new generics features doesn't work yet.
//
//   2. https://go-review.googlesource.com/c/go/+/295029 (mdempsky Feb23)
//      cmd/compile: renumber -G flag values
//
//       The current -G values were incrementally added as generics support was
//       plumbed further through the compiler to avoid breaking existing tests
//       that only worked yet with the earlier stages. Now that we have
//       reasonably robust end-to-end support for using types2, we can simplify
//       the flags somewhat to be more useful to end users.
//
//       In particular, this CL renumbers the -G values as proposed in #43931:
//
//         -G=0: continue using legacy typechecker
//         -G=1: use types2, but without support for generics
//         -G=2: use types2, with support for generics
//
//       The new -G=2 flag is equivalent to what was formerly -G=3, and
//       hopefully the new -G=1 flag can be enabled by default for Go 1.17
//       after we finish reconciling types2's error messages with the existing
//       regress tests.
//
//       While here, also add a -d=typecheckonly debug flag that can be used to
//       have the compiler gracefully exit after typechecking, so that we can
//       help distinguish frontend vs backend errors. This also happens to be
//       needed for smoketest.go, which typechecks at the moment, but does not
//       yet compile.
//
//   3. https://github.com/golang/go/issues/45597#issuecomment-822864432
//      griesemer:(Apr 20)
//         -G=3 won't be enabled for 1.17
//
// 问题：目前还不支持ssa分析输出：
//    1. 默认无法对范型函数使用ssa：
//    2. 如果把_SliceEqual改为SliceEqual，则无法编译。
//    ./graph.go:16:6: internal compiler error: Cannot export a generic function (yet): SliceEqual
//    3. 改为sliceEqual 不报错但无效
//      GOSSAFUNC=sliceEqual ../../bin/go build -gcflags "-G=3" graph.go
//    4. 改变G level 也会报错，只能用-G=3
//     $ GOSSAFUNC=sliceEqual ../../bin/go build -gcflags "-G=2" graph.go
//     go build command-line-arguments: open /var/folders/n7/mfyh01tx6rbd7cfvrqcn2b1h0000gn/T/go-build444510394/b001/_pkg_.a: no such file or directory
//
// 总结：目测1.17 release的计划应该是：
//    - 默认支持types2（-G=1），
//    - 然后支持范型（-G=2），
//    - -G=3供开发使用。
//    - 目前trunk上的代码还是正常用-G=2
