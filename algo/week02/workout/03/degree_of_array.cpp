#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    static int findShortestSubArray(vector<int>& nums) {
        unordered_map<int, vector<int>> frequency;  // key is num, v is {num_count,first_index,last_index}
        int max_frequency = 0;
        for (int i= 0; i< nums.size(); i++){
            int num = nums[i];
            if (frequency.count(num) == 0) {
                frequency[num] = {1,i,i};
            }else{
                frequency[num][0]++;  // increase num_count
                frequency[num][2]=i;  // update last_index
            }
            max_frequency = max(max_frequency,frequency[num][0]); //update max_frequency
        }
        int shortest_len = numeric_limits<int>::max();
        for (auto & it : frequency){
            if (max_frequency == it.second[0]) {
               shortest_len = min(shortest_len, it.second[2]-it.second[1]+1);
            }
        }
        return shortest_len;
    }
};

int main() {
    // 输入：[1, 2, 2, 3, 1]
    // 输出：2
    vector<int> nums = {1, 2, 2, 3, 1};
    auto result = Solution::findShortestSubArray(nums);
    cout << "nums=" << nums << ", shortest_subarray=" << result << endl;

    // 输入：[1,2,2,3,1,4,2]
    // 输出：6
    nums = {1,2,2,3,1,4,2};
    result = Solution::findShortestSubArray(nums);
    cout << "nums=" << nums << ", shortest_subarray=" << result << endl;

    return 0;
}
