#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<string> letterCombinations(string digits) {
        edges['2'] ="abc";
        edges['3'] ="def";
        edges['4'] ="ghi";
        edges['5'] ="jkl";
        edges['6'] ="mno";
        edges['7'] ="pqrs";
        edges['8'] ="tuv";
        edges['9'] ="wxyz";
        ans = {};
        dfs(digits, 0);
        return ans;
    }
private:
    unordered_map<char,string> edges; //其实是一个出边数组
    vector<string> ans;
    string s;
    void dfs(string& digits, uint level){
        // 结束条件
        if (level == digits.size()) {
            ans.push_back(s);
            return;
        }
        // 考虑所有的出边
        for (char& ch : edges[digits[level]]){
            s.push_back(ch);
            dfs(digits, level+1);
            s.pop_back();  //还原现场
        }
    }
};

int main(){
    struct Test {
        string digits;
        vector<string> expect;
    };
    vector<Test> tests = {
            {.digits = "23", .expect={"ad","ae","af","bd","be","bf","cd","ce","cf"}},
            {.digits = "", .expect={} },
            {.digits = "2",.expect={"a","b","c" }},
    };
    {
        Solution s;
        for (auto &test : tests) {
            auto result = s.letterCombinations(test.digits);
            cout << " digits=\"" << test.digits << "\",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}

