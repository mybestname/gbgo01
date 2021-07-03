#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    // 使用前缀最值法
    int trap(vector<int>& height) {
        int n = height.size();
        prefix[0] = suffix[n+1] = 0;
        for (int i = 1; i <= n; i++ ) {
            prefix[i] = max(prefix[i-1], height[i-1]); //前缀最值
        }
        for (int i = n; i > 0; i--) {
            suffix[i] = max(suffix[i+1], height[i-1]); //后缀最值
        }
        int ans = 0;
        for (int i = 1; i <=n ; i++) {
            ans += max(0, min(prefix[i-1], suffix[i+1])-height[i-1]); //前缀最值和后缀最值中的小者减去自己的高度，为该单元格上方雨水的最值。
        }
       return ans;
    }

    // 使用单调栈法
    int trapV2(vector<int>& height) {
        int ans = 0;
        stack<Rect> s;
        s.push({0,0});
        for (int h : height) {
            int w = 0;
            while (s.size() > 1 && s.top().height <= h) {   //递增
                w += s.top().width;
                int bottom = s.top().height;
                s.pop();
                ans += w*max(0,min(s.top().height,h)-bottom);
            }
            s.push({h,w+1});
        }
        return ans;
    }
private:
    int prefix[100];
    int suffix[100];

    struct Rect {
        int height;
        int width;
    };
};

int main() {
    vector<int> heights = {0,1,0,2,1,0,1,3,2,1,2,1};
    Solution s;
    int result = s.trap(heights);
    cout << "heights=" << heights << ",result=" << result << endl;
    // 输入：height = [0,1,0,2,1,0,1,3,2,1,2,1]
    // 输出：6
    {
        vector<int> heights = {0,1,0,2,1,0,1,3,2,1,2,1};
        Solution s;
        int result = s.trapV2(heights);
        cout << "heights=" << heights << ",result=" << result << endl;
        // 输入：height = [0,1,0,2,1,0,1,3,2,1,2,1]
        // 输出：6
    }

    return 0;
}

