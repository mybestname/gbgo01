#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    TreeNode* invertTree(TreeNode* root) {
        // 结束处理
        if (root == nullptr) return nullptr;
        // 先左再右
        invertTree(root->left);
        invertTree(root->right);
        // 左右交换
        TreeNode* temp = root->left;
        root->left = root->right;
        root->right = temp;
        return root;
    }
};

int main() {
    // Input: root = [4,2,7,1,3,6,9]
    //               4
    //             /  \
    //           2     7
    //          / \   / \
    //        1   3  6   9
    // Output: [4,7,2,9,6,3,1]

    TreeNode* root = new TreeNode(4);
    root->left = new TreeNode(2);
    root->right = new TreeNode(7);
    root->left->left = new TreeNode(1);
    root->left->right = new TreeNode(3);
    root->right->left = new TreeNode(6);
    root->right->right = new TreeNode(9);

    cout << "input tree\n"<< root << endl;
    Solution s;
    cout << "invert tree\n"<< s.invertTree(root) << endl;
    return 0;
}
