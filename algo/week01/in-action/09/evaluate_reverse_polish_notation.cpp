#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int evalRPN(vector<string> &tokens) {
        stack<int64_t> s;
        for (string& token : tokens){
           // is number
           if (token[0] >= '0' && token[0] <= '9' || token[0] == '-' ) {
               s.push(stoi(token));
           }else {
               int64_t b = s.top();  //注意栈顶是操作数b ，因为LIFO
               s.pop();
               int64_t a = s.top();
               s.pop();
               s.push(calc(a,b,token));
           }
        }
        return s.top();
    }

private:
    static int64_t calc(int64_t a, int64_t b, string op) {
        if (op == "+") return a+b;
        if (op == "-") return a-b;
        if (op == "*") return a*b;
        if (op == "/") return a/b;
        return 0;  //保证合法，不考虑。
    }
};

int main() {
    vector<string> tokens = {"10","6","9","3","+","-11","*","/","*","17","+","5","+"};
    int result = Solution::evalRPN(tokens);
    cout <<  tokens  << "=" << result << endl; // 22
    return 0;
}