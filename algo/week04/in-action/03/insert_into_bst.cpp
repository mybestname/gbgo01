#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"
using namespace std;

class Solution {
public:

    TreeNode* insertIntoBST(TreeNode* root, int val) {
        // 边界条件
        //  按题意不用考虑重复，所以边界为 root == nullptr
        if (root == nullptr) {
            //达到插入点，原地插入
            return new TreeNode(val);
        }
        // 中心思路：大的往左，小的往右
        if (val < root->val) {
            root->left = insertIntoBST(root->left,val);
        }else {
            root->right = insertIntoBST(root->right,val);
        }
        return root;
    }
};

int main() {
    struct Test {
        TreeNode* root;
        int val;
        TreeNode* expect ;

    };
    vector<Test> tests = {
            {
                    .root   =  new TreeNode(4, new TreeNode(2,new TreeNode(1), new TreeNode(3)), new TreeNode(7)),
                    .val      = 5,
                    .expect =  new TreeNode(4, new TreeNode(2,new TreeNode(1), new TreeNode(3)), new TreeNode(7, new TreeNode(5), nullptr))
            },
//             4                          4
//          /   \                      /   \
//        2      7                   2      7
//     /   \                      /   \    /
//    1    3                     1    3   5

    };

    Solution s;
    for (auto &test : tests) {
        cout << "root=\n" << test.root << "val=" << test.val
                << "\nexpect=\n" << test.expect << endl;
        auto result = s.insertIntoBST(test.root,test.val);
        cout << "got=\n" << result << endl;
    }
    return 0;
}

