#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"
using namespace std;

class Solution1 {
public:
    int findMin(vector<int>& nums) {
        int left = 0, right = nums.size()-1;
        while (left < right) {
            if (nums[left] < nums[right]){
                return nums[left];
            }
            int mid = left + ((right-left) >> 1);
            if (nums[mid] == nums[left]) {
                left++;
            }
            else if (nums[mid] < nums[left]){
                right = mid;
            }
            else {
                left = mid + 1;
            }
        }
        return nums[right];
    }
};

class Solution2 {
public:
    int findMin(vector<int>& nums) {
        int left = 0, right = nums.size()-1;
        while (left < right && nums[left] >= nums[right]) {
            int mid = left + ((right-left) >> 1);
            if (nums[mid] == nums[left]) {
                left++;
            }
            else if (nums[mid] < nums[left]){
                right = mid;
            }
            else {
                left = mid + 1;
            }
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
                    .nums = {1,3,5},
                    .expect = 1,
            },
            {
                    .nums = {2,2,2,0,1},
                    .expect = 0,
            },
            { .nums = {3,3,1,3}, .expect = 1},
            { .nums = {1,3,3}, .expect = 1},
            { .nums = {-1,-1,-1,-1}, .expect = -1},
            { .nums = {2,0,1,1,1}, .expect = 0},
    };

    {
        Solution1 s;
        cout << "== S1 ==" << endl;
        for (auto &test : tests) {
            auto result = s.findMin(test.nums);
            cout << "nums=" << test.nums << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    {
        Solution2 s;
        cout << "== S2 ==" << endl;
        for (auto &test : tests) {
            auto result = s.findMin(test.nums);
            cout << "nums=" << test.nums << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}


