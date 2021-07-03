#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int largestRectangleArea(vector<int>& heights) {
        stack<Rect> s;
        int ans = 0;
        heights.push_back(0); // 保证stack可以被清空。
        for (int h : heights) {
            int accumulated_width = 0;
            while(!s.empty() && s.top().height >= h) { //如果栈中的矩形的比当前高度高，说明当前为递减
                accumulated_width += s.top().width;
                ans = max(ans, accumulated_width * s.top().height); //计算面积
                s.pop(); //一直出栈，直到变成递增的情况。
            }
            // 开始递增，每个递增矩形入栈
            s.push({h,accumulated_width+1});
        }
        // 注意，这个解法只有在pop的时候，才会计算面积，那么对于结尾部分情况是单调递增，有可能stack没有全部清空
        // 简单的办法是在队尾手工添加一个0，用来保证stack可以被全部pop。
        return ans;
    }

private:
    struct Rect {
        int height;
        int width;
    };
};

int main() {
    vector<int> heights = {2,1,5,6,2,3};
    int area = Solution::largestRectangleArea(heights);
    cout << "heights=" << heights << ",area=" << area << endl;
    //输入: [2,1,5,6,2,3]
    //输出: 10

    // 全部递增的情况
    heights = {2,5,6};
    area = Solution::largestRectangleArea(heights);
    cout << "heights=" << heights << ",area=" << area << endl;
    return 0;
}
