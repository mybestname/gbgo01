#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    int maxDepth(TreeNode* root) {
        // 停止条件/边界处理
        if (root == nullptr) {
           return 0;
        }
        // 因为没有改任何共享的全局变量，所以没有状态恢复
        return max(maxDepth(root->right), maxDepth(root->left)) + 1;
    }
};

// 思路2写法1：
class Solution2_1 {
public:
    int maxDepth(TreeNode* root) {
        ans = 0;
        calcDepth(root, 1);
        return ans;
    }
    void calcDepth(TreeNode* root, int depth) {
        // 终止条件/边界处理
        if(root == nullptr) {
            return;
        }
        ans = max(ans,depth); //答案的最后更新结果即为最后答案，因为总是记录最大值，所以不需要恢复。
        calcDepth(root->left, depth+1);
        calcDepth(root->right, depth+1);
        // 因为depth本身不是共享的，而是做为参数传递，所以不需要恢复。
    }

private:
    int ans = 0;
};

// 思路2写法2：
class Solution2_2 {
public:
    int maxDepth(TreeNode* root) {
        depth = 1;
        ans = 0;
        calcDepth(root);
        return ans;
    }
    void calcDepth(TreeNode* root) {
        // 终止条件/边界处理
        if(root == nullptr) {
            return;
        }
        ans = max(ans,depth);
        depth++; //递归前
        calcDepth(root->left);
        // depth--; 递归后还原     ---
        //                         + 这两句合并等于被省略。
        // depth++; 递归前        ---
        calcDepth(root->right);
        depth--; // 递归后还原
    }

private:
    int depth = 0;
    int ans = 0;
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
        auto depth = s.maxDepth(root);
        cout << "the maximum depth:" << depth << endl;
    }
    {
        Solution2_1 s2_1;
        auto depth = s2_1.maxDepth(root);
        cout << "the maximum depth:" << depth << endl;
    }
    {
        Solution2_2 s2_2;
        auto depth = s2_2.maxDepth(root);
        cout << "the maximum depth:" << depth << endl;
    }
    return 0;
}
