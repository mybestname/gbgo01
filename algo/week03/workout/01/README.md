## 106. 从中序与后序遍历序列构造二叉树
- https://leetcode-cn.com/problems/construct-binary-tree-from-inorder-and-postorder-traversal/
- https://leetcode.com/problems/construct-binary-tree-from-inorder-and-postorder-traversal/
- 根据一棵树的中序遍历与后序遍历构造二叉树。
- 注意:
  - 你可以假设树中没有重复的元素。
    
```
例如，给出

中序遍历 inorder = [9,3,15,20,7]
后序遍历 postorder = [9,15,7,20,3]
返回如下的二叉树：
               3
              / \
             9  20
               /  \
              15   7
```

### 思路
- 和105思路一样
- postorder 后续，左右根： 可得根节点，有了根，那么根的左边，和右边就可以根据inorder得到焕发
- inorder 左根右，得到划分后，可以递归的求左子树和右子树。还是需要postoder帮助去确定根。

