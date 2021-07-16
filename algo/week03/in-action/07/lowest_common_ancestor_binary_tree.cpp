#include <iostream>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    TreeNode* lowestCommonAncestor(TreeNode* root, TreeNode* p, TreeNode* q) {
        father = unordered_map<int,TreeNode*>();
        redNode = unordered_set<int>();
        redNode.insert(root->val);
        calcFather(root);
        // build red nodes from p
        while(p->val!= root->val) {
            redNode.insert(p->val);
            p = father[p->val];
        }
        // find red node from q
        while (redNode.find(q->val) == redNode.end()) {
            q = father[q->val];
        }
        return q;
    }

private:
    // dfs求fater
    void calcFather(TreeNode* root){
        if (root == nullptr) return;
        if ( root->left != nullptr ) {
            father[root->left->val] = root;
            calcFather(root->left);
        }
        if ( root->right != nullptr ) {
            father[root->right->val] = root;
            calcFather(root->right);
        }
    }
    unordered_map<int,TreeNode*> father;
    unordered_set<int> redNode;
};

int main() {
    // 输入：root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 1
    //                         3
    //                     /    \
    //                    5       1
    //                  /  \    /   \
    //                 6   2    0   8
    //                    / \
    //                   7   4
    // 输出：3
    TreeNode* root = new TreeNode(3);
    TreeNode* p = new TreeNode(5);
    TreeNode* q = new TreeNode(1, new TreeNode(0), new TreeNode(8));
    root->left = p;
    root->right = q;
    p->left = new TreeNode(6);
    p->right = new TreeNode (2, new TreeNode(7), new TreeNode(4));
    cout << root << endl;
    Solution s;
    auto lca = s.lowestCommonAncestor(root, p, q);
    cout << lca << endl;
    return 0;
}


