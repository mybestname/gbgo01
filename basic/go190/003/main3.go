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
//   - https://github.com/golang/go/blob/master/src/go/parser/parser.go
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
// ==================================================================================================
// 附注1：关于两个parser(go/parser和compile/internal/xxx/parser.go)的问题
// ==================================================================================================
//
// 1. go/parser.go (这个是parser包，是open的api，可以供外界（go语言使用者）调用。例如fmt）
//    另外gc最早是用C实现的，这个包为gc用go重写（1.3~1.5)也提供了帮助 (这个包提供了native go lexer and parser)
// https://github.com/golang/go/commit/835cd46941683bce57689a184ab934b7739da036 (go0)
//  - https://github.com/golang/go/blob/835cd46941683bce57689a184ab934b7739da036/usr/gri/src/parser.go
//  -  first cut a Go parser in Go (gri,Jul 8, 2008)
// https://github.com/golang/go/commit/07513c259906359e9e5d6e2b3d05fd75f7e6f129 (go0)
//  - https://github.com/golang/go/blob/07513c259906359e9e5d6e2b3d05fd75f7e6f129/src/lib/go/parser.go
//  - move to src/lib/go folder (gri, Apr 1, 2009 )
// https://github.com/golang/go/commit/d90e7cbac65c5792ce312ee82fbe03a5dfc98c6f (go0)
//  - https://github.com/golang/go/blob/d90e7cbac65c5792ce312ee82fbe03a5dfc98c6f/src/pkg/go/parser/parser.go
//  - mv src/lib to src/pkg (r,Jun 10, 2009)
// https://github.com/golang/go/commit/c007ce824d9a4fccb148f9204e04c23ed2984b71 (go1.4beta1)
//  - https://github.com/golang/go/blob/c007ce824d9a4fccb148f9204e04c23ed2984b71/src/go/parser/parser.go
//  - move form src/pkg/go to src/go (rsc, Sep 8, 2014)
//
// 2. cmd/compile/internal/syntax/parser.go (这个是go的编译编译器真正在用的parser，供编译器内部使用）
// https://github.com/golang/go/commit/b569b87ca1f6905ed38cafc2bae6e26f4ec416b7 (go1.6beta1)
//  - https://github.com/golang/go/blob/b569b87ca1f6905ed38cafc2bae6e26f4ec416b7/src/cmd/compile/internal/gc/parser.go
//  - 初始创建cmd/compile/internal/gc/parser.go
//  - cmd/compile/internal/gc: recursive-descent parser (gri, Nov4,2015)
//  - 取递归下降法编译器，代原来的yacc编译器。
//       a translation of the yacc-based parser with adjustements
//       to make the grammar work for a recursive-descent parser followed
//       by cleanups and simplifications. The new parser is enabled by default.
// https://github.com/golang/go/commit/c8683ff7977c526fb48ae007971fed16ef32ff62 (go1.8beta1)
//  - https://github.com/golang/go/blob/c8683ff7977c526fb48ae007971fed16ef32ff62/src/cmd/compile/internal/syntax/parser.go
//  - 移动 cmd/compile/internal/gc/parser.go ->  cmd/compile/internal/syntax/parser.go
//  - cmd/compile/internal/syntax: fast Go syntax trees, initial commit. (gri, Aug 19, 2016)
//       Syntax tree nodes, scanner, parser, basic printers.
//       Builds syntax trees for entire Go std lib at a rate of ~1.8M lines/s
//       in warmed up state (MacMini, 2.3 GHz Intel Core i7, 8GB RAM)
//
// 3. 中间过程的y.go (yacc编译器）
// https://github.com/golang/go/commit/8c195bdf120ea86ddda42df89df7bfba80afdf10 (go1.5beta1)
//   - https://github.com/golang/go/blob/8c195bdf120ea86ddda42df89df7bfba80afdf10/src/cmd/internal/gc/y.go
//   - c代码转换的yacc go版编译器
//   - First draft of converted Go compiler, using rsc.io/c2go rev 83d795a. (rsc, Feb 18, 2015)
//
// https://github.com/golang/go/commit/17eba6e6b72b9dbf24d73a84be22edd65c229631 (go1.5beta1)
//   - https://github.com/golang/go/blob/17eba6e6b72b9dbf24d73a84be22edd65c229631/src/cmd/compile/internal/gc/y.go
//   - src/cmd/internal/gc/y.go -> src/cmd/compile/internal/gc/y.go
//   - cmd/compile, cmd/link: create from 5g, 5l, etc (rsc, May 22, 2015)
//
//
// ==================================================================================================
// 附注：关于gccgo，以及gccgo和gc的区别
// ==================================================================================================
//
// - The Gccgo compiler is a front end written in C++ with a recursive descent parser coupled to the standard GCC back end.
// - gccgo 时间线:
//    - 2007~ 2008: Robert Griesemer, Rob Pike and Ken Thompson laid out the goals and original specification of the language.
//       - gri, rob dnd ken sketching the goals on the white board on September 21, 2007.
//       - continued part-time in parallel with unrelated activities
//       - ken started work on a compiler with which to explore ideas. January 2008,
//          - it generated C code as its output.
//       - a full-time project and had settled enough to attempt a production compiler (mid-2008)
//    - 2008-mid: Ian Taylor read the draft specification and decided to write `gccgo` (an independent GCC front end.) (Meanwhile, mid-2008)
//       - b390608636231c2f443b55000db27580ae386727 (gri, added gccgo makefile Sep 23, 2008）
//          - +GO = /home/iant/go/bin/gccgo  注意显然：iant一直在做gcc方面的工作。
//       - a second compiler (gccgo) was not only encouraging, it was enabling. Having a second implementation
//         of the language was vital to the process of locking down the specification and libraries, helping
//         guarantee the high portability that is part of Go's promise. (robpike, go-ten-years-and-climbing sep 2017)
//    - 2008-late: Russ Cox joined later and helped move the language and libraries from prototype to reality. (late 2008)
//        - The “gc” Go toolchain is derived from the Plan 9 compiler toolchain.
//            - The assemblers, C compilers, and linkers are adopted essentially unchanged,
//            - uses a variant of the Plan 9 loader to generate ELF binaries.
//            - the `gc` compilers (in cmd/gc, cmd/5g, cmd/6g, and cmd/8g) are new C programs that fit into the toolchain.
//            - `gc` is written in C using `yacc` and `bison` for the parser.
//        - Why is the compiler called 6g
//          - (6g/8g/5g) compiler is named in the tradition of the Plan 9 C compilers
//          - 6 -> the architecture letter for amd64 (or x86-64, if you prefer). g stands for go.
//          - 8 -> 386 (8g,8l,8c,8a)
//          - 5 -> arm (5g,5l,5c,5a)
//        - Why C ?
//           - considered writing `6g` in Go itself but elected not to do so because of the difficulties of bootstrapping
//             and setting up a Go environment.
//           - also considered using LLVM for `6g`, but we felt it was too large and slow to meet our performance goals.
//    - 2009~2010 gccgo的大致样子
//        - https://github.com/golang/go/blob/2b1dbe8a4f9f214d7164abd18d99e2451efc5cdb/doc/go_lang_faq.html (rob Oct 6, 2009)
//        - https://github.com/golang/go/blob/3227445b75044cb94649e55b75e12656ca93f2d6/doc/go_faq.html (rsc Oct 22, 2009)
//        - https://github.com/golang/go/blob/8c40900fc28a7dda33a491902c5545b9abb37f58/doc/go_gccgo_setup.html (itaylor, Nov 7, 2009)
//        - https://github.com/golang/go/blob/8653acb191473126baa409d929841fe9fed3c734/doc/gccgo_contribute.html (itaylor,  Jan 30, 2010)
//        - https://github.com/golang/go/blob/659966a988e38f34b03a7c87780b594e4638f3b9/doc/gccgo_install.html (itaylor,  Aug 24, 2010)
//        - https://github.com/golang/go/blob/4164d60cc24c2889a9a40750891c61c537bbeb1b/doc/go_faq.html (adg Sep 29, 2010)
//        - gccgo run-time support uses glibc.
//        - gccgo implements goroutines using segmented stacks, supported by recent modifications to the gold linker.
//        - Under `gccgo` an OS thread will be created for each goroutine, and `GOMAXPROCS` is effectively equal to the number
//          of running goroutines. Under `gc` you must set `GOMAXPROCS` to allow the runtime to utilise more than one OS thread.
//        - Gccgo can, with care, be linked with GCC-compiled C or C++ programs. However, because Go is garbage-collected it
//          will be unwise to do so, at least naively.
//    - 2011 GCC4.6.0 正式加入对go的支持 （Nay25 2011)
//        - https://gcc.gnu.org/onlinedocs/gcc-4.6.0/gccgo/
//    - 2012 go build中可以直接调gccgo (go0之前）
//        - triggered by GC=gccgo in environment.
//        - https://github.com/golang/go/commit/45a8fae996700a40bc671bc48e78931d277dee0a
//        - go: introduce support for "go build" with gccgo. (Jan 28, 2012)
//    - 2013 GCC 4.8.2 release, gccgo implements Go 1.1.2 (Dec 2013)
//        - https://github.com/golang/go/blob/8d206d9d804d793727a668a21711dac69689cf23/doc/gccgo_install.html
//        - GCC 4.7.1 release and all later 4.7 releases include a complete go1
//    - 2013 gc开始使用go语言重写（go1.3~go1.5)
//       - Gc written in C move to Go (golang.org/s/go13compiler rsc)
//    -  GCC 4.9 -> Go 1.2
//    -  GCC 4.10 -> Go 1.3
//    - 2015-04 use gccgo as bootstrap compiler (go1.5beta1)
//      - https://github.com/golang/go/commit/67805eaa950c5318bfe02943cc175da6729919a9 (davecheney Apr 14, 2015)
//      - https://github.com/golang/go/blob/deb6c5b9200137423b9c594ff6a03bcc848a852e/doc/install-source.html ( Jun 10, 2015)
//      - 到了2016年10月bootstrap统一为下载包 (plus accumulated fixes to keep the tools running on newer operating systems)
//         - https://github.com/golang/go/commit/b990558162fa038f3651dc0f1821f64b282dda6f (go1.4-bootstrap-20161024.tar.gz)
//    - 2015-06 加入go/internal/gccgoimporter (go1.5beta1)
//      - https://github.com/golang/go/commit/f6ae5f96c7a1919d2d6f4d658737bc5082d9e996 (Jun 18, 2015)
//      - "go/importer":               {"L4", "go/internal/gcimporter", "go/internal/gccgoimporter", "go/types"},
//+       - "go/internal/gcimporter":    {"L4", "OS", "go/build", "go/constant", "go/token", "go/types", "text/scanner"},
//+       - "go/internal/gccgoimporter": {"L4", "OS", "debug/elf", "go/constant", "go/token", "go/types", "text/scanner"},
//    - 2015-06 GCC 5 -> Go 1.4 (user libraries. runtime is not fully merged, but that should not be visible to Go programs.)
//      - https://github.com/golang/go/blob/8668ac081a8596af97fbeed1bd3ae74d5e93e7d1/doc/gccgo_install.html (ian (Jun 16, 2015)
//      -  As of GCC 5 the gccgo installation also includes a version of the `go` command
//      -  GCC 5.1 -> April 22, 2015
//    - 2015-09 go/type的支持（类型检查，compiler的前端）支持gc/gccgo （go1.6beta1)
//      - https://github.com/golang/go/commit/7e1d1f899cbfb302447324e96da0e913ef94096c (gri Sep 30, 2015)
//      - type和importer这些，是gc和gccgo的共同依赖，他们属于compiler的前端。
//    - 2016-03 GCC 6 -> Go 1.6.1 user libraries （跳过了go 1.5)
//      - https://github.com/golang/go/commit/fb9aafac97649a11301b78ee9e2139804c52b528 (ian go1.7beta1, Mar 11, 2016)
//      - 6372c821c7cada20e662f9cc37f5d1c202a6b5fe 1.6.1
//      - GCC 6.1 -> April 27, 2016
//    - 2017-05 GCC 7 -> Go 1.8.3
//      - GCC 7.1 -> May 2, 2017
//    - 2017-05 go/types的变化 (应该解释为refactor，而不是不支持gccgo了）
//      - https://github.com/golang/go/commit/4be4da6331b4acfc379113bd5603079a4f36741a (gri go1.9beta1, Mar1, 2017)
//      - go/types: change local gotype command to use source importer, remove -gccgo flag (not supported after 1.5)
//    - 2017-09 importcfg sym links for gccgo (reading an importcfg directly)
//      - https://github.com/golang/go/commit/d8efa0e0ed8bbd5ed0780527652d86be2fba99dc (rsc go1.10beta1 Sep 29, 2017)
//      - builds in root a tree of symlinks implementing the directives from importcfg. This serves as a temporary
//        transition mechanism until we can depend on gccgo reading an importcfg directly. (The Go 1.9 and later gc
//        compilers already do.)
//    - 2017-10/11 rsc的cmd/go的refactor (gccgo.go独立源文件的首次出现）
//      - https://github.com/golang/go/commit/4e8be99590d54cc9ea949d9eadc560a1c2456539 (rsc go1.10beta1 Oct 12, 2017)
//        - cmd/go: clean up compile vs link vs shared library actions
//      - https://github.com/golang/go/commit/08362246b6a1dd7bffeb68f1d5116b9b2fe6209a (rsc go1.10beta1 Oct 21, 2017)
//        - cmd/go/internal/work: factor build.go into multiple files
//          -> src/cmd/go/internal/work/gc.go - gc toolchain
//          -> src/cmd/go/internal/work/gccgo.go - gccgo toolchain (这个结构一直沿用至今 May 30,2021)
//      - https://github.com/golang/go/commit/5993251c015dfa1e905bdf44bdb41572387edf90 (rsc go1.10beta1 Nov 9, 2017)
//        - cmd/go: implement per-package asmflags, gcflags, ldflags, gccgoflags
//    - 2018-01 iran -> cmd/go的进一步修复 -> go1.10
//       - https://github.com/golang/go/commit/fc408b620a1488323d8ac456a18685888541ac2d (ian go1.10beta2 Jan 6, 2018)
//         - cmd/go 支持 GNU build ID
//         - gccToolID returns the unique ID to use for a tool that is invoked by the GCC driver. This is in
//           particular gccgo, but this can also be used for gcc, g++, gfortran, etc.; those tools all use the GCC
//           driver under different names. The approach used here should also work for sufficiently new versions of clang.
//           Unlike toolID, the name argument is the program to run. The language argument is the type of input file as
//           passed to the GCC driver's -x option.
//       - https://github.com/golang/go/commit/9745eed4fd4160cfbf55e9dbbfa99aca5563b392 (ian go1.10rc1 Jan 12, 2018)
//         - cmd/go: make gccgo -buildmode=shared and -linkshared work again
//    - GCC 8, contain the Go 1.10.1 （跳过了go1.9)
//       -  GCC 8.1 -> May 2, 2018
//    - 2018-05/6 -> go1.11
//       - https://github.com/golang/go/commit/498c803c19caa94d9d37eb378deed786117bbeab (ian go1.11beta1 May 4, 2018)
//         - cmd/go, go/build: add support for gccgo tooldir
//         - The gccgo toolchain does not put tools (cgo, vet, etc.) in $GOROOT/pkg/tool, but instead in a directory
//           available at runtime.GCCGOTOOLDIR.
//       - https://github.com/golang/go/commit/d540da105c799a8fa010ee83419d6cb24d6627b4 (ian go1.11beta1 May 10, 2018)
//          - go/build, cmd/go: don't expect gccgo to have GOROOT packages
//          - When using gccgo the standard library sources are not available in GOROOT. Don't expect them to be there.
//           In the gccgo build, use a set of standard library packages generated at build time.
//       - https://github.com/golang/go/commit/30b6bc30b208299e4cb6598be854ec276db85661 (ian go1.11beta1 May 25, 2018)
//          - cmd/go, cmd/vet, go/internal/gccgoimport: make vet work with gccgo
//       - https://github.com/golang/go/commit/0c9be48a90bfafac68cde05c4d7db8eee17492f6 (ian go1.11beta1 Jun 20, 2018)
//          - go/internal/gccgoimporter: read export data from archives
//          - gccgo will normally generate archive files, needed by, cmd/vet, when typechecking packages.
//     - 2018-09 iran go/build 结构重整, gri,importer改进 -> go1.12
//       - https://github.com/golang/go/commit/b83ef36d6acb351ac50c5c7199fd683fb5226983 (ian go1.12beta1 Sep 26, 2018)
//          - go/build: move isStandardPackage to new internal/goroot package
//           -> src/internal/goroot/gccgo.go
//           -> src/internal/goroot/gc.go
//       - https://github.com/golang/go/commit/699da6bd134c22ac174ec1accae9ae8218f873f7 (ian go1.12beta1 Sep 26, 2018)
//          - go/build: support Import of local import path in standard library for gccgo
//       - https://github.com/golang/go/commit/9c81402f58ae83987f32153c1587c9f03b4a5769 (gri go1.12beta1 Sep 27, 2018)
//          - go/internal/gccgoimporter: fix updating of "forward declared" types
//     - GCC9 -> Go 1.12.2 (跳过了go1.11)
//       - https://github.com/golang/go/blob/77aea691762f46daeb56c2e1fe764fd3898fff6b/doc/gccgo_install.html (Alberto Donizetti go1.14beta1 Sep 13, 2019)
//       - GCC 9.1 -> May 3, 2019
//     - 2019-02/5/6/8 -> go1.13
//       - https://github.com/golang/go/commit/98cbf45cfc6a5a50cc6ac2367f9572cb198b57c7 (ian go1.13beta1 Feb 27, 2019)
//         - go/types: add gccgo sizes information
//       - https://github.com/golang/go/commit/451cf3e2cd8950571f436896a3987343f8c2d7f6 (gri go1.13beta1 May 14, 2019)
//         - spec: clarify language on package-level variable initialization
//         - cmd/compile, gccgo, and go/types produce different initialization orders
//       - https://github.com/golang/go/commit/7647fcd39292b5d36eb0f0be9750eecb03b1874c (ian go1.13beta1 Jun 7, 2019 )
//         - go/internal/gccgoimporter: update for gofrontend export data changes
//       - https://github.com/golang/go/commit/407010ef0b858a7fa6e6e95abe652fdff923da9a (ian go1.13rc1  Aug 1, 2019)
//         - cmd/go: only pass -fsplit-stack to gccgo if supported
//     - 2019-09 基本没什么改动的从go1.14到go1.15
//       - https://github.com/golang/go/commit/f668573a5e708db399688c9441cf5ec2eb2f29b0 (ian go1.14beta1 Sep 10, 2019)
//         - cmd/go: for gccgo, look for tool build ID before hashing entire file
//       - https://github.com/golang/go/commit/0a3b65c4926479c6ea2b8439cf073a43bfc2b9b6 (ian go1.14beta1 Sep 11, 2019)
//         - go/internal/gccgoimporter: support embedded field in pointer loop
//     - 2020-10/11/12 -> go1.16
//       - https://github.com/golang/go/commit/fe2cfb74ba6352990f5b41260b99e80f78e4a90a (randall go1.16beta1 Oct 2, 2020)
//         - go1.16删除387支持 dropping 387 floating-point support and requiring SSE2 support for GOARCH=386 in the
//           native gc compiler for Go 1.16. This would raise the minimum GOARCH=386 requirement to the Intel Pentium 4
//           (released in 2000) or AMD Opteron/Athlon 64 (released in 2003).
//         - https://github.com/golang/go/issues/40255
//       - https://github.com/golang/go/commit/72ee5bad9f9bd8979e14fab02fb07e39c5e9fd8c (ian go1.16beta1 Oct 6, 2020)
//         - cmd/cgo: split gofrontend mangling checks into cmd/internal/pkgpath
//         - 重构：增加src/cmd/internal/pkgpath/pkgpath.go
//       - https://github.com/golang/go/commit/a65bc048bf388e399af9bcfd726cd0f11bba7c8e (ian go1.16beta1 Oct 6, 2020)
//         - cmd/go: use cmd/internal/pkgpath for gccgo pkgpath symbol
//       - https://github.com/golang/go/commit/9fcb5e0c527337c830e95d48d4574930cac53093 (ian go1.16beta1 Oct 28, 2020)
//         - go/internal/gccgoimporter: support notinheap annotation
//           - 参考： runtime: update go:notinheap documentation
//           - https://github.com/golang/go/commit/2e0f8c379f91f77272d096929cf22391b64d0e34 (aclements go1.16beta1 Sep 25, 2020)
//           - 参考： `go:notinheap` 是在go1.8中加入的，
//           - https://github.com/golang/go/commit/77527a316b33d6f4c072c0774a1478bb53f42d35(Oct 16, 2016)
//       - https://github.com/golang/go/commit/c47eac7db00e03776c3975025184e1938fbced75 (ian go1.16beta1 Nov 21, 2020)
//         - cmd/cgo, cmd/internal/pkgpath: support gofrontend mangler v3
//         - The new version of gofrontend mangling scheme used by gccgo and GoLLVM .
//       - https://github.com/golang/go/commit/4f1b0a44cb46f3df28f5ef82e5769ebeac1bc493 (rsc go1.16beta1 Dec 10, 2020)
//         - rsc 整个结构的refactor, `io/ioutil` moved to `os`
//       - https://github.com/golang/go/commit/1b85e7c057e0ac20881099eee036cef1d7f2fbb0 (ian go1.16rc1 Jan 6 2021)
//         - cmd/go: don't scan gccgo standard library packages for imports
//     - 2021- now -> go1.17
//       - https://github.com/golang/go/commit/92c9f3a9b80eda50a55f8860587c2ed7734f4a29 (ian go1.17beta1 Apr 28 2021)
//         - cmd/go: include C/C++/Fortran compiler version in build ID
//
// - gofrontend 时间线
//   - libgo: update to Go1.16.3 release (Apr 13 2021)
//      - https://github.com/golang/gofrontend/commit/9782e85bef1c16c72a4980856d921cea104b129c
//   - libgo: update to Go1.16.2 release (Apr 13 2021)
//      - https://github.com/golang/gofrontend/commit/10b00ad87303d37c68b2d54dd25d655bd316946e
//   - libgo: update to Go1.16 release (Feb 20 2021)
//      - https://github.com/golang/gofrontend/commit/78a840e4940159a66072237f6b002ab79f441b79
//   - libgo: update to Go1.16rc1 (Jan 29, 2021)
//      - https://github.com/golang/gofrontend/commit/c7dcc3393018d55bf0e2740e1d423b80f48acf9f
//   - libgo: update to Go1.16beta1 release (Dec 31, 2020)
//      - https://github.com/golang/gofrontend/commit/49b07690a93a813585047d8265e326de8efd95c6
//   - libgo : update to Go1.15.6 (Dec 9,2020)
//      - https://github.com/golang/gofrontend/commit/0d0b423739b2fee9788cb6cb8af9ced29375e545
//   - libgo : update to Go 1.15.5 (Nov21 2020)
//      - https://github.com/golang/gofrontend/commit/36a7b789130b415c2fe7f8e3fc62ffbca265e3aa
//   - 1.15.4 (Nov10 2020)
//      - https://github.com/golang/gofrontend/commit/893fa057e36ae6c9b2ac5ffdf74634c35b3489c6
//   - libgo: update to Go 1.15.3 release (Oct 28, 2020)
//      - https://github.com/golang/gofrontend/commit/be0d2cc2df9f98d967c242594838f86362dae2e7
//   - 1.15.2 (Sep24 2020)
//      - https://github.com/golang/gofrontend/commit/6a7648c97c3e0cdbecbec7e760b30246521a6d90
//   - libgo: update to Go1.14beta1 (Jan 22, 2020)
//      - https://github.com/golang/gofrontend/commit/c2225a76d1e15f28056596807ebbbc526d4c58da
//   - libgo: update to Go1.13 (Sep 13, 2019)
//      - https://github.com/golang/gofrontend/commit/ceb1e4f5614b4772eed44f9cf57780e52f44753e
//   - libgo: update to Go 1.13beta1 release (Sep 7, 2019)
//      - https://github.com/golang/gofrontend/commit/8f2b844acda70330f7c50b360f8c983d2676ecbb
//   - libgo: update to Go 1.12 release (Feb 26, 2019)
//      - https://github.com/golang/gofrontend/commit/558fcb7bf2a6b78bdba87f20a8a4a95d27125d74
//   - libgo: update to Go 1.11 (Sep 25, 2018)
//      - https://github.com/golang/gofrontend/commit/7b25b4dff4778fc4d6b5d6e10594814146b3e5dd
//   - libgo: update to Go1.10beta1 (Jan 9, 2018)
//      - https://github.com/golang/gofrontend/commit/dbc0c7e4329aada2ae3554c20cfb8cfa48041213
//   - libgo: update to go1.9 (Sep 14, 2017)
//      - https://github.com/golang/gofrontend/commit/4e063a8eee636cce17aea48c7183e78431174de3
//   - libgo: update to final Go 1.8 release (Feb 17, 2017)
//      - https://github.com/golang/gofrontend/commit/893f0e4a707c6f10eb14842b18954486042f0fb3
//   - libgo: Update to go1.6rc1. (Feb 4, 2016)
//      - https://github.com/golang/gofrontend/commit/8cedea1bc3ac014843343ebb1b2e4e2a9a0f4d78
//   - libgo: Update to Go 1.5.1. (Oct 31, 2015)
//      - https://github.com/golang/gofrontend/commit/5e7ded0b52f32485722b3e475435f2b6d6187d5d
//   - libgo: Upgrade to Go 1.4.2 release. (Mar 7, 2015)
//      - https://github.com/golang/gofrontend/commit/c8cbd88101d3478ab4a844d3648f71892787bbce
//   - libgo: Update to Go 1.3.3 release. (Oct 28, 2014)
//      - https://github.com/golang/gofrontend/commit/d0b65bc98dfd968a103b9bd08e3c01f33498c220
//   - libgo: Update to Go 1.3 release. (Jul 19, 2014)
//      - https://github.com/golang/gofrontend/commit/b95c63324e743bd130eaf324bc916a1f29ceeef9
// - libgo: Update to Go 1.2.1 release. (Mar 4, 2014)
//      - https://github.com/golang/gofrontend/commit/08b18f1c15cf7c9996496b0144fcf7d660283ced
// - libgo: Update to Go 1.1.1. (Jul 16, 2013)
//      - https://github.com/golang/gofrontend/commit/093f55bc7f9087e20f03cca154dd329f2dcd0fe3
// - libgo: Update to Go 1.0.3.(Oct 3, 2012)
//      - https://github.com/golang/gofrontend/commit/44c49602f3e3f2cd65f6114c8f81275a3690ec74
// - libgo: Update to Go 1.0.2 release. (Jun 25, 2012)
//      - https://github.com/golang/gofrontend/commit/52effb1120f8d200b5e9f2407c4255057b02d36d
// - libgo: Update to Go 1.0.1 release. (May 4, 2012)
//      - https://github.com/golang/gofrontend/commit/7b3ed26d5d281e084fb9f222c6d02c27b7d6ef28
// - libgo: Update to weekly.2012-03-27 aka go1 release. (Mar 31, 2012)
//      - https://github.com/golang/gofrontend/commit/d3f03487312b546cba7f66ef43f9a38fca5d7ef0
// - Update to current Go library, part 8.(Aug 27, 2010)
//      - https://github.com/golang/gofrontend/commit/1c9427aef0e7b2bf4687f311e5416ce73cd27d2b
// - Initial import of gofrontend repository. (Jan 30, 2010)
//      - https://github.com/golang/gofrontend/commit/0ef89c4e8b1f5c66ab6c9a6c307bc153e0f1b0f4
//

//==================================================================================================
// 附注：范型支持 go1.17(trunk)对于型别的重大变化： type2/importer/noder2/irgren/
// ==================================================================================================
// 关于范型提案
//   - https://github.com/golang/go/issues/43651 (lanlancetaylor, Jan 13, 2021)
//      spec: add generic programming using type parameters
//   - https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md （March 19, 2021）
//     Type Parameters Proposal
//
// 为什么go1不支持范型 （Oct6, 2009, robpike）：
//    - https://github.com/golang/go/blob/2b1dbe8a4f9f214d7164abd18d99e2451efc5cdb/doc/go_lang_faq.html
//    - Generics may well be added at some point.  We don't feel an urgency for
//      them, although we understand some programmers do.
//      Generics are convenient but they come at a cost in
//      complexity in the type system and run-time.  We haven't yet found a
//      design that gives value proportionate to the complexity, although we
//      continue to think about it.  Meanwhile, Go's built-in maps and slices,
//      plus the ability to use the empty interface to construct containers
//      (with explicit unboxing) mean in many cases it is possible to write
//      code that does what generics would enable, if less smoothly.
//      This remains an open issue.
//
// types和types2的区别 (1.17 trunk 目的 go的generic支持)
//   - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/types/type.go
//   - https://github.com/golang/go/blob/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/types2/type.go
//
// 关于importer (cmd/compile/internal/importer)
//   - https://github.com/golang/go/tree/3075ffc93e962792ddf43b2a528ef19b1577ffb7/src/cmd/compile/internal/importer
//   - 这个是type2的一部分，和type2一起被加入。即原有的gcimporter(go/internal/gcimporter)针对types2的重写
//     - 参考 : https://github.com/golang/go/commit/ca36ba83ab86b9eb1ddc076f0ebfda648ce31d6b#diff-6b37a612c3d516286caf7a7a0e7a00a28c986040a2485ae573e3a3f1e33ed2bb
//     - 对比 : https://github.com/golang/go/tree/639acdc833bfd12b7edd43092d1b380d70cb2874/src/go/internal/gcimporter
//     - 代码是相同的：只是接口上的改变：去掉了fset *token.FileSet
//   - 功能是 Import for gc-generated object files.
//   - cmd/compile/internal/typecheck/iexport.go 定义了 export data format.

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
//
//
// ==================================================================================================
// 关于deps_test.go:
//  - 方便理解go包的依赖关系和主体结构
//  - 相当于使用go构造任何东西的基础
//  - 相当于从外界看go的整体的鸟瞰。
// ==================================================================================================
// trunk
//  - https://github.com/golang/go/blob/master/src/go/build/deps_test.go
// 没有修改前最后一版：
//  - https://github.com/golang/go/blob/6f52790a20a2432ae61e0ec9852a3df823a16d40/src/go/build/deps_test.go (Apr3,2020)
//  - L0 the lowest level, core, nearly unavoidable packages.
//      		"errors",
//      		"io",
//      		"runtime",
//      		"runtime/internal/atomic",
//      		"sync",
//      		"sync/atomic",
//      		"unsafe",
//      		"internal/cpu",
//      		"internal/bytealg",
//      		"internal/reflectlite",
//  - L1  simple functions and strings processing, but not Unicode tables.
//      		"math",
//      		"math/bits",
//      		"math/cmplx",
//      		"math/rand",
//      		"sort",
//      		"strconv",
//      		"unicode/utf16",
//      		"unicode/utf8",
//   - L2: Unicode and strings processing.
//      		"bufio",
//      		"bytes",
//      		"path",
//      		"strings",
//      		"unicode",
// - L3 adds reflection and some basic utility packages and interface definitions, but nothing that makes system calls.
//      		"crypto",
//      		"crypto/cipher",
//      		"crypto/internal/subtle",
//      		"crypto/subtle",
//      		"encoding/base32",
//      		"encoding/base64",
//      		"encoding/binary",
//      		"hash",
//      		"hash/adler32",
//      		"hash/crc32",
//      		"hash/crc64",
//      		"hash/fnv",
//      		"image",
//      		"image/color",
//      		"image/color/palette",
//      		"internal/fmtsort",
//      		"internal/oserror",
//      		"reflect",
//   - Operating system access. (syscall, time)
//   - OS enables basic operating system functionality, but not direct use of package syscall, nor os/signal.
//      		"io/ioutil",
//      		"os",
//      		"os/exec",
//      		"path/filepath",
//      		"time",
//   - Formatted I/O: few dependencies (L1) but we must add reflect and internal/fmtsort.
//   - Packages used by testing must be low-level (L2+fmt).
//
// - L4 is defined as L3+fmt+log+time, because in general once you're using L3 packages,
//   use of fmt, log, or time is not a big deal.
//      		"fmt",
//      		"log",
//      		"time",
//   - Go parser.
//      	"go/ast":     {"L4", "OS", "go/scanner", "go/token"},
//      	"go/doc":     {"L4", "OS", "go/ast", "go/token", "regexp", "internal/lazyregexp", "text/template"},
//      	"go/parser":  {"L4", "OS", "go/ast", "go/scanner", "go/token"},
//      	"go/printer": {"L4", "OS", "go/ast", "go/scanner", "go/token", "text/tabwriter"},
//      	"go/scanner": {"L4", "OS", "go/token"},
//      	"go/token":   {"L4"},
//   - Go type checking.
//      	"go/constant":               {"L4", "go/token", "math/big"},
//      	"go/importer":               {"L4", "go/build", "go/internal/gccgoimporter", "go/internal/gcimporter", "go/internal/srcimporter", "go/token", "go/types"},
//      	"go/internal/gcimporter":    {"L4", "OS", "go/build", "go/constant", "go/token", "go/types", "text/scanner"},
//      	"go/internal/gccgoimporter": {"L4", "OS", "debug/elf", "go/constant", "go/token", "go/types", "internal/xcoff", "text/scanner"},
//      	"go/internal/srcimporter":   {"L4", "OS", "fmt", "go/ast", "go/build", "go/parser", "go/token", "go/types", "path/filepath"},
//      	"go/types":                  {"L4", "GOPARSER", "container/heap", "go/constant"},
//