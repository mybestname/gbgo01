## 697. 数组的度
- https://leetcode-cn.com/problems/degree-of-an-array/
- https://leetcode.com/problems/degree-of-an-array/
- 给定一个非空且**只包含非负数**的整数数组 nums，
- 数组的度的定义是指**数组里任一元素出现频数的最大值**。
- 在 nums 中找到与 nums 拥有相同大小的度的**最短连续子数组**，返回其长度。


```
示例 1：

输入：[1, 2, 2, 3, 1]
输出：2
解释：
输入数组的度是2，因为元素1和2的出现频数最大，均为2.
连续子数组里面拥有相同度的有如下所示:
[1, 2, 2, 3, 1], [1, 2, 2, 3], [2, 2, 3, 1], [1, 2, 2], [2, 2, 3], [2, 2]
最短连续子数组[2, 2]的长度为2，所以返回2.
```
```
示例 2：

输入：[1,2,2,3,1,4,2]
输出：6
```
### 思路
- 首先统计数组各个元素出现频数。
- 最短的连续子数组代表什么
- 对任何一个元素记录如下数据，使用map维护：
  1。频次，2。第一次出现的index，3。最后一次出现的index
```  
  [1, 2, 2, 3, 1]
    index 0 -> 1 -> map k=1, v= (1,0,0)
    index 1 -> 2 -> map k=2, v= (1,1,1)
    index 2 -> 2 -> map k=2, v= (2,1,2)
    index 3 -> 3 -> map k=3, v= (1,3,3)
    index 4 -> 1 -> map k=1, v= (2,0,4)
  for index,num range numbs
    if exist map[num]
      map[num][0]++
      map[num][2]=index   
    if noexist
      map[num] = (index,index,index)
    max_frequency = max(max_frequency, map(num)[0])
  for num in map:  
  if max_frequency == map(num)[0]
      max_length = max(max_length, map(num)
  return max_length    
```
