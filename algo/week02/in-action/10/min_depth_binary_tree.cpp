#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    int minDepth(TreeNode* root) {
        // 停止条件/边界处理
        if (root == nullptr) {
            return 0;
        }
        // 注意：求minDepth的额外终止条件
        if (root->right == nullptr) {
            return minDepth(root->left)+1;
        }
        if (root->left == nullptr) {
            return minDepth(root->right)+1;
        }
        // 因为没有改任何共享的全局变量，所以没有状态恢复
        return min(minDepth(root->right), minDepth(root->left))+1;
    }
};
int main(){
    // 二叉树 [3,9,20,null,null,15,7]
    TreeNode* root = new TreeNode(3);
    root->left = new TreeNode(9);
    root->right = new TreeNode(20);
    root->right->left = new TreeNode(15);
    root->right->right = new TreeNode(7);
    cout << "input tree\n"<< root << endl;
    {
        Solution s;
        auto depth = s.minDepth(root);
        cout << "the minimum depth:" << depth << endl;
    }
    root = new TreeNode(2);
    root->right = new TreeNode(3);
    root->right->right = new TreeNode(4);
    root->right->right->right = new TreeNode(5);
    root->right->right->right->right = new TreeNode(6);
    cout << "input tree\n"<< root << endl;
    {
        Solution s;
        auto depth = s.minDepth(root);
        cout << "the minimum depth:" << depth << endl;
    }
    return 0;
}
