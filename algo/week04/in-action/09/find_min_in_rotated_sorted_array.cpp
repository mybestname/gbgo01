#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"
using namespace std;

class Solution {
public:
    int findMin(vector<int>& nums) {
        int left = 0, right = nums.size()-1;
        while (left < right) {
            int mid = (left + right) >> 1;
            if (nums[mid] <= nums[right])
                right = mid;
            else
                left = mid+1;
        }
        return nums[left];
    }
};

int main() {
    struct Test {
        vector<int> nums;
        int expect;
    };
    vector<Test> tests = {
            {
                    .nums = {3,4,5,1,2},
                    .expect = 1
            },
            {
                    .nums = {4,5,6,7,0,1,2},
                    .expect = 4
            },
            {
                    .nums = {11,13,15,17},
                    .expect = 11
            },
    };

    Solution s;
    for (auto &test : tests) {
        auto result = s.findMin(test.nums);
        cout << "nums=" << test.nums << ",expect=" << test.expect <<",got=" << result << endl;
    }
    return 0;
}

