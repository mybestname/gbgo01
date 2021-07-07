## 50. Pow(x, n)
- https://leetcode-cn.com/problems/powx-n/
- https://leetcode.com/problems/powx-n/
- 实现 pow(x, n) ，即计算 x 的 n 次幂函数（即，xn）。

```
示例 1：

输入：x = 2.00000, n = 10
输出：1024.00000
```
```
示例 2：

输入：x = 2.10000, n = 3
输出：9.26100
```
```
示例 3：

输入：x = 2.00000, n = -2
输出：0.25000
解释：2-2 = 1/22 = 1/4 = 0.25
```
### 思路
- 使用分治方法
- n为偶数: Pow(x,n) = Pow(x,n/2)*Pow(x,n/2)
  - 注意：可以使用temp变量，这样不用真的算两次Pow(x,n/2)
- n为奇数: Pow(x,n) = Pow(x,n/2)*Pow(x,n/2)*x  （注：这里n/2已经自动取整）
- 递归边界：Pow(x,1) = x
- 对于负数：1/Pow(x,-n);  
- 算法复杂度：O(log(n)) (基于不用算两次Pow(x,n/2)的情况)


