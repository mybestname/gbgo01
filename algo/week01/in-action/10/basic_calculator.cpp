#include <iostream>
#include <vector>
#include <stack>
#include <string>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int calculate(string s) {
        stack<char> ops;
        vector<string> tokens; // 做为结果的后缀表达式
        int64_t val = 0;
        bool num_started = false;
        bool need_zero = true; // 通过补零来处理正负号。
        for (char ch : s) {
            // parse number
            if (ch >= '0' && ch <= '9') {
                val = val*10 + ch-'0';  // char -> int
                num_started = true;
                continue;
            }else if (num_started) { //
                tokens.push_back(to_string(val));
                num_started = false;
                need_zero = false; //数的后面不需要
                val = 0;
            }
            // skip space
            if (ch == ' ') continue;
            // parse op
            if (ch == '(') {
                ops.push(ch);
                need_zero = true; // (后面需要补零。
                continue;
            }
            if (ch == ')') {
               // 如果是右括号，则不断弹栈，（注意empty判断）
               while(!ops.empty() && ops.top() != '(') {
                   tokens.push_back(string(1,ops.top()));
                   ops.pop();
               }
               ops.pop(); // pop
               need_zero = false; // )后面不需要0
               continue;
            }
            // 处理+-*/
            if (need_zero) tokens.push_back("0"); // +-*/之前，按需要补零
            while (!ops.empty() && getRank(ops.top()) >= getRank(ch)) { //栈顶级别大于当前级别。
                tokens.push_back(string(1, ops.top()));
                ops.pop();
            }
            // 最新字符入栈
            ops.push(ch);
            need_zero = true; //其它负号需要补零
        }
        if (num_started) tokens.push_back(to_string(val)); // 全部是数字？
        while(!ops.empty()) { //还是不空
           tokens.push_back(string(1,ops.top()));
           ops.pop();
        }
        return evalRPN(tokens);
    }

    static int getRank(char ch) {
        if (ch == '+' || ch == '-' ) return 1;
        if (ch == '*' || ch == '/' ) return 2;
        return 0;
    }

    static int evalRPN(vector<string> &tokens) {
        stack<int64_t> s;
        for (string& token : tokens){
            if (token[0] == '+' || token[0] == '-' || token[0] == '*' || token[0] == '/') {
                int64_t b = s.top();  //注意栈顶是操作数b ，因为LIFO
                s.pop();
                int64_t a = s.top();
                s.pop();
                s.push(calc(a,b,token));
            }else { // is number
                s.push(stoi(token));
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
    string input = "+48 + -48";
    int result = Solution::calculate(input);
    cout <<  input  << "=" << result << endl; // 22
    return 0;
}