#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int maxArea(vector<int> &height) {
        int i = 0, j = height.size()-1;
        int ans = 0;
        while (i < j) {
            if (height[i]<height[j]) {
                // 短边是i
                ans = max(ans, (j-i) * height[i]);
                i++;
            }else {
                // 短边是j
                ans = max(ans, (j-i) * height[j]);
                j--;
            }
        }
        return ans;
    }
    //解法2，小优化，如果height[i]==height[j]，则可以同时移动i,j，即两个边都可抛弃。
    static int maxArea2(vector<int> &height) {
        int i = 0, j = height.size()-1;
        int ans = 0;
        while (i < j) {
            ans = max (ans, min(height[i],height[j])*(j-i));
            if (height[i]==height[j]) i++, j--;
            else if (height[i]<height[j]) i++; else j--;
        }
        return ans;
    }
};

int main() {
    vector<int> nums = {1,8,6,2,5,4,8,3,7};
    int result = Solution::maxArea(nums);
    cout << "nums=" << nums << ",result=" << result << endl;
    //输入：[1,8,6,2,5,4,8,3,7]
    //输出：49
    result = Solution::maxArea2(nums);
    cout << "nums=" << nums << ",result=" << result << endl;
    return 0;
}
