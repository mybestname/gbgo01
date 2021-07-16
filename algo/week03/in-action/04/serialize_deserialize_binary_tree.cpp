#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"

using namespace std;
class Codec {
public:

    // Encodes a tree to a single string.
    string serialize(TreeNode* root) {
        ans = {};
        preorderTravel(root);
        string result;
        for (auto s: ans) result+=s+' ';
        return result.erase(result.rfind(' '), string::npos);
    }

    // Decodes your encoded data to tree.
    TreeNode* deserialize(string data) {
        vector<string> seq;
        // 1 2 null null 3 4 null null 5 null null
        split(seq, data, ' ');
        cout << seq << endl;
        int current = 0;
        return buildTree(seq, current);  //注意这里current变量是一个引用，代表一个共享变量。
    }

private:
    void split(vector<string>& seq, const string& str, char delim) {
        size_t current=0;
        size_t next = -1;
        do
        {
            current = next+1;
            next = str.find_first_of( delim, current );
            seq.push_back(str.substr( current, next - current ));
        }
        while (next != string::npos);
    }
    TreeNode* buildTree(vector<string>& seq, int& current) {
        // [1,2,null,null, 3, 4, null,null, 5, null,null]
        if (seq[current] == "null") {
            current++;
            return nullptr;
        }
        // 复原：根左右
        TreeNode* root = new TreeNode(stoi(seq[current]));
        current++;
        root->left = buildTree(seq, current);
        root->right = buildTree(seq, current);
        return root;
    }

    void preorderTravel(TreeNode* root) {
        // 边界
        if (root== nullptr) {
            ans.push_back("null");
            return;
        }
        // 根左右
        ans.push_back(to_string(root->val));
        preorderTravel(root->left);
        preorderTravel(root->right);
    }
    vector<string> ans;
};

// Your Codec object will be instantiated and called as such:
// Codec ser, deser;
// TreeNode* ans = deser.deserialize(ser.serialize(root));

int main() {
    // 输入：root = [1,2,3,null,null,4,5]
    // 输出：[1,2,3,null,null,4,5]
    TreeNode* root = new TreeNode(1);
    root->left = new TreeNode(2);
    root->right = new TreeNode(3);
    root->right->left = new TreeNode(4);
    root->right->right = new TreeNode(5);
    Codec ser, deser;
    cout << root << endl;
    TreeNode* ans = deser.deserialize(ser.serialize(root));
    cout << ans << endl;
}

