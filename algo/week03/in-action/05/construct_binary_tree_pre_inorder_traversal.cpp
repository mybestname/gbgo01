#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    TreeNode* buildTree(vector<int>& preorder, vector<int>& inorder) {
        // 结束条件

        // 3 is root from [3,9,20,15,7]
        // inorder [9,3,15,20,7]  root_index = 1, r_len=1, => [0:0], [1:1], [2:4]
        // split preorder [3,9,20,15,7],  => [0], [1], [2:4]
        //   [3]
        //   / \
        // [9] [20,15,7]
        // [9] [15,20,7]
        return _buildTree(preorder,inorder,0,preorder.size()-1, 0,inorder.size()-1);

    }
    // C++不支持slice，比较麻烦，建立一个辅助函数，传递下标
    // 模拟slice操作 即 preoder[l1:r1], inorder[l2:r2]
    TreeNode* _buildTree(vector<int>& preorder, vector<int>& inorder, int l1, int r1, int l2, int r2) {
        // 边界条件
        if (l1 > r1 ) return nullptr;
        TreeNode* root = new TreeNode(preorder[l1]);  //注意现在隐含的意思是传preoder[l1:r1], 所以是l1是per的第一个
        // 需要在inorder[l2:r2]中找root的位置
        int mid = l2;
        while(inorder[mid]!=root->val) mid++;
        // [9,3,15,20,7]  -> mid = 2 ,  left:[9], right:[15,20,7]
        // l2 mid     r2
        // [3,9,20,15,7]
        // l1(root)   r1
        int left_size = mid-l2;
        //int right_size = r2 -mid;
        // left : peroder[?:?] inorder[?:?]
        root->left = _buildTree(preorder,inorder,l1+1,l1+left_size, l2,mid-1);
        // right : preoder[?:?] inorder[?:?]
        root->right = _buildTree(preorder,inorder,l1+left_size+1, r1, mid+1, r2);
        return root;
    }

};

int main(){

    // 前序遍历 preorder = [3,9,20,15,7]
    // 中序遍历 inorder = [9,3,15,20,7]
    //                3
    //              / \
    //             9  20
    //               /  \
    //              15   7

    vector<int> preorder = {3,9,20,15,7};
    vector<int> inorder = {9,3,15,20,7};
    {
        Solution s;
        auto result = s.buildTree(preorder, inorder);
        cout << result << endl;
    }
    return 0;
}

