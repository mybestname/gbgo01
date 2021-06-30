## 155. 最小栈
- https://leetcode-cn.com/problems/min-stack/
- 设计一个支持 push ，pop ，top 操作，并能在常数时间内检索到最小元素的栈。
  - push(x) —— 将元素 x 推入栈中。
  - pop() —— 删除栈顶的元素。
  - top() —— 获取栈顶元素。
  - getMin() —— 检索栈中的最小元素。

```
示例:

输入：
["MinStack","push","push","push","getMin","pop","top","getMin"]
[[],[-2],[0],[-3],[],[],[],[]]

输出：
[null,null,null,null,-3,null,0,-2]

解释：
MinStack minStack = new MinStack();
minStack.push(-2);
minStack.push(0);
minStack.push(-3);
minStack.getMin();   --> 返回 -3.
minStack.pop();
minStack.top();      --> 返回 0.
minStack.getMin();   --> 返回 -2.
```

### 思路
- 前缀最小值
  - 多开一个栈，存储前缀最小值的stack，如果要getMin，那么就返回最小值。
  - 而主stack.pop()时候，也一起pop

```
[5, 3, 1, 4, 2 ] 主栈
[5, 3, 1, 1, 1 ] 前缀最小值栈 

返回min的时候走最小栈，pop时候，一起pop

push时候，主栈入值，最小值栈比较，入最小值。

如push 5

[5, 3, 1, 4, 2, 5 ]
[5, 3, 1, 1, 1, 1 ]

如push 0

[5, 3, 1, 4, 2, 5, 0 ]
[5, 3, 1, 1, 1, 1, 0 ]
```



