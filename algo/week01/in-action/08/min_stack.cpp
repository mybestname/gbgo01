#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class MinStack {
public:
    MinStack() = default;
    void push(int val) {
        preMin.push(stack.empty()? val : min (val, preMin.top()));
        stack.push(val);
    }
    void pop() {
        preMin.pop();
        stack.pop();
    }
    int top() {
       return stack.top();
    }
    int getMin() {
        return preMin.top();
    }
private:
    stack<int> preMin;
    stack<int> stack;
};

int main() {
    MinStack minStack = *new MinStack();
    minStack.push(-2);
    minStack.push(0);
    minStack.push(-3);
    int p1 = minStack.getMin();  // --> 返回 -3.
    minStack.pop();
    int p2 = minStack.top();     // --> 返回 0.
    int p3 = minStack.getMin();  // --> 返回 -2.
    vector<int> l = {p1,p2,p3};
    cout << l << endl;  // [-3,0,-2]
    return 0;
}