#include <iostream>
#include <vector>
#include <deque>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    static vector<int> maxSlidingWindow(vector<int>& nums, int k) {
        vector<int> ans;
        //存放下标的双端队列 (因为下标代表时间)
        deque<int> q;
        for (int i = 0 ; i< nums.size(); i++) {
            // 合法性判断
            while (!q.empty() && q.front() <= i-k) { // i-k为出界的判断，比i-k更小，说明已经expired
                q.pop_front(); // expire, 删除队头
            }
            // 维护单调性
            while( !q.empty() && nums[q.back()] <= nums[i]){  //队尾值比nums[i]还差，说明队尾为冗余
                q.pop_back();  // 冗余，删除队尾
            }
            // 新元素入队
            q.push_back(i);
            if (i >= k-1) {
                ans.push_back(nums[q.front()]);  //最优
            }
        }
        return ans;
    }
};


int main() {
    vector<int> nums = {1,3,-1,-3,5,3,6,7};
    int k = 3 ;
    vector<int> result = Solution::maxSlidingWindow(nums,k);
    cout << "nums=" << nums << ",result=" << result << endl;
    // 输入：nums = [1,3,-1,-3,5,3,6,7], k = 3
    // 输出：[3,3,5,5,6,7]

    return 0;
}
