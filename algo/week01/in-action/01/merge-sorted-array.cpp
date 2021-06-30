#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
    public:
    void static merge(vector<int>& nums1, int m, const vector<int>& nums2, int n) {
        //首先建两个索引i和j，那么考虑大小，从后往前。
        int i = m - 1, j = n - 1;
        // 一共需要 m + n 个元素
        for (int k = m + n - 1 ; k >= 0; k--) {
            if ( i<0 || (j >=0 && nums1[i] < nums2[j])) {
                nums1[k] = nums2[j];
                j--;
            } else {
                nums1[k] = nums1[i];
                i--;
            }
        }
    }
};

int main() {
    vector<int> nums1 = {1,2,3,0,0,0};
    vector<int> nums2 = {2,5,6};
    int m = 3;
    int n = 3;
    Solution::merge(nums1,m,nums2,n);
    cout << nums1 << endl; // [1,2,2,3,5,6]
    return 0;
}
