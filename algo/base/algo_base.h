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
            }
            if (i < v.size()-1) {
                out << "\n";
            }
        }
        out << "\n]";
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
#endif //ALGO2021_ALGO_BASE_H
