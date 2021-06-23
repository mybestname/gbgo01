#include <iterator> // needed for std::ostream iterator
#include <iostream>
#include <vector>
#include <deque>
#include <stack>
template <typename T>
std::ostream& operator<< (std::ostream& out, const std::vector<T>& v) {
    if ( !v.empty() ) {
        out << '[';
        std::copy (v.begin(), v.end(), std::ostream_iterator<T>(out, ","));
        out << "\b]";
    }
    return out;
}
using namespace std;

class Node {
public:
    int val;
    vector<Node*> children;
    Node(){ val = -1; }
    explicit Node(int _val){
        val = _val;
    }
    Node(int _val, vector<Node*> _children) {
        val = _val;
        children = std::move(_children);
    }
};

class Solution {
public :
    // using std::stack
    static vector<int> preorder_s(Node* root) {
        vector<int> result = {};
        // 空输入检查
        if (root == nullptr)
            return result;
        // 构造stack
        stack<Node*> s;
        s.push(root);                          // 根节点入栈  vector.push_back()
        while (!s.empty()) {                   // vector.size() == 0
            Node* node = s.top();              // c.back()  //取栈顶节点
            s.pop();  // c.pop_back()          // 弹出栈顶
            result.push_back(node->val);       // 存入到返回列表
            // 节点的所有子节点，按逆序入栈，这样才能保证出栈顺序是正确的。
            for (int i = node->children.size() - 1 ; i >= 0; i--){
                Node* n = node->children[i];
                if (n != nullptr) {             // 非空检查
                    s.push(n);
                }
            }
        }
        return result;
   }
   // using std::deque
   // std::stack 其实是deque的封装（默认）
   //  - stack 是容器适配器(adapter)
   //  - 也可适配 vector、list
   // 效率
   // - deque::push_back() 是O(1)
   // - vector::push_back() 是均摊O(1)
   static vector<int> preorder_d(Node* root) {
        vector<int> result = {};
        // 空输入检查
        if (root == nullptr)
            return result;
        // deque for stack
        deque<Node*> q;
        q.push_back(root);                      // 根节点入队列尾
        while (!q.empty()) {                    // vector.size() == 0
            Node* node = q.back();              // 取队尾节点，等效于取栈顶
            q.pop_back();                       // 弹出队尾，等效于弹出栈顶
            result.push_back(node->val);        // 存入到返回列表
            // 节点的所有子节点，按逆序入栈，这样才能保证出栈顺序是正确的。
            // 即使是使用deque，也需要按照逆序插入队尾，才能保证顺序，不要试图从队头插入。
            for (int i = node->children.size() - 1; i >= 0 ; i--){
                Node* n = node->children[i];
                if (n != nullptr) {             // 非空检查
                    q.push_back(n);
                }
            }
        }
        return result;
    }

    // using std::vector
    static vector<int> preorder_v(Node* root) {
        vector<int> result = {};
        // 空输入检查
        if (root == nullptr)
            return result;
        // vector as stack
        vector<Node*> stack;
        stack.push_back(root);                      // 根节点入队列尾
        while (!stack.empty()) {                    // vector.size() == 0
            Node* node = stack[stack.size()-1];     // 取队尾节点，等效于取栈顶
            stack.pop_back();                       // 弹出队尾，等效于弹出栈顶
            result.push_back(node->val);            // 存入到返回列表
            // 节点的所有子节点，按逆序入栈，这样才能保证出栈顺序是正确的。
            for (int i = node->children.size() - 1 ; i >= 0; i--){
                stack.push_back(node->children[i]);
            }
        }
        return result;
    }

};

int main() {
    Solution sol;

    Node* n5 = new Node(5);
    Node* n6 = new Node(6);
    Node* n4 = new Node(4);
    Node* n2 = new Node(2);
    Node* n3 = new Node(3, {n5,n6});
    Node* root = new Node(1, {n3, n2, n4});

    vector<int> result_s = Solution::preorder_s(root);
    vector<int> result_d = Solution::preorder_d(root);
    vector<int> result_v = Solution::preorder_v(root);

    cout << result_s << endl;
    cout << result_d << endl;
    cout << result_v << endl;


    // 输入：root = [1,null,2,3,4,5,null,null,6,7,null,8,null,9,10,null,null,11,null,12,null,13,null,null,14]
    //
    //            1
    //    2    3     4    5
    //       6   7   8   9  10
    //          11  12  13
    //         14
    //
    // 输出：[1,2,3,6,7,11,14,4,8,12,5,9,13,10]
    Node* n14 = new Node(14);
    Node* n11 = new Node(11, {n14});
    Node* n12 = new Node(12);
    Node* n13 = new Node(13);
    n6  = new Node(6);
    Node* n7  = new Node(7,{n11});
    Node* n8  = new Node(8, {n12});
    Node* n9  = new Node(9, {n13});
    Node* n10 = new Node(10);
    n2 = new Node(2);
    n3 = new Node(3, {n6,n7});
    n4 = new Node(4, {n8});
    n5 = new Node(5, {n9,n10});
    root = new Node(1, {n2,n3,n4,n5});
    result_s = Solution::preorder_s(root);
    result_d = Solution::preorder_d(root);
    result_v = Solution::preorder_v(root);
    cout << result_s << endl;
    cout << result_d << endl;
    cout << result_v << endl;

    return 0;
}


