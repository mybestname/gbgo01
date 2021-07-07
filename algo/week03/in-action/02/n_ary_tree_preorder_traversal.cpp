#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;

class Node {
public:
    int val;
    vector<Node*> children;

    Node() {}

    Node(int _val) {
        val = _val;
    }

    Node(int _val, vector<Node*> _children) {
        val = _val;
        children = _children;
    }
};

// 使用递归
class Solution {
public:

    vector<int> preorder(Node* root) {
        // 前序：根左右
        seq = {};
        find(root);
        return seq;
    }
private:
    vector<int> seq;
    void find(Node *root) {
        if (root == nullptr) return;
        seq.push_back(root->val);
        for (auto n : root->children) {
            find(n);
        }
    }
};

// 使用迭代法
class Solution2 {
public:
    vector<int> preorder(Node* root) {
        if (root == nullptr) return {};
        st.push(root);
        while(!st.empty()) {
            auto node = st.top();
            st.pop();
            ans.push_back(node->val);
            for (int i = node->children.size()-1; i >= 0; i--) {
                st.push(node->children[i]);
            }
        }
        return ans;
    }
private:
    stack<Node*> st;
    vector<int> ans;
};

int main() {
    //
    //                             1
    //                    /  /     |     \
    //                  2    3     4      5
    //                     /  \    |    /  \
    //                     6   7   8    9  10
    //                         |   |    |
    //                        11   12   13
    //                         |
    //                        14
    //
    // 输入：root = [1,null,2,3,4,5,null,null,6,7,null,8,null,9,10,null,null,11,null,12,null,13,null,null,14]
    // 输出：[1,2,3,6,7,11,14,4,8,12,5,9,13,10]
    //
    vector<Node*> children;
    Node* n2 = new Node(2);
    Node* n3 = new Node(3, {new Node(6), new Node(7, {new Node(11, {new Node(14)})})});
    Node* n4 = new Node(4, {new Node(8,{new Node(12)})});
    Node* n5 = new Node(5, {new Node(9,{new Node(13)}), new Node(10)});
    children.push_back(n2);
    children.push_back(n3);
    children.push_back(n4);
    children.push_back(n5);
    Node* root = new Node(1, children);

    {
        Solution s;
        auto preorder = s.preorder(root);
        cout << "n-ary preorder " << preorder << endl;
    }

    {
        Solution2 s;
        auto preorder = s.preorder(root);
        cout << "n-ary preorder " << preorder << endl;
    }
    return 0;
}

