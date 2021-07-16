## 105. 从前序与中序遍历序列构造二叉树
- https://leetcode-cn.com/problems/construct-binary-tree-from-preorder-and-inorder-traversal/
- https://leetcode.com/problems/construct-binary-tree-from-preorder-and-inorder-traversal/
- 根据一棵树的前序遍历与中序遍历构造二叉树。
- 注意:
  - 你可以假设树中没有重复的元素。
```
例如，给出

前序遍历 preorder = [3,9,20,15,7]
中序遍历 inorder = [9,3,15,20,7]
返回如下的二叉树：

    3
   / \
   9  20
     /  \
    15   7
```

### 思路
- 核心思路：通过找到根在中序中的位置来确定左右子树的大小，使得递归可以进行。
- [3,9,20,15,7] 告诉我们，3是root,但是左子树和右子树的长度未知。
- [9,3,15,20,7]，因为3是根，所以左子树1个节点，右子树3个节点。
- 综合：我们知道
- [3 | 9 | 20,15,7]
- [9 | 3 | 15 20 7]
