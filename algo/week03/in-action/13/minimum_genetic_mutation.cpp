#include <iostream>
#include <vector>
#include <queue>
#include <unordered_map>
#include <unordered_set>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    int minMutation(const string& start, const string& end, vector<string>& bank) {
        const char gene[4] = {'A','C','G','T'};
        unordered_set<string> bank_set;
        for (const auto& s: bank) {
           bank_set.insert(s);
        }
        unordered_map<string, int> depth;
        depth[start] = 0;
        queue<string> q;
        q.push(start);
        while (!q.empty()) {
            // 取出队头
            string s = q.front();
            q.pop();
            // 遍历 24条出边 (8种位置，3种可能选项）
            for (int i=0; i< 8 ; i++) {
                for (char j : gene) {
                    if (s[i] == j) continue;
                    string ns = s;
                    ns[i] = j;
                    if (bank_set.find(ns) == bank_set.end()) continue;
                    //考虑ns是否走过
                    if (depth.find(ns)==depth.end()){
                        depth[ns] = depth[s]+1 ;  //ns的深度+1 s->ns，从s走到了ns
                        q.push(ns); //ns入队
                        if (ns == end) {//提前结束
                            return depth[ns];
                        }
                    }
                }
            }
        }
        return -1; //没有找到
    }
};

int main(){
    struct Test {
        string start;
        string end;
        vector<string> bank;
        int expect ;
    };
    vector<Test> tests = {
            {
                    .start = "AACCGGTT",
                    .end =   "AACCGGTA",
                    .bank = {"AACCGGTA"},
                    .expect = 1,
            },
            {
                    .start= "AACCGGTT",
                    .end  = "AAACGGTA",
                    .bank = {"AACCGGTA", "AACCGCTA", "AAACGGTA"},
                    .expect = 2,
            },
            {
                    .start = "AAAAACCC",
                    .end =   "AACCCCCC",
                    .bank = {"AAAACCCC", "AAACCCCC", "AACCCCCC"},
                    .expect = 3,
            },
    };
    {
        Solution s;
        for (auto &test : tests) {
            auto result = s.minMutation(test.start, test.end, test.bank);
            cout << " start=" << test.start<< ",end="<<test.end << ",bank="<< test.bank
                 << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}


