#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"
using namespace std;

// 优先队列模版，key是用于比较的关键码，listNode可以是任何附带信息。
struct Node {
    int key;
    ListNode* listNode;
};
bool operator <(const Node& a, const Node& b) {
    // 大根堆的写法，因为priority_queue默认是优先排大值的。
    // return a.key < b.key;

    // 注意：这是小根堆的写法，这种重载<，使得priority_queue<Node> 变为小根堆。
    // 因为本题需要按升序排序，所以要找极小值，所以使用小根堆
    return a.key > b.key;
}

// 解法1，使用标准库的priority_queue
class Solution1 {
public:
    ListNode* mergeKLists(vector<ListNode*>& lists) {
        // 使用优先队列进行堆排序找到最小值：logK的复杂度
        // Node通过key进行排序。
        priority_queue<Node> pq;
        for (ListNode* node : lists) {
            if (node != nullptr) {
                pq.push(Node{node->val, node}); //这里push的是k个表头。
            }
        }
        ListNode head(-1, nullptr); // protect node
        ListNode* tail = &head ;
        while(!pq.empty()) {
            Node node = pq.top();
            pq.pop();
            // 在答案链表的末尾插入当前的最小值节点。
            tail->next = node.listNode;
            tail = tail->next;
            // 当最小值取出后，需要把最小值节点对应的链表的下一个节点（如果还有下一个话）入堆。
            ListNode* next = node.listNode->next;
            if ( next != nullptr) {
                // 这其实是pop一个最小值就push一个值的策略，这样会触发堆排序，自动把所有堆中节点的最小值推到根。
                // 因为下一个出队的节点一定是堆中的最小节点，所以它指向的链表的下一个节点一定是合适入堆的最小节点。
                // 所以可以保证最终合并的有序性。
                pq.push(Node{next->val, next});
            }
        }
        return head.next;
    }
};

// 解法2，使用自定义的2叉堆
class Solution2 {
public:
    ListNode *mergeKLists(vector<ListNode *> &lists) {
        // 使用优先队列进行堆排序找到最小值：logK的复杂度
        // Node通过key进行排序。
        BinaryHeap heap;
        for (ListNode* node : lists) {
            if (node != nullptr) {
                heap.push(Node{node->val, node}); //这里push的是k个表头。
            }
        }
        ListNode head(-1, nullptr); // protect node
        ListNode* tail = &head ;
        while(!heap.empty()) {
            Node node = heap.pop();
            // 在答案链表的末尾插入当前的最小值节点。
            tail->next = node.listNode;
            tail = tail->next;
            // 当最小值取出后，需要把最小值节点对应的链表的下一个节点（如果还有下一个话）入堆。
            ListNode* next = node.listNode->next;
            if ( next != nullptr) {
                // 这其实是pop一个最小值就push一个值的策略，这样会触发堆排序，自动把所有堆中节点的最小值推到根。
                // 因为下一个出队的节点一定是堆中的最小节点，所以它指向的链表的下一个节点一定是合适入堆的最小节点。
                // 所以可以保证最终合并的有序性。
                heap.push(Node{next->val, next});
            }
        }
        return head.next;
    }
private:
    class BinaryHeap {
    public:
        bool empty() {
            return heap.empty();
        }
        void push(const Node& n) {
            heap.push_back(n);  //首先push到尾部
            // 然后往上做调整
            int p = heap.size()-1;
            // 上溯到根
            while (p > 0 ) {
                int fa = (p-1)/2;  // 根在0，father = (child -1)/2
                // 堆性质检查
                if (heap[p].key < heap[fa].key) {  // 小根堆，谁小谁在前面
                    swap(heap[p],heap[fa]); // 满足则上调
                    p = fa;                        // 更新当前p的位置
                }else break;                       // 不满足直接终止
            }                                      // 或在根节点终止
        }
        Node pop(){
            Node top = heap[0];
            // 首先 堆顶 和 堆尾 交换
            heap[0] = heap[heap.size()-1];
            heap.pop_back(); // 同时删除尾部
            // 从堆顶向下调整。
            // `father * 2 + 1 = left_child`
            // `father * 2 + 2 = right_child`
            int p = 0;
            int child = p*2 + 1;
            while (child < heap.size()) {            // child没有出界，p有孩子
                // 注意：使用other 和 child，而不用left，right有原因，这样写法后面判断写法更简洁。
                int other_child = p*2 + 2;
                if (other_child < heap.size() && heap[child].key > heap[other_child].key)
                    child = other_child;             // 存在另外一个，并且另外一个更小，则更新为更小的。
                if (heap[p].key > heap[child].key) { // 满足小根堆交换规则
                    swap(heap[p],heap[child]);
                    p = child;
                    // 注意：这里不要忘记把child更新为下一个
                    child = p*2+1;
                }else break;                         // 不满足直接停止
            }                                        // 或在出界时候停止
            return top;
        }

    private:
        // 通过数组存储完全二叉树实现二叉堆
        // root为0
        vector<Node> heap;
    };
};

int main(){
    struct Test {
        vector<ListNode*> lists;
        ListNode* expect ;

    };

    {
        vector<Test> tests = {
            // 输入：lists = [[1,4,5],[1,3,4],[2,6]]
            // 输出：[1,1,2,3,4,4,5,6]
            {
                    .lists =
                            {
                                    new ListNode(1, new ListNode(4, new ListNode(5))),
                                    new ListNode(1, new ListNode(3, new ListNode(4))),
                                    new ListNode(2, new ListNode(6)),
                            },
                    .expect = new ListNode(1, new ListNode(1, new ListNode(2,new ListNode(3, new ListNode(4, new ListNode(4,new ListNode(5, new ListNode(6, nullptr)))))))),
            },
            // 输入：lists = []
            // 输出：[]
            {
                    .lists  = {},
                    .expect = nullptr,
            },
            //输入：lists = [[]]
            //输出：[]
            {
                    .lists  = { nullptr, },
                    .expect = nullptr,
            },
        };

        Solution1 s;
        for (auto &test : tests) {
            auto result = s.mergeKLists(test.lists);
            cout << "S1: lists=" << test.lists
                 << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }

    {
        vector<Test> tests = {
                // 输入：lists = [[1,4,5],[1,3,4],[2,6]]
                // 输出：[1,1,2,3,4,4,5,6]
                {
                        .lists =
                                {
                                        new ListNode(1, new ListNode(4, new ListNode(5))),
                                        new ListNode(1, new ListNode(3, new ListNode(4))),
                                        new ListNode(2, new ListNode(6)),
                                },
                        .expect = new ListNode(1, new ListNode(1, new ListNode(2,new ListNode(3, new ListNode(4, new ListNode(4,new ListNode(5, new ListNode(6, nullptr)))))))),
                },
                // 输入：lists = []
                // 输出：[]
                {
                        .lists  = {},
                        .expect = nullptr,
                },
                //输入：lists = [[]]
                //输出：[]
                {
                        .lists  = { nullptr, },
                        .expect = nullptr,
                },
        };
        Solution2 s;
        for (auto &test : tests) {
            auto result = s.mergeKLists(test.lists);
            cout << "S2: lists=" << test.lists
                 << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}