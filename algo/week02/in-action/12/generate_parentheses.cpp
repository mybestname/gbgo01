#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<string> generateParenthesis(int n) {
        // 边界条件/递归终止
        if (n == 0) return {""};

        // 拆分方法：(a)b
        // 示例：
        // ((())) a="(())" b=""
        // (())() a="()"   b="()"
        // ()()() a=""     b="()()"
        //
        // (a) ：k对括号，子问题a为k-1对括号
        // b   ：n-k对括号
        vector<string> result;
        // 不同的k之间：加法原理
        for (int k= 1; k<= n; k++) {
            vector<string> result_a = generateParenthesis(k-1);
            vector<string> result_b = generateParenthesis(n-k);
            // 乘法原理
            for (string& a : result_a) {
                for (string& b : result_b){
                    result.push_back( "(" + a + ")" + b);
                }
            }
        }
        return result;
    }
};

// 优化: 减少重复计算
class Solution2 {
public:
    vector<string> generateParenthesis(int n) {
        if (n == 0) return {""};
        // 直接返回中之前已经算过的结果 (因为针对一个固定的n，结果是确定的）
        if (cache.find(n) != cache.end()) return cache[n];
        vector<string> result;
        for (int k= 1; k<= n; k++) {
            vector<string> result_a = generateParenthesis(k-1);
            vector<string> result_b = generateParenthesis(n-k);
            for (string& a : result_a) {
                for (string& b : result_b){
                    result.push_back( "(" + a + ")" + b);
                }
            }
        }
        cache[n]=result;
        return result;
    }
private:
    unordered_map<int,vector<string>> cache;
};



int main() {
    struct Test {
        int n;
        vector<string> expect;
    };
    vector<Test> tests = {
            { .n = 3, .expect = {"((()))","(()())","(())()","()(())","()()()"}},
            { .n = 1, .expect = {"()"}},
    };
    {
        Solution s;
        for (auto &test : tests) {
            auto result = s.generateParenthesis(test.n);
            cout << "n=" << test.n << ", expect=" << test.expect << ", got=" << result << endl;
        }
    }
    {
        Solution2 s;
        for (auto &test : tests) {
            auto result = s.generateParenthesis(test.n);
            cout << "n=" << test.n << ", expect=" << test.expect << ", got=" << result << endl;
        }
    }
}


