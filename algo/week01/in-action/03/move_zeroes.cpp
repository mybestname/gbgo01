#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    void static moveZeroes(vector<int>& nums) {
        int n = 0 ;
        for (int i = 0; i< nums.size(); i++) {
            if (nums[i] != 0) {
                nums[n] = nums[i];
                n++;
            }
        }
        for (int i = n; i< nums.size(); i++) {
            nums[i]=0;
        }
    }
};

int main() {
    vector<int> nums = {0,1,0,3,12};  // [0,1,0,3,12]
    Solution::moveZeroes(nums);
    cout << nums << endl; //  [1,3,12,0,0]
    return 0;
}

