#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"
using namespace std;

// 使用堆（优先队列）的解法。对比使用单调队列的解法（见week01 例题18）
class Solution {
public:
    vector<int> maxSlidingWindow(vector<int>& nums, int k) {
        vector<int> ans;
        // 使用一个<值，下标>的pair，和优先队列配合完成lazy delete
        priority_queue<pair<int,int>> q;
        for (int i = 0; i<k-1; i++) q.push(make_pair(nums[i],i));  //前k-1个数先加入队列。
        for (int i = k-1; i<nums.size(); i++) { //从第k个数开始
            q.push(make_pair(nums[i],i)); //先数加入队列
            //现在需要判断是否是合法答案，如果不合法（超出窗口），则要pop掉（删除）
            while(q.top().second <= i-k) q.pop();
            //把合法答案加入答案列表
            ans.push_back(q.top().first);
        }
        return ans;
    }
};

int main() {
    struct Test {
        vector<int> nums;
        int k;
        vector<int> expect ;

    };
    vector<Test> tests = {
            {
                .nums   = {1,3,-1,-3,5,3,6,7},
                .k      = 3,
                .expect = {3,3,5,5,6,7},
            },
            {.nums = {1},    .k = 1, .expect = {1}},
            {.nums = {1,-1}, .k = 1, .expect = {1,-1}},
            {.nums = {9,11}, .k = 2, .expect = {11}},
            {.nums = {4,-2}, .k = 2, .expect = {4}},
    };

    Solution s;
    for (auto &test : tests) {
        auto result = s.maxSlidingWindow(test.nums,test.k);
        cout << "numbs=" << test.nums << ", k=" << test.k
             << ",expect=" << test.expect << ",got=" << result << endl;
    }

    return 0;
}
