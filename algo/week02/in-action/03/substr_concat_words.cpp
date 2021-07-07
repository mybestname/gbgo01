#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<int> findSubstring(string s, vector<string>& words) {
        vector<int> ans;
        wordsMap = countWords(words);
        int substr_size = words.size()*words[0].length();  //all words in the same length.
        for (int start=0; start < s.length(); start++) {
            if (isSame(s.substr(start,substr_size),words)) {
                ans.push_back(start);
            }
        }
        return ans;
    }

    //优化1：剔除掉某些非word开始的子串
    vector<int> findSubstringV2(string s, vector<string>& words) {
        vector<int> ans;
        wordsMap = countWords(words);
        int word_size = words[0].length();  //all words in the same length.
        int substr_size = words.size()*word_size;
        for (int start=0; start < s.length(); start++) {
            string a_word = s.substr(start,word_size);
            if (wordsMap.find(a_word) != wordsMap.end()) { // should start from a word
                if (isSame(s.substr(start, substr_size), words)) {
                    ans.push_back(start);
                }
            }
        }
        return ans;
    }
    // 优化2：
    // 没有必要每次重新构造tMap，可以删除一个word，再加入一个Word
    vector<int> findSubstringV3(string s, vector<string>& words) {
        vector<int> ans;
        wordsMap = countWords(words);
        int word_size = words[0].length();  //all words in the same length.
        int substr_size = words.size()*word_size;
        for (int first = 0; first < word_size; first++) {
            if (first + substr_size > s.size()) break;
            unordered_map<string,int> tMap;
            int curr = first;
            for (int i=0; i < words.size(); i++) {
                tMap[s.substr(curr, word_size)]++;
                curr+= word_size;
            }
            for (int start= first, end=curr; start+substr_size <= s.size(); start+=word_size, end+=word_size) {
                if (isMapEqual(tMap, wordsMap)) ans.push_back(start);
                tMap[s.substr(end,word_size)]++;
                tMap[s.substr(start,word_size)]--;
                if (tMap[s.substr(start,word_size)] == 0) {
                    tMap.erase(tMap.find(s.substr(start,word_size)));
                }
            }

        }
        return ans;
    }

private:
    unordered_map<string, int> wordsMap; // map(k=单词,v=出现次数)

    bool isSame(string t, vector<string>& words) {
        int m = words[0].length(); //所有单词长度相同。
        unordered_map<string,int> tMap;
        for (int i =0; i < t.length(); i+=m) {
            // i开始m个字符为一个word
            tMap[t.substr(i,m)]++;
        }
        return isMapEqual(tMap, wordsMap);
    }

    static bool isMapEqual(unordered_map<string, int>& a, unordered_map<string, int>& b) {
        // the 2 maps are in the same size
        if (a.size() != b.size()) return false;
        // (each em in a exists in b)
        for (auto& kv : a) {
            if (b.find(kv.first) == b.end() || b[kv.first] != kv.second) return false;
        }
        return true;
    }

    static unordered_map<string,int> countWords(vector<string>& words) {
        unordered_map<string, int> ans;
        for (string& word: words) {
            ans[word]++;
        }
        return ans;
    }
};

int main() {
    struct Test {
        string s;
        vector<string> words;
        vector<int> expect;
    };
    vector<Test> tests = {
            {.s = "barfoothefoobarman", .words =  {"foo","bar"}, .expect = {0,9}},
            {.s = "wordgoodgoodgoodbestword", .words =  {"word","good","best","word"}, .expect = {}},
            {.s = "barfoofoobarthefoobarman", .words =  {"bar","foo","the"}, .expect = {6,9,12}},
            {.s = "lingmindraboofooowingdingbarrwingmonkeypoundcake",.words={"fooo","barr","wing","ding","wing"},.expect = {13}}
    };
    Solution s;
    for (auto& test : tests) {
        vector<int> result = s.findSubstring(test.s, test.words);
        cout << "s=" << test.s << ", words=" << test.words << ", expect=" << test.expect <<",got=" <<result<< endl;
    }

    for (auto& test : tests) {
        vector<int> result = s.findSubstringV2(test.s, test.words);
        cout << "s=" << test.s << ", words=" << test.words << ", expect=" << test.expect <<",got=" <<result<< endl;
    }

    for (auto& test : tests) {
        vector<int> result = s.findSubstringV3(test.s, test.words);
        cout << "s=" << test.s << ", words=" << test.words << ", expect=" << test.expect <<",got=" <<result<< endl;
    }

}