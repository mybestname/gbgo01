#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int maxArea(vector<int> &height) {
        //首先考虑暴力法
        // 1. i < j
        // 2. for i=0 to n-1
        // 3.   for j = i+1 to n-1
        // 4.       ans = max(ans, area(i,j))
        return 0;
    }
};

int main() {
    vector<int> nums = {1,8,6,2,5,4,8,3,7};
    int result = Solution::maxArea(nums);
    cout << "nums=" << nums << ",result=" << result << endl;
    //输入：[1,8,6,2,5,4,8,3,7]
    //输出：49
    return 0;
}
