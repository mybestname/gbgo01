#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;

class Solution {
public:
    static vector<string> subdomainVisits(vector<string>& cpdomains) {
        unordered_map<string, int> domain_count;
        for (auto cpdomain : cpdomains) {
            // 格式为访问次数+空格+地址
            int space_i = cpdomain.find(" ");
            int count = stoi(cpdomain.substr(0,space_i));
            string addr = cpdomain.substr(space_i+1);
            domain_count[addr]+=count;
            for (int i = 0; i< addr.length(); i++){
                if ( addr[i] != '.' ) continue;
                domain_count[addr.substr(i+1)]+=count;
            }
        }
        vector<string> visits;
        for (const auto& count : domain_count) {
            visits.push_back(to_string(count.second)+" "+count.first);
        }
        return visits;
    }
};

int main() {
    vector<string> domains = {"9001 discuss.leetcode.com"};
    auto result = Solution::subdomainVisits(domains);
    cout << "domains=" << domains << ",visit=" << result << endl;
    //
    domains = {"900 google.mail.com", "50 yahoo.com", "1 intel.mail.com", "5 wiki.org"};
    result = Solution::subdomainVisits(domains);
    cout << "domains=" << domains << ",visit=" << result << endl;
    return 0;
}

