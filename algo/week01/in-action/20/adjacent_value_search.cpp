#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<int> adjacent_value(vector<int>& nums) {
        int a[SIZE], rk[SIZE], ans[SIZE];
        Node* pos[SIZE];

        for (int i = 0; i < SIZE; i++)
        {
           a[i]= 0, rk[i]=0, ans[i]=0;
        }

        uint n = nums.size();
        for (int i = 1; i<= n; i++) {
            a[i] = nums[i-1];
            rk[i] = i;
        }
        cout <<"rk="<< rk << ", a=" << a << endl;
        sort(rk+1,rk+n+1, [&a](int i, int j){ return a[i]< a[j];});
        cout <<"rk="<< rk << ", a=" << a << endl;

        head.val = numeric_limits<int>::min();  //保护节点，最小值
        tail.val = numeric_limits<int>::max();  //保护节点，最大值
        head.next = &tail;
        tail.pre = &head;

        for (int i = 1; i<= n; i++) {
            pos[rk[i]] = addNode(tail.pre, rk[i], a);
        }

        for (int i = n; i>1; i--) {
            Node* pre = pos[i]->pre;
            Node* next = pos[i]->next;
            // 比较前驱和后续，进行选择
            if (a[i]-pre->val < next->val - a[i] ||
                a[i]-pre->val == next->val - a[i] && pre->val < next->val){
                ans[i] = pre->idx;
            }else {
                ans[i] = next->idx;
            }
            //  然后删除该节点
            deleteNode(pos[i]);
        }
        cout <<"rk="<< rk << ", a=" << a << ", ans="<< ans <<endl;

        vector<int> result;
        for (int i = 2; i<=n; i++) {
           result.push_back(abs(a[i]-a[ans[i]]));
           result.push_back(ans[i]);
        }

        return result;
    }


private:
    static const int SIZE=10;
    struct Node{
        int val, idx;
        Node* pre; Node* next;
    }head, tail;

    Node* addNode(Node* p, int idx, const int* a) {
        Node* q = new Node();
        q->idx = idx, q->val = a[idx];
        q->pre = p, q->next = p->next;
        p->next->pre = q, p->next = q;
        return q;
    }
    void deleteNode(Node* p) {
        p->pre->next = p->next;
        p->next->pre = p->pre;
        delete p;
    }
};

int main() {
    {
        vector<int> nums = {1, 5, 3};
        Solution s{};
        vector<int> result = s.adjacent_value(nums);
        cout << "nums=" << nums << ",result=" << result << endl;
        // input [1 5 3]
        // output [4 1 2 1]
    }
    {
        vector<int> nums = {1,8,5,7,3,6};
        Solution s;
        vector<int> result = s.adjacent_value(nums);
        cout << "nums=" << nums << ",result=" << result << endl;
        // input [1 8 5 7 3 6]
        // output [7,1,3,2,1,2,2,1,1,3]
    }
    return 0;
}
