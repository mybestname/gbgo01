#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<vector<int>> subsets(vector<int>& nums) {
        findSubset(nums,0);
        return ans;
    }
private:
    vector<vector<int>> ans;
    vector<int> s;
    void findSubset(vector<int>& nums, int index){
        //终止条件
        if (index == nums.size()) {
            ans.push_back(s);
            return;
        }
        // 枚举nums[0],nums[1],nums[2], ... nums[n]，这n个数，选或者不选。
        // 每次都是两种情况
        // 不选
        //   对全局的状态没有任何影响。
        findSubset(nums,index+1);
        // 选
        s.push_back(nums[index]); //把数放到集合里面去。
        findSubset(nums,index+1); //然后再调用下一个, 注意：这两次调用，面对的全局变量s是不一样的。这也是为何第二次的s要恢复的原因。
        // 还原全局变量
        s.pop_back();
    }
};
int main() {
    vector<int> nums = {1,2,3};
    Solution s;
    vector<vector<int>> subsets = s.subsets(nums);
    // 输入：nums = [1,2,3]
    // 输出：[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
    cout << "nums=" << nums << ", subset=" << subsets << endl;
    return 0;
}

