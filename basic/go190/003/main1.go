package main

func m() {
	s := make([]int, 5)               // [0,0,0,0,0]  len=5, cap=5,
	s = append(s, 1, 2, 3)    // [0,0,0,0,0,1,2,3]  len=8,cap=10
	println(s)
}

//go:generate go tool compile -N -l main1.go
//go:generate go tool objdump main1.o
//go:generate rm main1.o

//
// 关于make slice的实现 :
//
// ====================================TL;DR==============================
// 编译器会替换原始代码，来生成中间代码，来最终构造出slice。
//  1. src - AST
//  2. walker -> OMAKE -> OMAKESLICE -> OSLICEHEADER
//  3. ssagen/ssa: OSLICEHEADER -> ssa.OpSliceMake -> 伪汇编码
//  4. asm : 加入arch -> obj
// ====================================TL;DR==============================
//
// 1. 首先从原始代码进行代码分析，生成语法树。
// 1.0 go tool compile 入口点
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/main.go#L43
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/main.go#L55
// 1.1 InitUniverse() initializes the universe block.
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/typecheck/universe.go#L96
// 1.2 noder阶段： 词法解析和型别检查 ->
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/main.go#L192
// 1.3 逃逸分析
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/main.go#L253
//
// 2. walk阶段：语法树处理。通过walk系列函数，处理语法树的所有节点。
//   - 通过cmd/compile/internal/walk过程完成
//   - make会由OMAKE节点，首先变成OMAKESLICE、OMAKEMAP 和 OMAKECHAN这种节点，因为make只能处理这3种特殊的数据类型。
//   - 针对make slice来说就是处理代码树上的OMAKESLICE节点。
// 2.0 入口点：
// - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/main.go#L281   enqueueFunc(fn)
// - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/compile.go#L66 prepareFunc(fn)
// - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/compile.go#L92 walk.Walk(fn)
//
// 2.1 OMAKE相关Node 定义
//
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ir/node.go#L203-L214
//  ```
//  	OMAKE          // make(Args) (before type checking converts to one of the following)
//  	OMAKECHAN      // make(Type[, Len]) (type is chan)
//  	OMAKEMAP       // make(Type[, Len]) (type is map)
//  	OMAKESLICE     // make(Type[, Len[, Cap]]) (type is slice)
//  	OMAKESLICECOPY // makeslicecopy(Type, Len, Cap) (type is slice; Len is length and Cap is the copied from slice)
//  	// OMAKESLICECOPY is created by the order pass and corresponds to:
//  	//  s = make(Type, Len); copy(s, Cap)
//  	//
//  	// Bounded can be set on the node when Len == len(Cap) is known at compile time.
//  	//
//  	// This node is created so the walk pass can optimize this pattern which would
//  	// otherwise be hard to detect after the order pass.
//  	.......
//  	.......
//  	.......
//  	OSLICE       // X[Low : High] (X is untypechecked or slice)
//  	OSLICEARR    // X[Low : High] (X is pointer to array)
//  	OSLICESTR    // X[Low : High] (X is string)
//  	OSLICE3      // X[Low : High : Max] (X is untypedchecked or slice)
//  	OSLICE3ARR   // X[Low : High : Max] (X is pointer to array)
//  	OSLICEHEADER // sliceheader{Ptr, Len, Cap} (Ptr is unsafe.Pointer, Len is length, Cap is capacity)
//  ```
// 2.1  transformBuiltin()/transformMake()  OMAKE ->OMAKESLICE、OMAKEMAP 和 OMAKECHAN
// ```
// // Corresponds to Builtin part of tcCall.
// func transformBuiltin(n *ir.CallExpr) ir.Node {
//  ....
//  	case ir.OMAKE:
//  		return transformMake(n)
//  ....
// }
//  func {
//  .......
//  	case types.TSLICE:
// 			......
//  		nn = ir.NewMakeExpr(n.Pos(), ir.OMAKESLICE, l, r)
//  		.......
//  	case types.TMAP:
//  		.......
//  	case types.TCHAN:
//  .......
//  }
// ```
// 2.2 walkExpr()/walkExpr1()  遍历节点，发现并处理OMAKESLICE
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/expr.go#L261-L275
// ```
//  	func walkExpr1(n ir.Node, init *ir.Nodes) ir.Node {
//  	switch n.Op() {
//  		......
//  		case ir.OMAKECHAN:
//  		......
//  		case ir.OMAKEMAP:
//  		......
//  		case ir.OMAKESLICE:
//  			n := n.(*ir.MakeExpr)
//  			return walkMakeSlice(n, init)
//  		case ir.OMAKESLICECOPY:
//  		......
//  }
// ```
//
// 2.3. `walkMakeSlice` len和cap的各种处理，并生成SliceHeader (OSLICEHEADER)
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/builtin.go#L364-L431
// ```
// // walkMakeSlice walks an OMAKESLICE node.
//  func walkMakeSlice(n *ir.MakeExpr, init *ir.Nodes) ir.Node { ...
//  		......
//  		len, cap := l, r
//  		......
//  		sh := ir.NewSliceHeaderExpr(base.Pos, t, ptr, len, cap)
//  		......
//  }
// ```
//
// 2.4 `walkSliceHeader` 处理slice头，完成`SliceHeaderExpr`构造。
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/expr.go#L840
// ```
//  // walkSliceHeader walks an OSLICEHEADER node.
//   func walkSliceHeader(n *ir.SliceHeaderExpr, init *ir.Nodes) ir.Node {
//   	n.Ptr = walkExpr(n.Ptr, init)
//   	n.Len = walkExpr(n.Len, init)
//   	n.Cap = walkExpr(n.Cap, init)
//   	return n
//   }
// ```
//
//
// 3. SSA生成。语法树变成中间代码，经过几十轮过程，最终生成类似于汇编的中间代码。
// 3.0 入口 ssagen.Compile
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/compile.go#L153  ssagen.Compile(fn, worker)
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/pgen.go#L165 Compile(fn,worker)
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L642  buildssa(fn,worker)
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssa/compile.go#L29  ssa.Compile(fn)
// 3.1 ir.OSLICEHEADER -> ssa.OpSliceMake
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L3131-L3152
//```
//  	case ir.OSLICEHEADER:
//  		n := n.(*ir.SliceHeaderExpr)
//  		p := s.expr(n.Ptr)
//  		l := s.expr(n.Len)
//  		c := s.expr(n.Cap)
//  		return s.newValue3(ssa.OpSliceMake, n.Type(), p, l, c)
//  	case ir.OSLICE, ir.OSLICEARR, ir.OSLICE3, ir.OSLICE3ARR:
//  		return s.newValue3(ssa.OpSliceMake, n.Type(), p, l, c)
//```
//
//执行命令：
//GOSSAFUNC=m go build -gcflags "-N -l" main1.go
//可以从生成的ssa.html中可以观察到这一系列的中间过程代码。其中expend-call过程中可以看到针对`s := make([]int, 5)`变成了ssa IR的SliceMake描述
//```
//.....
//   v7 (4) = LocalAddr <*[5]int> {.autotmp_1} v2 v6
//   v8 (?) = Const64 <int> [5]
//  ......
//  v16 (4) = SliceMake <[]int> v7 v8 v8
//```
// SliceMake的描述也只是中间码生成的一步，会被替换为类汇编代码，最终的输出已经非常接近汇编代码。
// ```
// ......
//   00003 (+4) MOVQ $0, ""..autotmp_1-64(SP)
//   00004 (4) XORPS X0, X0
//   00005 (4) MOVUPS X0, ""..autotmp_1-56(SP)
//   00006 (4) MOVUPS X0, ""..autotmp_1-40(SP)
//   00007 (4) LEAQ ""..autotmp_1-64(SP), AX
//   00008 (4) TESTB AX, (AX)
//   00009 (4) JMP 10
//   00010 (4) MOVQ AX, "".s-24(SP)
//   00011 (4) MOVQ $5, "".s-16(SP)
//   00012 (4) MOVQ $5, "".s-8(SP)
// ......
// ```
// 4. 汇编过程，加入Arch相关的内容。生成真正的汇编代码。然后被汇编器翻译为机器语言。可以发现ssa的最终输出已经非常接近最终的汇编代码。
//
// 4.0 入口
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/pgen.go#L190 pp.Flush() -> assemble, boilerplate, etc.
//
// 4.1 Flushplist, (Build symbols, assign instructions )
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/internal/obj/plist.go#L22
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/internal/obj/plist.go#L126-L142
// ```
//  	// Turn functions into machine code images.
//  	for _, s := range text {
//  		...
//  		ctxt.Arch.Preprocess(ctxt, s, newprog)
//  		ctxt.Arch.Assemble(ctxt, s, newprog)
//  		...
//  	}
// ```
// 4.2 Arch.Preprocess/Arch.Assemble 加入arch相关代码
// amd64 Preprocess
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/asm/internal/arch/arch.go#L102
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/internal/obj/x86/obj6.go#L573
// amd64 Assemble
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/internal/obj/x86/obj6.go#L1385
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/internal/obj/x86/asm6.go#L2038 span6()
//
// 4.3 最终的机器码
// ```
//  main1.go:4		0x4e2			48c744244000000000	MOVQ $0x0, 0x40(SP)
//  main1.go:4		0x4eb			0f57c0			XORPS X0, X0
//  main1.go:4		0x4ee			0f11442448		MOVUPS X0, 0x48(SP)
//  main1.go:4		0x4f3			0f11442458		MOVUPS X0, 0x58(SP)
//  main1.go:4		0x4f8			488d442440		LEAQ 0x40(SP), AX
//  main1.go:4		0x4fd			8400			TESTB AL, 0(AX)
//  main1.go:4		0x4ff			eb00			JMP 0x501
//  main1.go:4		0x501			4889442468		MOVQ AX, 0x68(SP)
//  main1.go:4		0x506			48c744247005000000	MOVQ $0x5, 0x70(SP)
//  main1.go:4		0x50f			48c744247805000000	MOVQ $0x5, 0x78(SP)
// ```