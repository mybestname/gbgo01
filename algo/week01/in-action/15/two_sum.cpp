#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static vector<int> twoSum(vector<int> &nums, int target) { //输入无序
        // 使用一个pair，分别值和存储下标。
        vector<pair<int,int>> numbers;
        for (int i = 0 ; i< nums.size(); i++) {
            numbers.push_back(make_pair(nums[i],i));
        }
        // 排序
        sort(numbers.begin(), numbers.end());
        // 再按167处理
        int j = numbers.size()-1;
        for (int i = 0; i<numbers.size(); i++) { //i的外层枚举不变，目标是去掉内层循环
            while(i<j && numbers[i].first+numbers[j].first > target) j--;
            if (i<j && numbers[i].first+numbers[j].first == target)
                return {numbers[i].second,numbers[j].second};  //这里下标不用再加一了，因为题意是按0下标返回。
        }
        return {};
    }
};

int main() {
    vector<int> nums = {25,5,75};
    int target = 100;
    vector<int> result = Solution::twoSum(nums, target);
    cout << "nums=" << nums << ", target=" << target <<", result=" << result << endl;
    // [25, 5, 75], 100
    // [0,2]
    return 0;
}