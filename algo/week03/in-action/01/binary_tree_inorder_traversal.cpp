#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    // 中序：左 中 右
    vector<int> inorderTraversal(TreeNode* root) {
        // 边界条件
        if (root == nullptr) {
            return {};
        }
        vector<int> result;
        auto result_l = inorderTraversal(root->left);
        for (auto r_l : result_l) {
            result.push_back(r_l);
        }
        result.push_back(root->val);
        auto result_r = inorderTraversal(root->right);
        for (auto r_r : result_r) {
            result.push_back(r_r);
        }
        return result;
    }
};
// 解法2: 使用全局共享变量存放返回值
class Solution2 {
public:
    vector<int> inorderTraversal(TreeNode* root) {
        seq = {};
        find(root);
        return seq;
    }
private:
    vector<int> seq;
    // 中序：左 中 右
    void find(TreeNode* root) {
        // 边界条件
        if (root== nullptr) return;
        find(root->left);
        seq.push_back(root->val);
        find(root->right);
        // 注意，因为只是把结果按顺序放入返回列表，这个状态不参与到递归运算，所以不需要恢复状态。
    }
};

int main(){
    // 输入：root = [1,null,2,3]
    // 输出：[1,3,2]
    TreeNode* root = new TreeNode(1);
    root->right = new TreeNode(2);
    root->right->left = new TreeNode(3);
    cout << "input tree\n"<< root << endl;
    {
        Solution s;
        auto inorder = s.inorderTraversal(root);
        cout << "in order " << inorder << endl;
    }

    {
        Solution2 s;
        auto inorder = s.inorderTraversal(root);
        cout << "in order " << inorder << endl;
    }
    return 0;
}

