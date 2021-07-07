#ifndef ALGO2021_ALGO_BASE_H
#define ALGO2021_ALGO_BASE_H
#include <iterator> // needed for std::ostream iterator
template <unsigned int N>
std::ostream& operator<< (std::ostream& out, const int (&arr)[N]) {
    out << '[';
    for(int i = 0; i < N; i++) out << arr[i] << ",";
    out << "\b]";
    return out;
}

template <typename T>
std::ostream& operator<< (std::ostream& out, const std::vector<T>& v) {
    if ( !v.empty() ) {
        out << '[';
        std::copy (v.begin(), v.end(), std::ostream_iterator<T>(out, ","));
        out << "\b]";
    } else {
        out << "[]";
    }
    return out;
}

template <typename T>
std::ostream& operator<< (std::ostream& out, const std::vector<std::vector<T>>& v) {
    if ( !v.empty() ) {
        out << "[";
        for (int i= 0; i< v.size(); i++ ){
            std::vector<T> in = v[i];
            if (!in.empty()) {
                out << '[';
                std::copy(in.begin(), in.end(), std::ostream_iterator<T>(out, ","));
                out << "\b],";
            }else{
                out << "[],";
            }
            if (i < v.size()-1) {
                out << "\n";
            }
        }
        out << "\n]";
    }else{
        out <<"[]";
    }
    return out;
}

// Definition for singly-linked list.
class ListNode {
public:
    int val;
    ListNode *next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode *next) : val(x), next(next) {}
};
std::ostream& operator<< (std::ostream& out, const ListNode* v) {
    out << '[';
    if ( v != nullptr) { out << v->val ;}
    out << ",";
    for (ListNode* next= v->next ; next != nullptr; next = next->next) {
        out << next->val << ",";
    }
    out << "\b]";
    return out;
}

class TreeNode {
public:
    int val;
    TreeNode *left;
    TreeNode *right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode *left, TreeNode *right) : val(x), left(left), right(right) {}
};

void printTreeNode(std::ostream& out, const std::string& prefix, const TreeNode* node, bool isLeft);

std::ostream& operator<< (std::ostream& out, const TreeNode* t) {
    printTreeNode(out, "", t, false);
    return out;
}

void printTreeNode(std::ostream& out, const std::string& prefix, const TreeNode* node, bool isLeft)
{
    if( node != nullptr )
    {
        out << prefix;

        out << (isLeft ? "├──" : "└──" );

        // print the value of the node
        out << node->val << std::endl;

        // enter the next tree level - left and right branch
        printTreeNode( out, prefix + (isLeft ? "│   " : "    "), node->left, true);
        printTreeNode( out, prefix + (isLeft ? "│   " : "    "), node->right, false);
    }
}

#endif //ALGO2021_ALGO_BASE_H
