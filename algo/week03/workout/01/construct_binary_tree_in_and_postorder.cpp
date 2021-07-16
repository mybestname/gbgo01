#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;

// 清晰版的思路
class Solution {
public:

    TreeNode *buildTree(vector<int> &inorder, vector<int> &postorder) {
        return _buildTree(postorder,inorder,0,postorder.size()-1, 0,inorder.size()-1);
    }

private:
    TreeNode* _buildTree(vector<int>& postorder, vector<int>& inorder, int l1, int r1, int l2, int r2) {
        // 边界条件
        if (l1 > r1 ) return nullptr;
        TreeNode* root = new TreeNode(postorder[r1]);  //注意现在隐含的意思是传postoder[l1:r1], 所以是r1是post的最后一个
        // vector<int> idx = {l1,r1,l2,r2};
        // cout <<  "root=" << root->val << ", idx=" <<idx <<endl;
        // 在inorder[l2:r2]中找root的位置
        int mid = l2;
        while(inorder[mid]!=root->val) mid++;
        int left_size = mid-l2;
        root->left = _buildTree(postorder,inorder,l1,l1+left_size-1, l2,mid-1);
        root->right = _buildTree(postorder,inorder,l1+left_size, r1-1, mid+1, r2);
        return root;
    }
};

// 使用hashmap优化在inorder中寻找root的位置。
class Solution2 {
public:
    TreeNode* buildTree(vector<int>& inorder, vector<int>& postorder) {
        inOrderMap = unordered_map<int,int>(inorder.size());
        for (int i = 0; i< inorder.size(); i++) {
            inOrderMap[inorder[i]]=i;
        }
        TreeNode* root = buildTreeNode(inorder,postorder,0, postorder.size()-1, 0);
        return root;
    }
private:
    TreeNode* buildTreeNode(vector<int>& inorder, vector<int>& postorder, int postStart, int postEnd, int inStart) {
        // 边界条件
        if (postStart > postEnd) return nullptr;
        int rootVal = postorder[postEnd];
        TreeNode* root = new TreeNode(rootVal);
        int rootPos = inOrderMap[rootVal];
        int len = rootPos-inStart;
        TreeNode* tree_l = buildTreeNode(inorder, postorder, postStart, postStart+len-1, inStart);
        TreeNode* tree_r = buildTreeNode(inorder, postorder,postStart+len, postEnd-1, rootPos+1);
        root->left = tree_l;
        root->right = tree_r;
        return root;
    }
    unordered_map<int,int> inOrderMap;
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
    {
        Solution2 s;
        auto result = s.buildTree(inorder,postorder);
        cout << result << endl;
    }
    return 0;
}