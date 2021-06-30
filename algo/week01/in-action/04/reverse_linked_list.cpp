#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static ListNode* reverseList(ListNode* head) {
        ListNode* last = nullptr;
        // 1. 因为要修改链表的每一条边，所以需要遍历链表
        // 2. 什么时候用`head != nullptr`, 什么时候用`head.next != nullptr`
        //    - 取决于是遍历n次，还是n-1次，因为要改n条边，所以应该用head!=nullptr
        while(head != nullptr) {
            ListNode* nextHead = head->next; // 先保存原来的next信息，供遍历使用。因为在本题中这个边是需要修改的，所以需要保存。
            head->next = last;   // ------- 这是使用的是上一次迭代的结果-- 这里改了一条边
                                 //                                 +--- 这两个语句，连同前后两个迭代，达到了翻转的目的。
            last = head;         // ------- 这里存起来为下一次迭代用----- 这里保存的值，用来在下一次迭代去修改下一个边
            head = nextHead;                 // 使用之前保存的值进行遍历。
        }
        return last;  //当 head == nullptr 时候，last正是我们的结果（完成的翻转）
        // 开始的样子
        //  last  head
        //    v   v
        // nil    1 -> 2 -> 3 -> 4 -> 5 -> nil
        //
        // 最后的样子
        //                             last head
        //                             v     v
        // nil <- 1 <- 2 <- 3 <- 4 <- 5    nil
        //
        //
        // 如果把具体的箭头都想象成独一无二的内存地址，那么就很好理解
        // head (0xa) 1 (0xb) 2 (0xc) 3 (0xd) 4 (0xe) 5 (0xf) nil  ( abcdef)
        //
        // last (0xe) 5 (0xd) 4 (0xc) 3 (0xb) 2 (0xa) 1 (0xf) nil  ( edcbaf)
        //
        // 那么这个问题其实还是可以reduce为对一个数组的元素进行替换的问题。
    }
};

int main() {
    ListNode l1n1 = ListNode(4, nullptr);
    ListNode l1n2 = ListNode(2, &l1n1);
    ListNode l1n3 = ListNode(1, &l1n2);
    ListNode* l1 = &l1n3;   // [1,2,4]
    ListNode* result = Solution::reverseList(l1);
    cout << result << endl; // [4,2,1]
    return 0;
}

