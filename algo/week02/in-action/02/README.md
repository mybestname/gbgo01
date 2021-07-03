## 49. 字母异位词分组
- https://leetcode-cn.com/problems/group-anagrams/
- https://leetcode.com/problems/group-anagrams/
- 给定一个字符串数组，将字母异位词组合在一起。字母异位词指**字母相同，但排列不同**的字符串。

```
示例:
输入: ["eat", "tea", "tan", "ate", "nat", "bat"]
输出:
[
["ate","eat","tea"],
["nat","tan"],
["bat"]
]
```
- 说明：
  - 所有输入均为小写字母。
  - 不考虑答案输出的顺序。

### 思路1
- 先对字符排序，建立key，同key即为一组。
- 例如
  - eat -> aet
  - tea -> aet
  - ate -> aet

### 思路2
- 统计每个字符串中字符出现次数
- 使用一个计算数组 int[26] 