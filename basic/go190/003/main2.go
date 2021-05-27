package main

func main() {

}

// 关于new的实现
//
//
//
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/walk/builtin.go#L486-L500
// `func walkNew(n *ir.UnaryExpr, init *ir.Nodes) ir.Node`
//
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L3207-L3209
// ```
// 	case ir.ONEW:
//		n := n.(*ir.UnaryExpr)
//		return s.newObject(n.Type().Elem())
// ```
// https://github.com/golang/go/blob/d050238bb653711b47335583c5425c9efec30e4e/src/cmd/compile/internal/ssagen/ssa.go#L764-L769
// ```
// 	// newObject returns an SSA value denoting new(typ).
//  func (s *state) newObject(typ *types.Type) *ssa.Value {
//  	if typ.Size() == 0 {
//  		return s.newValue1A(ssa.OpAddr, types.NewPtr(typ), ir.Syms.Zerobase, s.sb)
//  	}
//  	return s.rtcall(ir.Syms.Newobject, true, []*types.Type{types.NewPtr(typ)}, s.reflectType(typ))[0]
//  }
// ```
//
// 注：根据最新trunk, ONEWOBJ(go1.16和之前）已经被删除，转而使用 ONEW
// - https://github.com/golang/go/commit/ab3b67abfd9bff30fc001c966ab121bacff3de9b
// - https://go-review.googlesource.com/c/go/+/284117
//   cmd/compile: remove ONEWOBJ
//   After CL 283233, SSA can now handle new(typ) without the frontend to
//   generate the type address, so we can remove ONEWOBJ in favor of ONEW
//   only.



