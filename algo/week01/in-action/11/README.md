## 1248. 统计「优美子数组」
- https://leetcode-cn.com/problems/count-number-of-nice-subarrays/
- 给你一个整数数组 nums 和一个整数 k。
- 如果某个 **连续** 子数组中恰好有 k 个奇数数字，我们就认为这个子数组是「优美子数组」。
- 请返回这个数组中「优美子数组」的数目。
```
示例 1：
输入：nums = [1,1,2,1,1], k = 3
输出：2
解释：包含 3 个奇数的子数组是 [1,1,2,1] 和 [1,2,1,1] 。

示例 2：
输入：nums = [2,4,6], k = 1
输出：0
解释：数列中不包含任何奇数，所以不存在优美子数组。

示例 3：
输入：nums = [2,2,2,1,2,2,1,2,2,2], k = 2
输出：16
```

### 思路
- 理解什么是优美，子数组中恰好有k个奇数
- 理解什么是连续，子数组是连续的，即顺序不能改变。  
  - 例如[1,1,2,1,1]，如果k=3，那么子数组长度需要>=3，<5，也就是4或者3。
  - 长度为3的情况下，不存在有三个奇数的**连续**子数组。
  - 长度为4的**连续**子数组有两个：[1,1,2,1],[1,2,1,1]，这两个都满足条件。
  - 所以例1的结果为2。
``` 
  5+4+3+2+1 = n(n+1)/2 = 5*6/2 = 15
  sub = n(n+1)/2 - 1 = (n(n+1)-2)/2 =(n^2 +n -2) / 2 = (25 + 5-2) /2 = 28/2 =  14  
```
- 首先对原数组进行处理，先把原数据变成`0`和`1`，只保留问题的实质：奇偶属性。
```
[1, 1, 2, 1, 1 ]
[1, 1, 0 ,1, 1 ]   1表示奇数，0表示偶数

[2, 4, 6 ]
[0, 0, 0 ]
```
- 再次等效变换，求奇数个数变成求 **和==k** 
```
[1, 1, 0, 1, 1] 有多少个1，变为 和==k
```
- 最终变成求和问题，这样和本题的考点：**前缀和** 挂钩了。
  - `子数组`其实就是`子段`，即求`子段和`的问题
- 简单思想是暴力解法
  ```
  for r in (1..n) 
     for l in (1..r)
        if (s[r] - s [l-1] == k) 
          ans += 1
  ```
- 凡是遇到两重循环的暴力解法，都可以通过：
  - 把两个循环变量分离：
    - 把r固定，观察内层循环在做什么？
      - 固定外层循环变量，考虑内层循环需要满足什么条件。
    - 条件总结为：
      - 对于每一个r（1..n），考虑有几个l（1..r）使得 `s[r] - s [l-1] == k`
    - 改写为条件：
      - 对于每一个i（1..n），考虑有几个j（0..i-1）使得 `s[i] - s [j] == k`
      - 对于每一个i（1..n），考虑有几个j（0..i-1）使得 `s [j] == s[i] - k `
      - 对于每一个i，有几个`s[j]`等于`s[i]-k`
    - 变为：在一个数组中统计"等于某个数"的数的数量。  
      - 使用一个count数组，来记录`等于某个数`的数的数量（等于一个map）
      - count[key]=val
        - key = `s[j]`，对每一个j，count[s[j]]++;
        - 问题等于求count[s[i]-k]=？