#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static ListNode *detectCycle(ListNode *head) {
        ListNode* fast = head;
        ListNode* slow = head;
        while (fast != nullptr && fast->next != nullptr) { // while循环中限制fast的不为空条件
            fast = fast->next->next; // 一次走两步
            slow = slow->next;       // 一次走一步
            if (fast == slow) {      // 快慢相遇
                while (head != slow) {
                    head = head->next;  //头指针
                    slow = slow->next;  //和慢指针一起走
                }
                return head;            //直到头慢相遇，即为环的起点。
            }
        }
        return nullptr;
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
    ListNode* detected = Solution::detectCycle(l);
    cout << detected->val <<  endl; // 2
    return 0;
}

