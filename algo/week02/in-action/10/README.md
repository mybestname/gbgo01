## 104. 二叉树的最大深度
- 给定一个二叉树，找出其最大深度。
- 二叉树的深度为根节点到最远叶子节点的最长路径上的节点数。
- 说明: 叶子节点是指没有子节点的节点。
- https://leetcode-cn.com/problems/maximum-depth-of-binary-tree/
- https://leetcode.com/problems/maximum-depth-of-binary-tree/

```
示例：
给定二叉树 [3,9,20,null,null,15,7]，

    3
   / \
   9  20
     /  \
    15   7
返回它的最大深度 3 。
```

### 思路1
- 最大深度其实和上题类似，都是要比较min/max
- 每一个节点有一个深度属性，最大深度为max（左深度，右深度）然后+1
- 这种是自底向上：统计的思路，使用分治的思想。

### 思路2
- 将"深度"做为全局变量：跟随节点移动而动态变化。
- 递归一层，变量+1，在叶子处更新答案
- 需要注意恢复状态（保护/还原现场）（因为有全局共享变量的存在）


## 111. 二叉树的最小深度
- https://leetcode-cn.com/problems/minimum-depth-of-binary-tree/
- https://leetcode.com/problems/minimum-depth-of-binary-tree/
- 给定一个二叉树，找出其最小深度。
- 最小深度是从根节点到最近叶子节点的最短路径上的节点数量。
- 说明：叶子节点是指没有子节点的节点。

```
示例 1：
                    3
                   / \
                  9  20
                     / \
                    15 7 

输入：root = [3,9,20,null,null,15,7]
输出：2
```

```
示例 2：
                     2
                      \
                       3
                        \
                         4
                          \
                           5
                            \
                             6 
输入：root = [2,null,3,null,4,null,5,null,6]
输出：5
```

### 思路
- 和求Max一样
- 但是需要额外注意终止条件。


