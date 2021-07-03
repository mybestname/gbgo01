#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;


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
   static ListNode* mergeTwoList2(ListNode* l1, ListNode* l2) {
        ListNode protect = ListNode(-1, nullptr); //保护节点
        ListNode* tail = &protect;                         //保护节点保存到tail
        while(l1!= nullptr && l2 != nullptr) {
            if (l1->val < l2->val) {
                tail->next = l1;
                l1 = l1->next;
            }else {
                tail->next = l2;
                l2 = l2->next;
            }
            tail = tail->next;
        }
        if (l1 != nullptr) tail->next = l1; else tail->next = l2;
        return protect.next;
    }

};

int main() {
    {
        ListNode l1n1 = ListNode(4, nullptr);
        ListNode l1n2 = ListNode(2, &l1n1);
        ListNode l1n3 = ListNode(1, &l1n2);
        ListNode *l1 = &l1n3;

        ListNode l2n1 = ListNode(4, nullptr);
        ListNode l2n2 = ListNode(3, &l2n1);
        ListNode l2n3 = ListNode(1, &l2n2);
        ListNode *l2 = &l2n3;

        ListNode *result = Solution::mergeTwoList(l1, l2);
        cout << result << endl; //[1,1,2,3,4,4]
    }

    {
        ListNode l1n1 = ListNode(4, nullptr);
        ListNode l1n2 = ListNode(2, &l1n1);
        ListNode l1n3 = ListNode(1, &l1n2);
        ListNode* l1 = &l1n3;

        ListNode l2n1 = ListNode(4, nullptr);
        ListNode l2n2 = ListNode(3, &l2n1);
        ListNode l2n3 = ListNode(1, &l2n2);
        ListNode* l2 = &l2n3;

        ListNode* result = Solution::mergeTwoList2(l1,l2);
        cout << result << endl; //[1,1,2,3,4,4]
    }

}
