## 433. 最小基因变化
- https://leetcode-cn.com/problems/minimum-genetic-mutation/
- https://leetcode.com/problems/minimum-genetic-mutation/
- 一条基因序列由一个带有8个字符的字符串表示，
  其中每个字符都属于 "A", "C", "G", "T"中的任意一个。
- 假设我们要调查一个基因序列的变化。一次基因变化意味着这个基因序列中的一个字符发生了变化。
  - 例如，基因序列由"AACCGGTT" 变化至 "AACCGGTA" 即发生了一次基因变化。
- 与此同时，每一次基因变化的结果，都需要是一个合法的基因串，即该结果属于一个基因库。
- 现在给定3个参数 — start, end, bank，
  - 分别代表起始基因序列，目标基因序列及基因库，
- 请找出能够使起始基因序列变化为目标基因序列所需的最少变化次数。
- 如果无法实现目标变化，请返回 -1。
- 注意：
  - 起始基因序列默认是合法的，但是它并不一定会出现在基因库中。
  - 如果一个起始基因序列需要多次变化，那么它每一次变化之后的基因序列都必须是合法的。
  - 假定起始基因序列与目标基因序列是不一样的。

```
示例 1：

start: "AACCGGTT"
end:   "AACCGGTA"
bank: ["AACCGGTA"]

返回值: 1
```
```
示例 2：

start: "AACCGGTT"
end:   "AAACGGTA"
bank: ["AACCGGTA", "AACCGCTA", "AAACGGTA"]

返回值: 2
```
```
示例 3：

start: "AAAAACCC"
end:   "AACCCCCC"
bank: ["AAAACCCC", "AAACCCCC", "AACCCCCC"]

返回值: 3
```

### 思路
- 最少变化次数。-> 求最小步数，
- 一般使用BFS, 第K步走完才走K+1步，在任意时刻队列中至多只有两层。 
  - 广搜队列的2段性和单调性，适合求层数，最小步数。
- A,C,G,T 只有四种选择，8个字符，4^8=2^16=65536种可能性。
- 由65536个点，每一个点有8*3=24条出边，构成的有向图

```
                AAAAAAAA
             /  /      \     \ ...      \   
     CAAAAAAA GAAAAAAA TAAAAAAA .... 24个出边。
```

如果要遍历整张图，那么O(2^16*8*3)

