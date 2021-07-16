#include <iostream>
#include <vector>
#include <queue>
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

class Solution {
public:
    vector<vector<int>> levelOrder(Node* root) {
        vector<vector<int>> seq;
        queue<pair<Node*,int>> q;  //每个Node在队列里面存两个信息，该Node指针和深度（层数）
        q.push(make_pair(root,0));
        while (!q.empty()) {
            Node* head = q.front().first;
            int depth = q.front().second;
            q.pop();
            for(auto child : head->children){
                q.push(make_pair(child,depth+1));  //child按从左到右入队。
            }
            /*
            if (depth + 1> seq.size() ) seq.push_back({});
            // 注意：size()是uint, 不能用0 和 `seq.size()-1`去比较，因为右项永远大于0。
            // 语句 `depth > seq.size() - 1` 的隐含条件是seq.size()-1不为负，和逻辑直觉相谬。
            */
            if (seq.size() <= depth) seq.push_back({}); // 用 <= 更舒服
            seq[depth].push_back(head->val);
        }
        return seq;
    }
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
    // 输出：[[1],[2,3,4,5],[6,7,8,9,10],[11,12,13],[14]]
    //
    vector<Node *> children;
    Node *n2 = new Node(2);
    Node *n3 = new Node(3, {new Node(6), new Node(7, {new Node(11, {new Node(14)})})});
    Node *n4 = new Node(4, {new Node(8, {new Node(12)})});
    Node *n5 = new Node(5, {new Node(9, {new Node(13)}), new Node(10)});
    children.push_back(n2);
    children.push_back(n3);
    children.push_back(n4);
    children.push_back(n5);
    Node *root = new Node(1, children);

    {
        Solution s;
        auto preorder = s.levelOrder(root);
        cout << "n-ary preorder " << preorder << endl;
    }

    return 0;
}
