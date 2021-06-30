#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static vector<vector<int>> threeSum(vector<int>& nums) {
        // 先排序，因为不需要考虑下标，直接排序即可
        sort(nums.begin(),nums.end());
        // 设 i < j < k
        // nums[i]+nums[j]+nums[k] = 0
        // nums[j]+nums[k] = -nums[i]
        // 等效于target==-nums[i]的two_sum
        vector<vector<int>> ans;
        for (int i = 0; i< nums.size(); i++ ){  // 外层枚举所有的target，内层调用two sum

            // 特别注意，输入数据需要去重！！！nums[i]代表target，如果有重复的target，那么就会有重复的答案。
            if (i>0 && nums[i] == nums[i-1]) continue; // 如果target和上一个重复，则skip

            // 注意，三个数不能重复。又 i<j<k, 那么从i+1开始搜索j，k。
            auto all_jks = twoSum(nums, i+1, -nums[i]);
            for (auto jk : all_jks) {
               ans.push_back({nums[i], jk[0],jk[1]});
            }
        };
        return ans;
    }
private:
    static vector<vector<int>> twoSum(vector<int>& numbers, int start, int target) { //设定一个start值，而不是从0开始
        vector<vector<int>> ans; //找所有的解
        int j = numbers.size()-1;
        for (int i = start; i<numbers.size(); i++) { // 从start开始寻找
            // ！！！！另外需要特别注意，这里找所有的值，而不是下标，那么有可能下标不同，但是值是相同的。
            if (i>start && numbers[i] == numbers [i-1]) continue;
            // numbers[i-1]代表上一个可能的返回值，只要和这个值相同，那么不管是不是，都不用继续了。
            // 如果是，答案已经存过了，如果不是，也不用再重复判断一遍了。
            while(i<j && numbers[i]+numbers[j] > target) j--;
            if (i<j && numbers[i]+numbers[j] == target) {
                ans.push_back({numbers[i],numbers[j]});   //存值，而不是下标。
            }
        }
        return ans;
    }
};

int main() {

    //intput [0,0,0,0]
    //output [[0,0,0]]
    vector<int> nums = {0,0,0,0};
    vector<vector<int>>  result = Solution::threeSum(nums);
    cout << "nums=" << nums << ", result=" << result << endl;
    return 0;

    //Input: nums = [-1,0,1,2,-1,-4]
    nums = {-1,0,1,2,-1,-4};
    result = Solution::threeSum(nums);
    cout << "nums=" << nums << ", result=" << result << endl;
    //Output: [[-1,-1,2],[-1,0,1]]
}


