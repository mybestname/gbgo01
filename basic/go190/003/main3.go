package main

func m2() {
	s1 := make([]int,5)
	for i := range s1 {
		s1[i] = i+1
	}
	s2 := []int{1,2,3,4,5}

	var a [5]int
	s3 := a[:0]
	_ = s1
	_ = s2
	_ = s3
}

//go:generate go tool compile -N -l main3.go
//go:generate go tool objdump main3.o
//go:generate rm main3.o
//GOSSAFUNC=m2 go build -gcflags "-N -l" main3.go

// 对于go来说，如果想初始化slice，同时指定默认值，只能靠slice literal
// 因为从make的实现来说，参数只是指定backend的初始数组长度（cap）以及slice的长度。
// 默认填充的是零值。
// 而ssa阶段的OpSliceMake操作，输入为：slice的元素类型、backend的对应数组的指针、slice大小，slice容量。
