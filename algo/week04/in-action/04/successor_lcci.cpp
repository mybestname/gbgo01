#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"
using namespace std;

class Solution {
public:
    TreeNode* inorderSuccessor(TreeNode* root, TreeNode* p) {
        return findSucc(root, p->val);
    }

private:
    TreeNode* findSucc(TreeNode* root, int val) {
        TreeNode* cur = root;
        TreeNode* ans = nullptr;
        while (cur != nullptr) {
            // 后继需要记录一下经过的点，把进过点中的最小值点记录下来。
            // 这个经过最小值在没有找到或没有右子树的情况下，就是后继点。
            if (cur->val > val) {
                if (ans == nullptr || ans->val > cur->val)  ans = cur;
            }
            if (cur->val == val) { //找到val
                //看一下右子树
                if (cur->right != nullptr) { // 如果有右子树
                    //右子树的向左到底为后继
                    cur = cur->right;
                    while (cur->left != nullptr) cur = cur->left;
                    return cur;  //直接返回到底的结果为后继
                }
                break;  // 直接break即可，因为没有右子树且cur.val == val 的后继，一定在前面的ans中了。
            }
            // 遍历下一个
            if (val < cur->val ) cur = cur->left;
            else cur= cur->right;

        }
        return ans;
    }
};
int main() {
    struct Test {
        TreeNode* root;
        TreeNode* p;
        TreeNode* expect ;

    };
    auto* node2 = new TreeNode(2,new TreeNode(1), new TreeNode(3));
    auto* node7 = new TreeNode(7);
    vector<Test> tests = {
            {
                    .root   =  new TreeNode(4, node2, node7),
                    .p      =  new TreeNode(2),
                    .expect =  new TreeNode(3),
            },
//             4
//          /   \
//        2      7
//     /   \
//    1    3

            {
                    .root   =  new TreeNode(4, node2, node7),
                    .p      =  new TreeNode(4),
                    .expect =  new TreeNode(7),
            },
            {
                    .root = new TreeNode(2,new TreeNode(1), new TreeNode(3)),
                    .p = new TreeNode(1),
                    .expect = new TreeNode(2),
            }

    };

    Solution s;
    for (auto &test : tests) {
        cout << "root\n" << test.root << "p=" << test.p->val
             << "\nexpect\n" << test.expect << endl;
        auto result = s.inorderSuccessor(test.root,test.p);
        cout << "got=\n" << result << endl;
    }
    return 0;
}
