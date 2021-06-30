#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static bool hasCycle(ListNode *head) {
        ListNode* fast = head;
        while (fast != nullptr && fast->next != nullptr) { // while循环中限制fast的不为空条件
           fast = fast->next->next; // 一次走两步
           head = head->next;       // 一次走一步
           if (fast == head) {      // 相遇
              return true;          // 有环
           }
        }
        return false;
    }
};

int main() {

    ListNode* n4 = new ListNode(-4, nullptr);
    ListNode* n3 = new ListNode(0, n4);
    ListNode* n2 = new ListNode(2, n3);
    ListNode* n1 = new ListNode(3, n2);
    n4->next = n2;
    ListNode* l = n1;
    // head = [3,2,0,-4], pos = 1
    bool has_cycle = Solution::hasCycle(l);
    cout << "has_cycle=" << has_cycle << endl; //
    return 0;
}