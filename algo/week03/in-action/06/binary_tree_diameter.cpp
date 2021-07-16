#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    int diameterOfBinaryTree(TreeNode* root) {
        diameter = 0;
        dfsBinaryTree(root);
        return diameter;
    }
private:
    int dfsBinaryTree(TreeNode* root) {
        //边界条件
        if (root == nullptr) return 0;
        int l = dfsBinaryTree(root->left);
        int r = dfsBinaryTree(root->right);
        diameter = max(diameter,l+r);
        return 1+ max(l,r);
    }
    int diameter;
};
int main(){
        TreeNode* root = new TreeNode(1);
        root->left = new TreeNode(2, new TreeNode(4), new TreeNode(5));
        root->right = new TreeNode(3);
        //           1
        //         / \
        //        2   3
        //       / \
        //      4   5
        cout << "\ninput binary tree is" << endl;
        cout << root << endl;
        Solution s;
        auto d = s.diameterOfBinaryTree(root);
        cout << "diameter of binary tree is " << d << endl;
    return 0;
}

