#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<vector<int>> permute(vector<int>& numbers) {
        nums = numbers;
        n = numbers.size();
        for (int i=0; i<n; i++) {
            used.push_back(false);
        }
        find(0);
        return ans;
    }
    // 考虑0..n-1个位置
    void find(int index) {
        //结束条件
        if (index == n) {
           //cout << " s = " << s << endl;
           ans.push_back(s);
           return;
        }
        for (int i = 0; i< n; i++) {
            if (!used[i]) {
                used[i] = true;
                s.push_back(nums[i]);
                //下一个
                find(index+1);   //注意，每一次的调用，其全局变量used和s是不同的，这也是调用完要恢复的原因。
                //恢复状态
                used[i] = false;
                s.pop_back();
            }
        }
    }

private:
    vector<vector<int>> ans;
    vector<int> s;
    int n{};
    vector<int> nums;
    vector<bool> used;
};

int main() {
    Solution s;
    vector<int> nums = {1,2,3};
    auto result = s.permute(nums);
    cout << "nums=" << nums <<",permutations="<<result << endl;
    // 输入：nums = [1,2,3]
    // 输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
    return 0;
}