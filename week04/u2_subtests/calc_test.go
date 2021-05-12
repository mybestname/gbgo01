package main

import "testing"

func TestMul(t *testing.T) {
	t.Run("pos", func(t *testing.T) {
		if Mul(2, 3) != 6 {
			t.Fatal("fail")
		}

	})
	t.Run("neg", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("fail")
		}
	})
}
// $ go test -run TestMul/pos -v
// === RUN   TestMul
// === RUN   TestMul/pos
// --- PASS: TestMul (0.00s)
//     --- PASS: TestMul/pos (0.00s)
// PASS
// ok      u2_subtests     0.053s

// table-driven tests
func TestMul2(t *testing.T) {
	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"pos", 2, 3, 6},
		{"neg", 2, -3, -6},
		{"zero", 2, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := Mul(c.A, c.B); ans != c.Expected {
				t.Fatalf("%d * %d expected %d, but %d got",
					c.A, c.B, c.Expected, ans)
			}
		})
	}
}

type calcCase struct{
	Name           string
	A, B, Expected int }

func createMulTestCase(t *testing.T, c *calcCase) {
	t.Helper()                                       //加了help错误返回在63行, 不见返回在54行，显然加help更方便定位错误。
	if ans := Mul(c.A, c.B); ans != c.Expected {
		t.Fatalf("%d * %d expected %d, but %d got",
			c.A, c.B, c.Expected, ans)
	}
}

// 对一些重复的逻辑，抽取出来作为公共的帮助函数(helpers)，可以增加测试代码的可读性和可维护性
func TestMul3(t *testing.T) {
	createMulTestCase(t, &calcCase{"pos", 2, 3, 6})
	createMulTestCase(t, &calcCase{"neg", 2, -3, -6})
	createMulTestCase(t, &calcCase{"zero",2, 0, 2}) // wrong case
}