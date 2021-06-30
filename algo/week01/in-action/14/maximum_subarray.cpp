#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int maxSubArray(vector<int>& nums) {
        int n = nums.size();
        vector<int64_t> sum(n+1,0); // 0~n
        for (int i = 1; i<=n; i++) sum[i] = sum[i-1] + nums[i-1]; //求前缀和 1~n
        int64_t ans = 0;
        int64_t prefix_min = 0 ;  // 中间变量用于存储前缀最小值
        for (int i= 1; i<= n ; i++) { //1~n
           ans = max(ans, sum[i]-prefix_min);
           prefix_min = min(prefix_min, sum[i]);
        }
        return ans;
    }
};

int main() {
    vector<int> nums = {-2,1,-3,4,-1,2,1,-5,4} ;
    int max_sum = Solution::maxSubArray(nums);
    cout << "nums=" << nums << ", max_sum=" << max_sum << endl;
    return 0;
}