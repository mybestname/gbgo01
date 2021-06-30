#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int numberOfNiceSubarrays(vector<int>& nums, int k) {
        // 建立下标为1..n的s（前缀和）数组和count数组（原数据转化位0/1的数组）
        uint n = nums.size();
        vector<int> s(n+1, 0);      // 0..n 前缀和数组s
        vector<int> counts(n+1,0);  // 0..n 记录前缀和统计数据，表示s中存在的数的数量。
        for (int i = 1;  i <= n ; i++ ) { // 1..n 所以直接可以用i-1
           s[i] = s[i-1] + nums[i-1]%2;   // 填充前缀和数组，注意使用 % 2，使得nums中的数据被过滤为0和1
        }
        for (int i = 0; i <=n ; i++) {
            counts[s[i]]++;                // 统计前缀和中元素的数量
        }
        // 例：输入为  [1,1,2,1,1]
        //           [1,1,0,1,1]
        //      s为[0,1,2,2,3,4]
        //  count为[1,1,2,1,1] 解读为前缀和为0的数量为1, 前缀和1的数量为1，前缀和2为数量为2，前缀和为3的数量为2
        //
        // 问题实质是求：对于每一个i，前缀和为s[i]-k的数量有多少个？注意s[i]-k应该大于等于0
        int ans = 0;
        for (int i=1; i<=n; i++){
            int m = s[i]-k;
            // cout << "m=" << m  << ", s[i]=" << s[i] << endl;
            if (m >= 0) {
                ans += counts[m];
                //cout << "ans=" << ans <<",m=" <<m <<",counts[m]="<<counts[m]<<",i="<<i<< ",counts=" << counts << ",s=" << s << endl;
            }
        }
        return ans;
        // 例：
        //      s [0,1,2,2,3,4]
        //  count [1,1, 2, 1,1]
        //      k=3
        // s[1]-2 = 1-3 = -2 (drop)
        // s[2]-2 = 2-3 = -1 (drop)
        // s[3]-2 = 2-3 = -1 (drop)
        // s[4]-2 = 3-3 = 0
        // s[5]-2 = 4-3 = 1
        // count[0] = 1
        // count[1] = 1
        // answer = 1 + 1 = 2
    }
};

int main() {
    vector<int> nums = {1,1,2,1,1};
    int k = 3;
    int result = Solution::numberOfNiceSubarrays(nums, k);          // nums = [1,1,2,1,1], k = 3
    cout <<  "nums=" << nums  << ",k=" << k << "; nice=" << result << endl; // 2
    return 0;
}