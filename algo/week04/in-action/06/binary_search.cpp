#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"
using namespace std;

// 原解法，使用的其实是模版2
// - 两侧都不包含，用 ans 维护答案，终止于 left > right
// - 需要在mid被排除掉之前进行判断，如果mid是一个解，更新一下ans。
class Solution {
public:
    int search(vector<int>& nums, int target) {
        int ans = -1, left = 0, right = nums.size() -1;
        while ( left <= right ) {
            int mid = (left + right)/2;
            if (nums[mid] == target) {
                ans = mid;
                break;
            }
            if (nums[mid] < target) {
                left = mid + 1;
            }else {
                right = mid - 1;
            }
        }
        return ans;
    }
};

// 使用模版1.1
class Solution1_1 {
public:
    int search(vector<int>& nums, int target) {
        int left = 0, right = nums.size() ;
        while ( left < right ) {
            int mid = (left + right ) >> 1;
            if (nums[mid] == target) {
                return mid;
            }
            if (nums[mid] < target) {
                right = mid;
            }else {
                left = mid + 1;
            }
        }
        //return right;
        return -1; // not found
    }
};

// 使用模版1.2
class Solution1_2 {
public:
    int search(vector<int>& nums, int target) {
        int left = -1, right = nums.size() -1;
        while ( left < right ) {
            int mid = (left + right + 1) >> 1;
            if (nums[mid] == target) {
                return mid;
            }
            if (nums[mid] < target) {
                left = mid;
            }else {
                right = mid - 1;
            }
        }
        return left;
    }
};

int main() {
    struct Test {
        vector<int> nums;
        int target;
        int expect;
    };
    vector<Test> tests = {
            {
                    .nums   = {-1,0,3,5,9,12},
                    .target =  9,
                    .expect =  4,
            },
            {
                    .nums   = {-1,0,3,5,9,12},
                    .target =  2,
                    .expect =  -1,
            },
    };
    {
        Solution s;
        cout << "== 原解法 ==" << endl;
        for (auto &test : tests) {
            auto result = s.search(test.nums, test.target);
            cout << "nums=" << test.nums << ",target=" << test.target << ",expect=" << test.expect <<",got=" << result << endl;
        }
    }
    {
        Solution1_1 s;
        cout << "== 模版1.1 ==" << endl;
        for (auto &test : tests) {
            auto result = s.search(test.nums, test.target);
            cout << "nums=" << test.nums << ",target=" << test.target << ",expect=" << test.expect <<",got=" << result << endl;
        }
    }
    {
        Solution1_2 s;
        cout << "== 模版1.2 ==" << endl;
        for (auto &test : tests) {
            auto result = s.search(test.nums, test.target);
            cout << "nums=" << test.nums << ",target=" << test.target << ",expect=" << test.expect <<",got=" << result << endl;
        }
    }

    return 0;
}
