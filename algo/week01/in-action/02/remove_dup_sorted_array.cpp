#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    int static removeDuplicates(vector<int>& nums) {
        int j = 1 ; //元素索引，有条件增加，最后返回。从1开始。少做一次
        // 空间索引i，最大不超过num的size, 无条件增加。可以从1开始，这样就不用处理i=0时候，i-1的判定。
        // 同时少做一次`nums[j] = nums[i]`
        for (int i = 1; i < nums.size() ; i++)  {
            if ( nums[i] != nums[i-1]) {
                nums[j] = nums[i];
                j++;
            }
        }
        // 比较 0 开始
        // int j = 0;
        // for (int i = 0; i < nums.size() ; i++)  {
        //     if ( i==0 || nums[i] != nums[i-1]) {
        //         nums[j] = nums[i];
        //         j++;
        //     }
        // }

        return j;
    }

};

int main() {
    vector<int> nums = {1,1,2};

    int n = Solution::removeDuplicates(nums);
    nums.resize(n);
    cout << n << " " << nums << endl; // 2 [1,2]

    nums = {0,0,1,1,1,2,2,3,3,4};
    n = Solution::removeDuplicates(nums);
    nums.resize(n);
    cout << n << " " << nums << endl; // 5 [0,1,2,3,4]
    return 0;
}


