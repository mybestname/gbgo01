#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"
using namespace std;

// 使用递归求解
class Solution {
public:
    TreeNode* deleteNode(TreeNode* root, int key) {
        if (root == nullptr) return nullptr; //找不到，直接返回空
        if (root->val == key) { //如果可以找到
            // 判断子树的情况
            // 如果只有一棵子树, 直接删除val，再把子树和父节点相连即可。
            if (root-> left == nullptr) return root->right;
            if (root-> right == nullptr) return root->left;
            // 两棵都有的情况
            // 需要找到后继，先删除后继，然后再用后继节点代替
            TreeNode* next = root->right;  //后继看右子树，这是有右子树的常见
            while(next->left != nullptr) next = next->left; // 一直下到左底
            // 删除后继，同时替换
            //            5                       5
            //           / \                     / \
            //          3   6                   4   6
            //         / \   \                 /     \
            //        2   4   7               2       7
            root->right = deleteNode(root->right, next->val);  // 删除后继4。
            // 替换为后继4
            root->val = next->val;
            return root;
        }
        // 首先递归找key，使用递归模版
        if (key < root->val) {
            root->left = deleteNode(root->left, key);
        }else {
            root->right = deleteNode(root->right, key);
        }
        return root;
    }
};

int main() {
    struct Test {
        TreeNode* root;
        int key;
    };
    auto* node3 = new TreeNode(3,new TreeNode(2), new TreeNode(4));
    auto* node6 = new TreeNode(6, nullptr, new TreeNode(7));
    vector<Test> tests = {
            {
                    .root   =  new TreeNode(5, node3, node6),
                    .key    =  3,
            },
    };
    // root = [5,3,6,2,4,null,7]
    // key = 3
    //
    //        5
    //       / \
    //      3   6
    //     / \   \
    //    2   4   7

    Solution s;
    for (auto &test : tests) {
        cout << "root\n" << test.root << "key=" << test.key << endl;
        auto result = s.deleteNode(test.root,test.key);
        cout << "got=\n" << result << endl;
    }
    return 0;
}

