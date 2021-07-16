#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    vector<vector<string>> solveNQueens(int n) {
        if (n <= 0) return {};
        N = n;
        used = vector<bool>(n,false);
        find(0);
        vector<vector<string>> result;
        for (auto& p : ans) {
            vector<string> r;
            for (int row = 0; row < n; row++) {
                int col = p[row];
                string s(n, '.');
                s[col] = 'Q';
                r.push_back(s);
            }
            result.push_back(r);
        }
        return result;
    }

    // 考虑0..n-1个位置
    void find(int row) {
        //结束条件
        if (row == N) {
            ans.push_back(per);
            return;
        }
        for (int col = 0; col< N; col++) {
            if (!used[col] && !usedCPlusR[col+row] && !usedCMinusR[col-row]) {
                used[col] = true;
                usedCPlusR[col+row] = true;
                usedCMinusR[col-row] = true;
                per.push_back(col);
                find(row+1);
                used[col] = false;
                usedCPlusR[col+row] = false;
                usedCMinusR[col-row] = false;
                per.pop_back();
            }
        }
    }

private:
    vector<vector<int>> ans;
    vector<int> per;   //a permutation of numbers
    int N{};
    vector<bool> used;
    unordered_map<int,bool> usedCPlusR;  // col + row used
    unordered_map<int,bool> usedCMinusR; // col - row used
};

int main(){
    struct Test {
        int n;
        vector<vector<string>> expect;
    };
    vector<Test> tests = {
            {.n = 4, .expect={{".Q..","...Q","Q...","..Q."},{"..Q.","Q...","...Q",".Q.."}}},
    };
    {
        Solution s;
        for (auto &test : tests) {
            auto result = s.solveNQueens(test.n);
            cout << " n=" << test.n<< ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}

