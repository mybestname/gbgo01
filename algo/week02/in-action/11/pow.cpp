#include <iostream>
#include <vector>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    double myPow(double x, int n) {
        if (n<0) return 1/ myPow(x,-n); //负数转化为正数 //注意，这里隐含着问题，n和-n的范围不对称。
        // 递归边界/终止条件
        if (n==0) return 1;
        if (n==1) return x;
        double temp = myPow(x,n/2);
        if (n%2 == 0) return temp*temp;
        else return temp*temp*x;
        // 因为没有任何共享全局变量，所以不需要恢复状态。
    }
};

int main(){
    struct Test {
        float x;
        int n;
        float expect;
    };
    vector<Test> tests = {
            {.x = 2.00000, .n = 10, .expect = 1024.00000},
            {.x = 2.10000, .n = 3,  .expect = 9.26100},
            {.x = 2.00000, .n = -2, .expect = 0.25000},
    };
    Solution s;
    for (auto& test : tests) {
        auto result = s.myPow(test.x, test.n);
        cout << "x=" << test.x << ", n=" << test.n << ", expect=" << test.expect <<", got=" <<result<< endl;
    }
}

