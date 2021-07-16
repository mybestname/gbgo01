## 1245. 树的直径 
- https://leetcode-cn.com/problems/tree-diameter/
- https://leetcode.com/problems/tree-diameter/
- Given an undirected tree, return its diameter: 
   - the number of edges in a longest path in that tree.
- The tree is given as an array of edges where `edges[i] = [u, v]`
  is a bidirectional edge between nodes u and v.  
- Each node has labels in the set `{0, 1, ..., edges.length}`.




- 一棵「无向树」，请你测算并返回它的「直径」：这棵树上最长简径的边数。
- 树用一个由所有「边」组成的数组 edges 来表示
  - 其中 edges[i] = [u, v] 表示节点 u 和 v 之间的双向边。
- 树上的节点都已经用 set `{0, 1, ..., edges.length}` 中的数做了标记
  每个节点上的标记都是独一无二的。


```
示例 1
          0
         / \
        1  2

输入：edges = [[0,1],[0,2]]
输出：2
解释：
这棵树上最长的路径是 1 - 0 - 2，边数为 2。
```
```
示例 2：

                       1 
                    /  |  \
                   2   4   0
                 /    /
                3    5

输入：edges = [[0,1],[1,2],[2,3],[1,4],[4,5]]
输出：4
解释： 
这棵树上最长的路径是 3 - 2 - 1 - 4 - 5，边数为 4。
```

### 思路
- 两次深度优先遍历（DFS）
  - 第一次：从任意一点出发，找到距离它最远的点p
  - 第二次：从p点出发，找到距离它最远的点q
  - 则连接p，q两点的路径即为树的直径。

```

                       1 
                    /  |  \
                   2   4   0
                 /    /
                3    5

输入：edges = [[0,1],[1,2],[2,3],[1,4],[4,5]]

选择 1 到 0 ， 得到最远点为3，距离为3
1st dfs
dfs (1,0)
  P=0,Dist=[0,0,0,0,0,0],Dist[P]=0
dfs (2,1)
  P=2,Dist=[0,0,1,0,0,0],Dist[P]=1
dfs (3,2)
  P=3,Dist=[0,0,1,2,0,0],Dist[P]=2
dfs (4,1)
  P=3,Dist=[0,0,1,2,1,0],Dist[P]=2
dfs (5,4)
  P=3,Dist=[0,0,1,2,1,2],Dist[P]=2
  
再从 3 到 0， 得到最远点为5，距离为4。所有树的直径为4  
2nd dfs
dfs (3,0)
  P=3,Dist=[0,0,1,0,1,2],Dist[P]=0
dfs (2,3)
  P=2,Dist=[0,0,1,0,1,2],Dist[P]=1
dfs (1,2)
  P=1,Dist=[0,2,1,0,1,2],Dist[P]=2
dfs (0,1)
  P=0,Dist=[3,2,1,0,1,2],Dist[P]=3
dfs (4,1)
  P=0,Dist=[3,2,1,0,3,2],Dist[P]=3
dfs (5,4)
  P=5,Dist=[3,2,1,0,3,4],Dist[P]=4


选择1到4，最远点为3，距离为2
1st dfs
dfs (1,4)
  P=0,Dist=[0,0,0,0,0,0],Dist[P]=0
dfs (0,1)
  P=0,Dist=[1,0,0,0,0,0],Dist[P]=1
dfs (2,1)
  P=0,Dist=[1,0,1,0,0,0],Dist[P]=1
dfs (3,2)
  P=3,Dist=[1,0,1,2,0,0],Dist[P]=2
  
在从3到4，最远点为5，距离为4，所以直径是4  
2nd dfs
dfs (3,4)
  P=3,Dist=[1,0,1,0,0,0],Dist[P]=0
dfs (2,3)
  P=2,Dist=[1,0,1,0,0,0],Dist[P]=1
dfs (1,2)
  P=1,Dist=[1,2,1,0,0,0],Dist[P]=2
dfs (0,1)
  P=0,Dist=[3,2,1,0,0,0],Dist[P]=3
dfs (4,1)
  P=0,Dist=[3,2,1,0,3,0],Dist[P]=3
dfs (5,4)
  P=5,Dist=[3,2,1,0,3,4],Dist[P]=4


  
选择2到3，则最远点为5，距离为3                           对比清除整个Dist数组的情况
1st dfs                                     1st dfs
dfs (2,3)                                   dfs (2,3)
  P=0,Dist=[0,0,0,0,0,0],Dist[P]=0            P=0,Dist=[0,0,0,0,0,0],Dist[P]=0
dfs (1,2)                                   dfs (1,2)
  P=1,Dist=[0,1,0,0,0,0],Dist[P]=1            P=1,Dist=[0,1,0,0,0,0],Dist[P]=1
dfs (0,1)                                   dfs (0,1)
  P=0,Dist=[2,1,0,0,0,0],Dist[P]=2            P=0,Dist=[2,1,0,0,0,0],Dist[P]=2
dfs (4,1)                                   dfs (4,1)
  P=0,Dist=[2,1,0,0,2,0],Dist[P]=2            P=0,Dist=[2,1,0,0,2,0],Dist[P]=2
dfs (5,4)                                   dfs (5,4)
  P=5,Dist=[2,1,0,0,2,3],Dist[P]=3            P=5,Dist=[2,1,0,0,2,3],Dist[P]=3
                                              
再从5到3，则距离为4，所以直径为4                          
2nd dfs                                    2nd dfs                            
dfs (5,3)                                  dfs (5,3)                          
  P=5,Dist=[2,1,0,0,2,0],Dist[P]=0           P=5,Dist=[0,0,0,0,0,0],Dist[P]=0 
dfs (4,5)                                  dfs (4,5)                          
  P=4,Dist=[2,1,0,0,1,0],Dist[P]=1           P=4,Dist=[0,0,0,0,1,0],Dist[P]=1 
dfs (1,4)                                  dfs (1,4)                          
  P=1,Dist=[2,2,0,0,1,0],Dist[P]=2           P=1,Dist=[0,2,0,0,1,0],Dist[P]=2 
dfs (0,1)                                  dfs (0,1)                          
  P=0,Dist=[3,2,0,0,1,0],Dist[P]=3           P=0,Dist=[3,2,0,0,1,0],Dist[P]=3 
dfs (2,1)                                  dfs (2,1)                          
  P=0,Dist=[3,2,3,0,1,0],Dist[P]=3           P=0,Dist=[3,2,3,0,1,0],Dist[P]=3 
dfs (3,2)                                  dfs (3,2)                          
  P=3,Dist=[3,2,3,4,1,0],Dist[P]=4           P=3,Dist=[3,2,3,4,1,0],Dist[P]=4

```

## 543. 二叉树的直径
- https://leetcode-cn.com/problems/diameter-of-binary-tree/
- https://leetcode.com/problems/diameter-of-binary-tree/
- 给定一棵二叉树，你需要计算它的直径长度。
- 一棵二叉树的直径长度是任意两个结点路径长度中的最大值。
  - 这条路径可能穿过也可能不穿过根结点。

```
示例 :
给定二叉树

          1
         / \
        2   3
       / \     
      4   5    
返回 3, 它的长度是路径 [4,2,1,3] 或者 [5,2,1,3]。
```


