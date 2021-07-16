## 34. 在排序数组中查找元素的第一个和最后一个位置
- https://leetcode-cn.com/problems/find-first-and-last-position-of-element-in-sorted-array/
- https://leetcode-cn/problems/find-first-and-last-position-of-element-in-sorted-array/
- 给定一个按照升序排列的整数数组 nums，和一个目标值 target。
  找出给定目标值在数组中的开始位置和结束位置。
- 如果数组中不存在目标值 target，返回 [-1, -1]。
- 要求 
  - 时间复杂度为 O(log n)
```
示例 1：

输入：nums = [5,7,7,8,8,10], target = 8
输出：[3,4]
```
```
示例 2：

输入：nums = [5,7,7,8,8,10], target = 6
输出：[-1,-1]
```
```
示例 3：

输入：nums = [], target = 0
输出：[-1,-1]
```

### 思路
- 有重复值的情况
- 二分法求解，注意二分模版的记忆和使用
- 什么是开始位置？查询第一个 `>= target` ，即查询low_bound
- 什么是结束位置？查询最后一个 `<= target`，
- 然后判断两个位置的数据的合法性。
```c
      ------<=target------>=target------
      
```

