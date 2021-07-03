#include <iostream>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <array>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static vector<vector<string>> groupAnagrams(vector<string>& strs) {
        unordered_map<string,vector<string>> group;
        for (const auto& str : strs) {
            string copy = str;
            sort(copy.begin(),copy.end());
            // find key from group
            //
            // if (group.find(copy)!= group.end()) {
            //     group[copy].push_back(str);
            // }else {
            //     group[copy] = {str};
            // }
            group[copy].push_back(str); //for C++ API, redundancy of the if-exist test, above 5 line reduce to 1 line.
        }
        vector<vector<string>> result;
        for (const auto& pair : group){
            result.push_back(pair.second);
        }
        return result;
    }
};

using namespace std;

class ArrayHash {
public:
    size_t operator()(const array<int, 26>& alphabet) const
    {
        size_t h = 0;
        for (auto e : alphabet) {
            h ^= std::hash<int>{}(e)  + 0x9e3779b9 + (h << 6) + (h >> 2);
            // phi = (1 + sqrt(5)) / 2
            // 2^32 / phi = 0x9e3779b9
            // the golden ratio
            // http://burtleburtle.net/bob/hash/doobs.html
            // https://stackoverflow.com/questions/35985960/c-why-is-boosthash-combine-the-best-way-to-combine-hash-values
        }
        return h;
    }
};

class Solution2 {
public:
    static vector<vector<string>> groupAnagrams(vector<string>& strs) {
        unordered_map<array<int,26>,vector<string>, ArrayHash> group;
        for (const auto& str : strs) {
            array<int,26> a{0};
            for (const auto& c : str) {
                a[tolower(c)-'a']++;
            }
            group[a].push_back(str);
        }
        vector<vector<string>> result;
        result.reserve(group.size());
        for (const auto& pair : group){
            result.push_back(pair.second);
        }
        return result;
    }
};

int main() {
    {
        vector<string> words = {"eat", "tea", "tan", "ate", "nat", "bat"};
        auto anagrams = Solution::groupAnagrams(words);
        // 输入: ["eat", "tea", "tan", "ate", "nat", "bat"]
        // 输出:
        // [
        // ["ate","eat","tea"],
        // ["nat","tan"],
        // ["bat"]
        // ]
        cout << "words=" << words << ", anagrams_group=" << anagrams << endl;
    }

    {
        vector<string> words = {"eat", "tea", "tan", "ate", "nat", "bat"};
        auto anagrams = Solution2::groupAnagrams(words);
        cout << "words=" << words << ", anagrams_group=" << anagrams << endl;
    }
    return 0;
}