### 206. 反转链表
- https://leetcode-cn.com/problems/reverse-linked-list/
- 给你单链表的头节点 head ，请你反转链表，并返回反转后的链表。

```
示例 1：

输入：head = [1,2,3,4,5]
输出：[5,4,3,2,1]
```

```
示例 2：

输入：head = [1,2]
输出：[2,1]
```

```
示例 3：

输入：head = []
输出：[] 
```

- 思路
  - 首先回答问题：需要改多少个边？`[1,2,3,4,5]`为例子。
    ```
           1 -> 2 -> 3 -> 4 -> 5 -> nil
    =>
    
    nil <- 1 <- 2 <- 3 <- 4 <- 5 
    ```
    需要改5条边，也就是**需要改`n`条边！而不是`n-1`条!**
    要特别注意`nil`节点。
  - 然后回答如何改：
    ```
           node0.next -> node1
                         node1.next -> node2
                                       node2.next ->
    
           =>
    
           node2.next -> node1 
                         node1.next -> node0
    ```
    通过观察，可以看到，必须有一个暂存的节点。否则无法直接交换。
    所以我们需要多一个变量，来临时存放last节点
    