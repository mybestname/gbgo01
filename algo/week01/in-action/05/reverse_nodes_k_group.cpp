#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static ListNode *reverseKGroup(ListNode *head, int k) {
        ListNode *protect = new ListNode(-1,head); // ------- (4.1)
        ListNode *last = protect;                     // ------- (4.2)
        while (head != nullptr) {                     // ------- (1)
            // 找到分组的结束节点
            ListNode* end = getEndNodeOfK(head,k);    // ------- (1.1)
            // 需要额外考虑边界问题。
            if (end == nullptr) { break; }            // --------(5)
            // 先存一下到下一组的指针
            ListNode* nextGroupHead = end->next;      // ------  (2.1)
            // 反转分组（处理head到end之间的k-1条边的反转）
            reverseList(head, end);                   // ------- (1.2)
            // 处理该组的head和end                      // ------- (1.3)
            // 1.) 该组的end，现在是头，需要连给上一组的尾
            //     ？上组尾 = end;
            last->next = end;                          // ------ (3.1)
            // 2.) 该组的head，现在是尾，需要连到下一组的头
            //    head.next= ？下组头
            //    本来应该是end.next，但是end.next已经在反转函数中被修改过了
            //    需要有一个额外存的地方，所以提前保存到nextGroupHead
            head->next= nextGroupHead;                 // ------ (3.2)
                                                       //
            // 迭代到下一组                              //
            last = head;                               // ------ (2.2)
            head = nextGroupHead;                      // ------ (2.3)
        }
        // 注意！返回什么？
        // last是上一组的尾，并不是真正我们想返回的。
        return protect->next;                          // ------ (4.3)
        // 思考顺序
        // (1) 大的while框架：对每个分组进行处理
        //     - 1. 先找分组的结束
        //     - 2. 对分组进行反转
        //     - 3. 处理分组的head和end
        // (2) while如何对分组进行迭代，必须通过中间变量的帮助，从而引出nextGroupHead
        // (3) 对于分组的head和end如何处理
        // (4) 思考返回什么，引出关于保护节点的作用。
        // (5) 需要额外考虑的边界问题。
    }
private:
    static void reverseList(ListNode* head, ListNode* end) { //不需要返回，修改为void
        // 边界检查问题
        if (head == end) return; // - (4) 另外需要处理边界检查，如果k=1的情况，head==end，直接返回。
        ListNode *last = head;   // - (2.1) 第一个需要改的边的内容应该是head（即head的下一个需要指向head）
        head = head->next;       // - (2.2) head不需要改，直接从head的下一个开始。
        while (head != end) {    // - (1) 终止标志修改为end节点。
            ListNode *nextHead = head->next;
            head->next = last;
            last = head;
            head = nextHead;
        }
        // 当 head == end 时候，即返回条件的时候，
        // 还需要修改终止节点的指向，因为end节点还没有反转，end节点目前还是指向下一个节点
        end->next = last;        // - (3) end节点需要指向上一个节点。
        // 思考顺序:
        // (1): 首先建立大的while迭代模版和终止标志
        // (2): 把修改的开始节点限定为head的下一个节点，因为head不需要修改。
        // (3): 终止条件时需要额外处理end节点
        // (4): 添加额外的输入检查
    }
    static ListNode* getEndNodeOfK(ListNode* head, int k) {
        while (head != nullptr) {
            k--;                     //先--
            if (k==0) break;         //这样调用时候直接用k， `getEndNodeOfK(head,k)` 不用 `getEndNodeOfK(head, k-1)`
            head = head->next;
        }
        return head;
    }
};

int main() {

    ListNode* l1 = new ListNode(1,
   new ListNode(2,
   new ListNode(3,
   new ListNode(4,
   new ListNode(5, nullptr))))); // [1,2,3,4,5]
    int k = 2;
    ListNode* result = Solution::reverseKGroup(l1, k);
    cout << result << endl; // [2,1,4,3,5]
    return 0;
}



