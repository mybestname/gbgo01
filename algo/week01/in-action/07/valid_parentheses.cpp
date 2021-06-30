#include <iostream>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static bool isValid(string s) {
        stack<int> st;
        for (char ch : s) {
            if (ch == '(') st.push(')');
            else if (ch == '[') st.push(']');
            else if (ch == '{') st.push('}');
            else if (!st.empty() && st.top() == ch ) st.pop();
            else return false;
        }
        return st.empty();
    }
};

int main() {
    string s = "()[]{}";
    bool valid = Solution::isValid(s);
    cout << s << ",valid=" << valid << endl;
    s = "([)]\"";
    valid = Solution::isValid(s);
    cout << s << ",valid=" << valid << endl;
    return 0;
}
