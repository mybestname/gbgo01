#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    // 左子树的最大节点比root小
    // 右子树的最小节点比root大
    // 所有左子树和右子树自身也是二叉搜索树。
    bool isValidBST(TreeNode* root) {
        return check(root).isValid;
    }
private:
    struct subTreeInfo {
        bool isValid;
        int minVal;
        int maxVal;
    };
    subTreeInfo check(TreeNode* root){
        // 终止条件
        if (root == nullptr) {
            subTreeInfo info;
            info.isValid = true;
            info.minVal = numeric_limits<int>::max();   //int_MAX
            info.maxVal = numeric_limits<int>::min();   //int_MIN
            return info;
        }
        subTreeInfo left = check(root->left);
        subTreeInfo right = check(root->right);
        subTreeInfo result;
        result.isValid = left.isValid && right.isValid && left.maxVal < root->val && right.minVal > root->val;
        result.minVal = min(min(left.minVal,right.minVal), root->val);
        result.maxVal = max(max(left.maxVal, right.maxVal), root->val);
        // 没有任何全局变量，所以不用恢复。
        return result;
    };
};

int main() {
    // 输入:
    //             5
    //            / \
    //           1   4
    //              / \
    //             3   6
    // 输出: false
    // 解释: 输入为: [5,1,4,null,null,3,6]。
    TreeNode* root = new TreeNode(5);
    root->left = new TreeNode(1);
    root->right = new TreeNode(4);
    root->right->left = new TreeNode(3);
    root->right->right = new TreeNode(6);
    cout << "input tree\n"<< root << endl;
    Solution s;
    auto result = s.isValidBST(root);
    cout << "is_validate_BST:" << result << endl;
    return 0;
}

