## 25. K个一组翻转链表
- https://leetcode-cn.com/problems/reverse-nodes-in-k-group/
- 给你一个链表，每 k 个节点一组进行翻转，请你返回翻转后的链表。
- k 是一个正整数，它的值小于或等于链表的长度。
- 如果节点总数不是 k 的整数倍，那么请将最后剩余的节点保持原有顺序。
- 进阶：
   - 你可以设计一个只使用常数额外空间的算法来解决此问题吗？
   - 你不能只是单纯的改变节点内部的值，而是需要实际进行节点交换。

```
示例 1：

输入：head = [1,2,3,4,5], k = 2
输出：[2,1,4,3,5]
```
```
示例 2：

输入：head = [1,2,3,4,5], k = 3
输出：[3,2,1,4,5]
```
```
示例 3：

输入：head = [1,2,3,4,5], k = 1
输出：[1,2,3,4,5]
```
```
示例 4：

输入：head = [1], k = 1
输出：[1]
```

### 思路
- 仔细理解题意：注意是每k组，所以有可能有多个k组的情况：
  ```
  [1,2,3,4,5] 2 
  有两个k组  [1,2] 和 [3,4] 
  [2,1,4,3,5] 
  ```
- 1、首先要找到所有分组（找到每一组的开始和结尾）
     - 需要找头和结尾吗？
     - 上一个结尾，就是下一个头，那么reduce找结尾
     - 单独建立函数
- 2、然后再翻转每个分组（这部分可以复用例题4翻转链表的内容）
    - 第二步中间的部分和例题4一样，基本可以复用该函数
      - 第一点：
        - 但是以前只要传一个参数：头部，现在需要修改为接受2个参数。头部和尾部。
        - 因为**现在的结尾是尾部节点**，而不是整个的链表的尾。 
          - 对于尾部需要特殊处理
          - 注意`while(head!=nil)`和 `while(head!=end)`
            的不同
          - `head==nil`尾部是已经反转的，而`head==end`时候，end并没有反转。
      - 第二点：
        - 例题4是改n条边，这里是改n-1条边，头部节点不需要修改
        - 头部还要另外处理。
      - 第三点：
        - 不需要返回。
    - 翻转完成后，**每个k组的头尾部分**还需要进一步处理
      - 头部
      - 尾部
- 3、特别注意是，反转函数内部和分组部分都需要处理头部和尾部。
     - 因为两种的视角不一样。
       - 反转函数内部要保证，头和尾都顺利建立反转关系。
       - 分组的头和尾，需要重新建立组和组之间的连续。
       - 两种关系都需要处理。