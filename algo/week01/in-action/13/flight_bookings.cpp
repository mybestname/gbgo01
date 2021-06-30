#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static vector<int> corpFlightBookings(vector<vector<int>>& bookings, int n) {
        vector<int> delta(n+2, 0);  //差分数组，0～n+1
        for (auto& booking : bookings) {
            int first = booking[0];
            int last = booking[1];
            int seats = booking[2];
            // 差分公式 `B[l] += d; B[r+1] -= d`
            delta[first] += seats;
            delta[last+1] -= seats;
        }
        vector<int> a(n+1,0); //原数组 0~n（多一位）
        // 1~n 对差分求前缀和，得到原数组
        for (int i = 1; i <= n; i++) a[i] = delta[i]+a[i-1];
        // 向前移动一位，满足0～n-1，以满足返回条件
        for (int i = 1; i <= n; i++) a[i-1] = a[i];
        a.pop_back();
        return a;
    }
};

int main() {
    vector<vector<int>> bookings = {
            {1, 2, 10},
            {2, 3, 20},
            {2, 5, 25},};
    // 输入：bookings = [[1,2,10],[2,3,20],[2,5,25]], n = 5
    // 输出：[10,55,45,25,25]
    int n = 5;
    vector<int> result = Solution::corpFlightBookings(bookings, n);
    cout << "bookings=" << bookings << ", result=" << result << endl;
    return 0;
}

