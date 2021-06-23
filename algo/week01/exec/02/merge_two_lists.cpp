#include <iostream>
#include <vector>
#include <iterator> // needed for std::ostream iterator


using namespace std;

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

class Solution {
public :
   static ListNode* mergeTwoList(ListNode* l1, ListNode* l2) {
        if (l1 == nullptr) return l2;
        if (l2 == nullptr) return l1;
        if (l1->val < l2->val) {
            l1->next = mergeTwoList(l1->next, l2);
            return l1;
        }
        l2->next = mergeTwoList(l1, l2->next);
        return l2;
    }
};

int main() {
    Solution sol;

    ListNode l1n1 = ListNode(4, nullptr);
    ListNode l1n2 = ListNode(2, &l1n1);
    ListNode l1n3 = ListNode(1, &l1n2);
    ListNode* l1 = &l1n3;

    ListNode l2n1 = ListNode(4, nullptr);
    ListNode l2n2 = ListNode(3, &l2n1);
    ListNode l2n3 = ListNode(1, &l2n2);
    ListNode* l2 = &l2n3;

    ListNode* result = Solution::mergeTwoList(l1,l2);
    cout << result << endl; //[1,1,2,3,4,4]

}
