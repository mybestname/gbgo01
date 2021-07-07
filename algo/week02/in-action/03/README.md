## 30. 串联所有单词的子串
- https://leetcode-cn.com/problems/substring-with-concatenation-of-all-words/
- https://leetcode.com/problems/substring-with-concatenation-of-all-words/
- 给定一个字符串 s 和一些 **长度相同** 的单词 words 。
- 找出 s 中恰好可以由 words 中所有单词串联形成的子串的起始位置。
- 注意子串要与 words 中的单词完全匹配，**中间不能有其他字符**，
  但不需要考虑 words 中单词串联的顺序。
```
示例 1：

输入：s = "barfoothefoobarman", words = ["foo","bar"]
输出：[0,9]
解释：
从索引 0 和 9 开始的子串分别是 "barfoo" 和 "foobar" 。
输出的顺序不重要, [9,0] 也是有效答案。
```
```
示例 2：

输入：s = "wordgoodgoodgoodbestword", words = ["word","good","best","word"]
输出：[]
```
```
示例 3：

输入：s = "barfoofoobarthefoobarman", words = ["bar","foo","the"]
输出：[6,9,12]
```
### 思路
- 思考🤔 words需要考虑全排列的可能性，那么是一个n！种排列方法的问题。
  - 如果考虑words的每一种可能性，显然太慢。
- 当想不到解法时候，可以仔细考虑判定：
  - 给出一个s的子串，判定该子串是不是words的串联
- 暴力的解法
  - 判断一个字符串t，是否由words构成，
    - 把t分解成若干单词，看跟words数组是否相同（顺序无关，进行排序即可）
  - 字符串排序慢，如何改进？
- 思考上一道题，顺序无所谓，字符一致即可。即对所有字符排序建立key
  - 本道题，顺序无所谓，单词一致就可以。
- 那么现在本题变为
  - words是一个hash表 {"foo":1,"bar":1}
  - 子串为另一个hash表，按单词长度划词。
    {"bar":1,"foo":1}, {"arf":1,"oot":1,},...,{"foo":1,"bar":1},...{'oob':1,'arm':1}
  - 找到`字串hash表` == `wordshash表` 的情况。
  - 或者两个集合相等的情况
  - 这样不用对字符串排序。
- 进一步考虑如何设计hash表
  - hash表存储每一个单词出现的次数。

### 优化
- 没有必要一个字符一个字符的移动，可以一个word一个word的移动。
- 即`barfoothefoobarman`的子串不用是
  - `barfoo`,`arfoot`,`rfooth`,...这样逐字符构造子串。
  - 可以
  `barfoo`,`foothe`,`thefoo`,...这样逐word的构造字串。
- 对于任何一个字串，已知word长度的情况下（例如4），那么只有如下四种分隔方法。
- lingmindraboofooowingdingbarrwingmonkeypoundcake  
  - ling"'mind'rabo'ofoo'owin'"gdingbarrwingmonkeypoundcake  
  - l'ingm'indr'aboofooowingdingbarrwingmonkeypoundcake  
  - li'ngmi'ndra'boofooowingdingbarrwingmonkeypoundcake  
  - ling'mind'raboofooowingdingbarrwingmonkeypoundcake  