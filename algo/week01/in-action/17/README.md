## 84. 柱状图中最大的矩形
- https://leetcode-cn.com/problems/largest-rectangle-in-histogram/
- https://leetcode.com/problems/largest-rectangle-in-histogram/
- 给定 n 个非负整数，用来表示柱状图中各个柱子的高度。每个柱子彼此相邻，且宽度为 1 。
- 求在该柱状图中，能够勾勒出来的矩形的最大面积。
```
6             x--x
5          x--x  x
4          x  x  x
3          x  x  x  x--x
2    x--x  x  x  x--x  x
1    x  x--x  x  x  x  x
0    x  x  x  x  x  x  x
  0  1  2  3  4  5  6  7
柱状图的示例，其中每个柱子的宽度为 1，给定的高度为 [2,1,5,6,2,3]。

6             x--x
5          x*****x
4          x*****x
3          x*****x  x--x
2    x--x  x*****x--x  x
1    x  x--x*****x  x  x
0    x  x  x*****x  x  x

图中阴影部分为所能勾勒出的最大矩形面积，其面积为 10 个单位。
```

```
示例:

输入: [2,1,5,6,2,3]
输出: 10
```
### 思路
- 首先考虑暴力：
  - 依旧寻找i，j, i<j
    ```
    for i= 0 to n-1
      for j = i+1 to n;
           area = (j-i)* min(height[i]..height[j])
    ```
  - 和接水不一样，接水只和i，j相关，而矩形则要求i，j中间所以取值的min，所以不能简单使用双指针
  - 和一段相关，无法用双指针。（是只跟两个端点相关，还是和中间都相关）
- 简化情况，从简单到复杂，直到找到通用的解法。
  - 考虑柱子为单调递增的情况。
    ```c
     4             x--x
     3          x--x  x
     2       x--x  x  x
     1    x--x  x  x  x
     0    x  x  x  x  x
       0  1  2  3  4  5
    ```
    - 此时可能的答案为4种：
    ```
     x--x                                              
     x  x    x--x--x                                 
     x  x    x  x  x    x--x--x--x                 
     x  x    x  x  x    x  x  x  x    x--x--x--x--x
     x  x    x  x  x    x  x  x  x    x  x  x  x  x
    ```
    - 可以观察出，当高度递增时候，针对每一个高度，只有最长的矩形是一个可选值。
    - 即同高度，其它长度的矩形的可能性都是**冗余**的。
  - 考虑某个递增区间后减小的情况
    ```c
     4             x--x  
     3          x--x  x   
     2       x--x  x  x--x
     1    x--x  x  x  x  x
     0    x  x  x  x  x  x
       0  1  2  3  4  5  6
    ```   
     - 通过分析，可以发现
     ```c
     4             x--x                             
     3          x--x  x        对于最后一块儿来说，等效于求                   
     2       x--x  x  x--x     其实还是递增        x--x--x--x--x
     1    x--x  x  x  x  x                    x--x  x  x  x  x
     0    x  x  x  x  x  x                    x  x  x  x  x  x
       0  1  2  3  4  5  6                                   
     ```
     - 那么可以归纳为，首先求单调递增，当减小后，把减小的块儿转化为另外的递增形式。
- 考虑如何实现
  - 上述分析，只在末尾发生变化，可以考虑stack的解决方案。
  - 本题是单调栈的经典题。