#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"
using namespace std;

class Solution {
public:
    int mySqrt(int x) {
        // 寻找最大的x，满足 x^2 <= target
        // 使用long long避免越界。
        long long left = 0, right = 1<<16;  //2^16 x为32位
        while (left < right ) {
            long long mid = ( left + right + 1 ) >> 1;
            if (mid*mid <= x)
                left = mid;
            else
                right = mid-1;
        }
        return int(left);
    }
};

int main() {
    struct Test {
        int input;
        int expect;
    };
    vector<Test> tests = {
            {.input = 4, .expect = 2},
            {.input = 8, .expect = 2},
    };

    Solution s;
    for (auto &test : tests) {
        auto result = s.mySqrt(test.input);
        cout << "input=" << test.input << ",expect=" << test.expect <<",got=" << result << endl;
    }
    return 0;
}

