package main

// https://stackoverflow.com/questions/9619479/go-what-determines-the-iteration-order-for-map-keys
// Map is implemented in Go as a hashmap.
//
// The Go run-time uses a common hashmap implementation which is implemented in C.
// The only implementation differences between map[string]T and map[byte]T are:
// hash function, equivalence function and copy function.
//
// Unlike (some) C++ maps, Go maps aren't fully specialized for integers and for strings.
//
// In Go release.r60, the iteration order is independent from insertion order as long as there are
// no key collisions. If there are collisions, iteration order is affected by insertion order.
// This holds true regardless of key type. There is no difference between keys of type string and
// keys of type byte in this respect, so it is only a coincidence that your program always printed
// the string keys in the same order. The iteration order is always the same unless the map is modified.
//
// However, in the newest Go weekly release (and in Go1 which may be expected to be released this month),
// the iteration order is randomized (it starts at a pseudo-randomly chosen key, and the hashcode
// computation is seeded with a pseudo-random number). If you compile your program with the weekly release
// (and with Go1), the iteration order will be different each time you run your program. That said,
// running your program an infinite number of times probably wouldn't print all possible permutations of
// the key set.
//
// Example outputs:
// ```
// stringMap keys: b 0 hello c world 10 1 123 bar foo 100 a
// stringMap keys: hello world c 1 10 bar foo 123 100 a b 0
// stringMap keys: bar foo 123 100 world c 1 10 b 0 hello a
// ...
// ```
// https://github.com/golang/go/blob/release.r60.3/src/pkg/runtime/hashmap.c
// https://github.com/golang/go/blob/weekly.2012-03-04/src/pkg/runtime/hashmap.c
//  - https://github.com/golang/go/commit/85aeeadaecbe48ecf0be44f030c06feb85e71eab
//
// Latest-code: `mapiterinit`
// https://github.com/golang/go/blob/cca23a73733ff166722c69359f0bb45e12ccaa2b/src/runtime/map.go#L798-L833
// ```
//     	// decide where to start
//    	r := uintptr(fastrand())                              //用于决定起点的r是个随机数
//    	if h.B > 31-bucketCntBits {
//    		r += uintptr(fastrand()) << 31
//    	}
//    	it.startBucket = r & bucketMask(h.B)                  //所以每次起点都是随机的
//    	it.offset = uint8(r >> h.B & (bucketCnt - 1))
//
//    	// iterator state
//    	it.bucket = it.startBucket
// ```
// map的本质是散列表，而map的增长扩容会导致重新进行散列，这就可能使map的遍历结果在扩容前后变得不可靠，
// Go设计者为了让大家不依赖遍历的顺序，故意在实现map遍历时加入了随机数，让每次遍历的起点--即起始bucket的位置不一样，
// 所以即使未扩容时我们遍历出来的map也总是无序的。

// map无序性的内部实现
func main() {

}
