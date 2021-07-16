## 17. 电话号码的字母组合
- https://leetcode-cn.com/problems/letter-combinations-of-a-phone-number/
- https://leetcode.com/problems/letter-combinations-of-a-phone-number/
- 给定一个仅包含数字 2-9 的字符串，返回所有它能表示的字母组合。答案可以按 **任意顺序** 返回。
- 给出数字到字母的映射如下（与电话按键相同）。
- 注意 1 不对应任何字母。
```c
    1     2     3
    4     5     6
    7     8     9
    *     0     #
    
2:abc 3:def  4:ghi 5:jkl
6:mno 7:pqrs 8:tuv 9:wxyz
```
```
示例 1：

输入：digits = "23"
输出：["ad","ae","af","bd","be","bf","cd","ce","cf"]
```
```
示例 2：

输入：digits = ""
输出：[]
```
```
示例 3：

输入：digits = "2"
输出：["a","b","c"]
```
### 思路
- DFS解决
