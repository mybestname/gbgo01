#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"
using namespace std;

// 使用模版1.1
class Solution {
public:
    vector<int> searchRange(vector<int>& nums, int target) {
        // 开始位置？查询第一个 `>= target` ，即查询low_bound
        // 范围 [0..n-1],n (额外多一个表示不存在）
        vector<int> ans;
        int left = 0, right = nums.size();
        while (left < right) {
            int mid = (left + right) >> 1;
            // 满足条件分支
            if (nums[mid] >= target ) {
                right = mid;
            }else {
                left = mid+1;
            }
        }
        ans.push_back(right) ;
        // 结束位置？查询最后一个 `<= target`，
        left = -1, right = nums.size()-1;
        while (left < right) {
            int mid = (left+right+1) >> 1;
            if (nums[mid] <= target) {
                left = mid;
            }else {
                right = mid-1;
            }
        }
        ans.push_back(left);
        // 处理返回结果
        if (ans[0] > ans [1]) return {-1,-1};
        return ans;
    }
};

int main() {
    struct Test {
        vector<int> nums;
        int target;
        vector<int> expect;
    };
    vector<Test> tests = {
            {
                    .nums   = {5,7,7,8,8,10},
                    .target =  8,
                    .expect =  {3,4},
            },
            {
                    .nums   = {5,7,7,8,8,10},
                    .target =  6,
                    .expect =  {-1,-1},
            },
            {
                    .nums   = {},
                    .target =  0,
                    .expect =  {-1,-1},
            },
    };

    Solution s;
    for (auto &test : tests) {
        auto result = s.searchRange(test.nums, test.target);
        cout << "nums=" << test.nums << ",target=" << test.target << ",expect=" << test.expect <<",got=" << result << endl;
    }
    return 0;
}
