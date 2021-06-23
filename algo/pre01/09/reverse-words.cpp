#include <iostream>
#include <vector>

using namespace std;
class Solution {
public :
    static string reverseWords(string s) {
        s.erase (3,4);
        return s;
    }
};

int main() {
    Solution sol;
    string input = "the sky is blue";
    string words = Solution::reverseWords(input);
    cout << words << endl; //4

}
