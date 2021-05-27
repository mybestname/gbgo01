package main

func m() {
	s1 := make([]int, 5)               // [0,0,0,0,0]  len=5, cap=5,
	_ = s1
	s2 :=  make([]int, 8192)          // force to heap
	// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ir/cfg.go#L19
	// MaxImplicitStackVarSize = int64(64*1024) = 65536
	// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/escape/escape.go#L2036
	// []int/[]int64 ->  maxImplicitStackVarSize/t.Elem().width = 65536/8(int/int64) = 8192
	_ = s2
}

//go:generate go tool compile -N -l -m=2 main1.go
//go:generate go tool objdump main1.o
//go:generate rm main1.o

//
// 关于make slice的实现 :
//
// ====================================TL;DR==============================
// 编译器会替换原始代码，来生成中间代码，来最终构造出slice。
//  1. parser（语法分析，构造AST）: src -> AST
//  2. walker（优化/遍历AST，生成IR）: OMAKE -> OMAKESLICE -> OSLICE(stack)/OSLICEHEADER(heap) (如果涉及到runtime，加入runtime，如堆上分配）
//  3. ssagen/ssa (IR进一步优化）: OSLICE/OSLICEHEADER -> ssa.OpSliceMake etc. -> 50 多个优化过程 -> 伪汇编码
//  4. asm （真正汇编） :  加入arch -> obj
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
// - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/walk.go#L24
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
// 2.3. `walkMakeSlice` len和cap的各种处理，根据是否逃逸，完成需要分配内存空间的计算并生成SliceHeader (OSLICEHEADER)
//  - 不逃逸的情况下，构造成arr声明加slice引用的形式，构造OSLICE
//  - 逃逸的情况下，调用runtime.makeslice，堆上内存分配，构造OSLICEHEADER
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/builtin.go#L364-L431
// ```
//  // walkMakeSlice walks an OMAKESLICE node.
//  func walkMakeSlice(n *ir.MakeExpr, init *ir.Nodes) ir.Node {
//  	......
//  	//不逃逸的情况下，构造成arr声明加slice引用的形式，构造OSLICE
//  	// var arr [r]T
//  	// n = arr[:l]
//  	if n.Esc() == ir.EscNone {
//  	......
//  		i := typecheck.IndexConst(r)
//  	......
//  		t = types.NewArray(t.Elem(), i) // [r]T
//  	......
//  		r := ir.NewSliceExpr(base.Pos, ir.OSLICE, var_, nil, l, nil) // arr[:l]
//  		return walkExpr(typecheck.Expr(typecheck.Conv(r, n.Type())), init)
//  	}
//  	......
//  	// 逃逸的情况下，调用runtime.makeslice，准备进行堆上内存分配。
//  	len, cap := l, r
//  	fnname := "makeslice64"
//  	argtype := types.Types[types.TINT64]
//  	......
//  	if (....) {
//  		fnname = "makeslice"
//  		argtype = types.Types[types.TINT]
//  	}
//  	fn := typecheck.LookupRuntime(fnname)                      //调用runtime.makeslice/makeslice64 进行内存计算，构造指针结构
//  	ptr := mkcall1(fn, types.Types[types.TUNSAFEPTR], init, reflectdata.TypePtr(t.Elem()), typecheck.Conv(len, argtype), typecheck.Conv(cap, argtype))
//  	ptr.MarkNonNil()
//  	len = typecheck.Conv(len, types.Types[types.TINT])
//  	cap = typecheck.Conv(cap, types.Types[types.TINT])
//  	sh := ir.NewSliceHeaderExpr(base.Pos, t, ptr, len, cap)    //构造OSLICEHEADER
//  	return walkExpr(typecheck.Expr(sh), init)
//  }
// ```
// 2.3.1 如果不逃逸（slice小，可以栈上分配）：
//  构造OSLICE的情况
//
// 2.3.2 关于runtime.slice的调用，该方法可以构造出指针（unsafe.Pointer结构），用来处理在堆上分配的情景。
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/runtime/slice.go#L83-L113
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/runtime/malloc.go#L903 mallocgc
// 这个指针如何alloc（Memory allocation）栈上如何分配，堆上如何分配。
//
// 2.3.2.1 注意mkcall1的用法，最终在代码树中构造了CallExpr
// 这里的目的是构造slice的runtime结构。
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/builtin.go#L425 makcall1(makeslice64,ptr)
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/walk.go#L130 mkcall1()
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/walk.go#L103 vmkcall()
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/expr.go#L546
// -> 构造CallExpr
//
// 2.3.2.2 typecheck.Expr(SliceHeaderExpr)
// -> 构造SliceHeaderExpr
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/typecheck/typecheck.go#L706 case ir.OSLICEHEADER:
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/typecheck/expr.go#L797
//
//
//
//
// 2.4 walkSlice和walkSliceHedler，进一步处理OSLICE(栈上）或 OSLICEHEADER（堆上）
// 2.4.1 `walkSlice` 处理slices(OSLICE), 构造 `SliceExpr`
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/expr.go#L802-L837
// ```
//  // walkSlice walks an OSLICE, OSLICEARR, OSLICESTR, OSLICE3, or OSLICE3ARR node.
//  func walkSlice(n *ir.SliceExpr, init *ir.Nodes) ir.Node {
//  	......
//  	n.Low = walkExpr(n.Low, init)
//  	n.High = walkExpr(n.High, init)
//  	n.Max = walkExpr(n.Max, init)
//  	......
//     // Reduce x[i:j:cap(x)] to x[i:j].
//  	......
//  	return reduceSlice(n)
//  }
//   ```
//
// 2.4.2 `walkSliceHeader` 处理slice头(OSLICEHEADER)，构造`SliceHeaderExpr`
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
// 2.4.2.1 `walkExpr`继续处理OSPTR，构造`UnaryExpr`
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ir/node.go#L310
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/expr.go#L100
//
//
//
// 3. SSA生成。语法树变成中间代码，经过几十轮过程，最终生成类似于汇编的中间代码。
// 3.0 入口 ssagen.Compile
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/gc/compile.go#L153  ssagen.Compile(fn, worker)
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/pgen.go#L165 Compile(fn,worker)
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L642  buildssa(fn,worker)
//  - https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssa/compile.go#L29  ssa.Compile(fn)
//
// 3.1 ir.OSLICE 或 ir.OSLICEHEADER -> ssa.OpSliceMake
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L3131-L3152
// ```
//  	......
//  	case ir.OSLICEHEADER:
//  		n := n.(*ir.SliceHeaderExpr)
//  		p := s.expr(n.Ptr)
//  		l := s.expr(n.Len)
//  		c := s.expr(n.Cap)
//  		return s.newValue3(ssa.OpSliceMake, n.Type(), p, l, c)
//  	......
//  	case ir.OSLICE, ir.OSLICEARR, ir.OSLICE3, ir.OSLICE3ARR:
//  		......
//  		n := n.(*ir.SliceExpr)
//  		......
//  		i = s.expr(n.Low)
//  		......
//  		j = s.expr(n.High)
//  		......
//  		k = s.expr(n.Max)
//  		p, l, c := s.slice(v, i, j, k, n.Bounded())
//  		......
//  		return s.newValue3(ssa.OpSliceMake, n.Type(), p, l, c)
//  	......
// ```
// 3.1.1 如果是ir.OSLICE的，在加入SliceMake前，还会加入OpSlicePtr，OpAddPtr等。
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L5784
// ```
//  // slice computes the slice v[i:j:k] and returns ptr, len, and cap of result.
//  // i,j,k may be nil, in which case they are set to their default value.
//  // v may be a slice, string or pointer to an array.
//  func (s *state) slice(v, i, j, k *ssa.Value, bounded bool) (p, l, c *ssa.Value) {
//  	 t := v.Type
//  	......
//  	case t.IsSlice():
//  		ptr = s.newValue1(ssa.OpSlicePtr, types.NewPtr(t.Elem()), v)
//  		len = s.newValue1(ssa.OpSliceLen, types.Types[types.TINT], v)
//  		cap = s.newValue1(ssa.OpSliceCap, types.Types[types.TINT], v)
//  	......
//  	subOp := s.ssaOp(ir.OSUB, types.Types[types.TINT])
//  	......
//  	// Calculate the length (rlen) and capacity (rcap) of the new slice.
//  	rlen := s.newValue2(subOp, types.Types[types.TINT], j, i)
//  	rcap := rlen
//  	if j != k && !t.IsString() {
//  		rcap = s.newValue2(subOp, types.Types[types.TINT], k, i)
//  	}
//  	if (i.Op == ssa.OpConst64 || i.Op == ssa.OpConst32) && i.AuxInt == 0 {
//  		// No pointer arithmetic necessary.  //普通情况
//  		return ptr, rlen, rcap
//  	}
//  	......
//  	// 需要指针运算的复制情况，忽略
//  	// Calculate the base pointer (rptr) for the new slice.
//  	......
//  	rptr := s.newValue2(ssa.OpAddPtr, ptr.Type, ptr, delta)
//  	return rptr, rlen, rcap
//  }
// ```
// 3.2 如果是OSLICEHEADER, 在加入SliceMake前, 分别处理n.Ptr, n.Len, n.Cap
//    - ir.OSPTR -> ssa.OpSlicePtr
//    - ir.OLen  -> ssa.OpSliceLen
//    - ir.OCap  -> ssa.OpSliceCap
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L3106-L3113
// ```
//  	case ir.OSPTR:
//  		n := n.(*ir.UnaryExpr)
//  		a := s.expr(n.X)
//  		if n.X.Type().IsSlice() {
//  			return s.newValue1(ssa.OpSlicePtr, n.Type(), a)
// 		......
//```
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L3089-L3104
// ```
//  		......
//   	case ir.OLEN, ir.OCAP:
//  		n := n.(*ir.UnaryExpr)
//  		switch {
//  		case n.X.Type().IsSlice():
//  			op := ssa.OpSliceLen
//  			if n.Op() == ir.OCAP {
//  				op = ssa.OpSliceCap
//  			}
//  			return s.newValue1(op, types.Types[types.TINT], s.expr(n.X))
//   		......
// ```
//
// 3.3 关于runtime.makeslice 加入的CallExpr的处理
//   ir.OCALLFUNC (CallExpr) -> ssa.OpStaticLECall
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L1443-L1453
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L4914
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L5143-L5145
//
// 3.4 ssa.html输出分析
// 执行命令：
// GOSSAFUNC=m go build -gcflags "-N -l" main1.go
// 可以从生成的ssa.html中可以观察到这一系列的中间过程代码。其中expend-call过程中可以看到:
// 3.3.1 `s := make([]int, 5)` 分析
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
//
// 3.3.2 `make([]int, 8192)` 分析
// 可以看出 runtime.makeslice 被包裹在 StaticLECall op中。同时比较 SliceMake 的情况。这回的指针不在是一个本地地址（stack），而是 unsafe.Pointer 结构（代表heap）分配。
// ```
// ......
//  v26 (6) = StaticLECall <unsafe.Pointer,mem> {AuxCall{runtime.makeslice([*byte,0],[int,8],[int,16])[unsafe.Pointer,24]}} [32] v24 v25 v25 v21
//  v27 (6) = SelectN <mem> [1] v26
//  v28 (6) = SelectN <unsafe.Pointer> [0] v26
//  v29 (6) = SliceMake <[]int> v28 v25 v25
//  v30 (6) = VarDef <mem> {s2} v27
//  v31 (6) = LocalAddr <*[]int> {s2} v2 v30
//  v32 (6) = Store <mem> {[]int} v31 v29 v30
// ......
// ```
// 同样最终的伪汇编码也会清除这些中间过程的 StaticLECall 和 SliceMake 等等。
// ```
//   00013 (+6) LEAQ type.int(SB), AX
//   00014 (6) MOVQ AX, (SP)
//   00015 (6) MOVQ $8192, 8(SP)
//   00016 (6) MOVQ $8192, 16(SP)
//   00017 (+6) PCDATA $1, $0
//   00018 (+6) CALL runtime.makeslice(SB)
//   00019 (6) MOVQ 24(SP), AX
//   00020 (6) MOVQ AX, "".s2-48(SP)
//   00021 (6) MOVQ $8192, "".s2-40(SP)
//   00022 (6) MOVQ $8192, "".s2-32(SP)
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
//  main1.go:4			48c744244000000000	MOVQ $0x0, 0x40(SP)
//  main1.go:4			0f57c0			XORPS X0, X0
//  main1.go:4			0f11442448		MOVUPS X0, 0x48(SP)
//  main1.go:4			0f11442458		MOVUPS X0, 0x58(SP)
//  main1.go:4			488d442440		LEAQ 0x40(SP), AX
//  main1.go:4			8400			TESTB AL, 0(AX)
//  main1.go:4			eb00			JMP 0x501
//  main1.go:4			4889442468		MOVQ AX, 0x68(SP)
//  main1.go:4			48c744247005000000	MOVQ $0x5, 0x70(SP)
//  main1.go:4			48c744247805000000	MOVQ $0x5, 0x78(SP)
// ```
// ```
//  main1.go:6			488d0500000000		LEAQ 0(IP), AX		[3:7]R_PCREL:type.int
//  main1.go:6			48890424		MOVQ AX, 0(SP)
//  main1.go:6			48c744240800200000	MOVQ $0x2000, 0x8(SP)
//  main1.go:6			48c744241000200000	MOVQ $0x2000, 0x10(SP)
//  main1.go:6			e800000000		CALL 0x4d6		[1:5]R_CALL:runtime.makeslice<1>
//  main1.go:6			488b442418		MOVQ 0x18(SP), AX
//  main1.go:6			4889442448		MOVQ AX, 0x48(SP)
//  main1.go:6			48c744245000200000	MOVQ $0x2000, 0x50(SP)
//  main1.go:6			48c744245800200000	MOVQ $0x2000, 0x58(SP)
// ```