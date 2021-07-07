#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<vector<int>> combine(int n, int k) {
        this->n = n;
        this->k = k;
        findSubset(1); //从第一开始，而不是从0个开始，因为
        return ans;
    }
private:
    vector<vector<int>> ans;
    vector<int> s;
    int n{},k{};
    //枚举n个数，1,2,3...n，选还是不选
    void findSubset(int index) {
        // cout << "findSubset :" << index << ", s="<< s << endl;
        // 终止条件
        // 如果subset的长度已经超过k个，或者剩下的肯定无法凑足k个。那么肯定不满足，退出
        if (s.size()>k || s.size()+n-index+1 <k) return;
        if (index == n+1) {  //已经到n+1，结束。
            //cout << "ans add " << s <<endl;
            ans.push_back(s);
            return;
        }
        // 处理
        // 不选，直接下一个。
        findSubset(index+1);
        // 或选到subset
        s.push_back(index);
        findSubset(index+1);
        // 恢复状态
        s.pop_back();
    }
};

int main() {
    Solution s;
    int n=4, k=2;
    auto result = s.combine(n,k);
    cout << "n=" <<n << ",k=" << k <<",combinations="<<result << endl;
    return 0;
}
