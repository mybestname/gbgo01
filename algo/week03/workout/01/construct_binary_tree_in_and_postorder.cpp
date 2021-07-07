#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:

    TreeNode* buildTree(vector<int>& inorder, vector<int>& postorder) {
        unordered_map<int,int> inOrderMap;
        for (int i = 0; i< inorder.size(); i++) {
            inOrderMap[inorder[i]]=i;
        }
        TreeNode* root = buildTreeNode(inorder,postorder,0, postorder.size()-1, 0,inOrderMap);
        return root;
    }
private:
    TreeNode* buildTreeNode(vector<int>& inorder, vector<int>& postorder, int postStart, int postEnd, int inStart, unordered_map<int,int> inorderMap) {
        // 边界条件
        if (postStart > postEnd) return nullptr;
        int rootVal = postorder[postEnd];
        TreeNode* root = new TreeNode(rootVal);
        int rootPos = inorderMap[rootVal];
        int len = rootPos-inStart;
        TreeNode* tree_l = buildTreeNode(inorder, postorder, postStart, postStart+len-1, inStart, inorderMap);
        TreeNode* tree_r = buildTreeNode(inorder, postorder,postStart+len, postEnd-1, rootPos+1, inorderMap);
        root->left = tree_l;
        root->right = tree_r;
        return root;
    }

};

int main(){
    // 中序遍历 inorder = [9,3,15,20,7]
    // 后序遍历 postorder = [9,15,7,20,3]
    //                3
    //              / \
    //             9  20
    //               /  \
    //              15   7

    vector<int> inorder = {9,3,15,20,7};
    vector<int> postorder = {9,15,7,20,3};
    {
        Solution s;
        auto result = s.buildTree(inorder,postorder);
        cout << result << endl;
    }
    return 0;
}