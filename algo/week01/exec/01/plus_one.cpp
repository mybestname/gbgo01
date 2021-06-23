#include <iostream>
#include <vector>
#include <iterator> // needed for std::ostream iterator
template <typename T>
std::ostream& operator<< (std::ostream& out, const std::vector<T>& v) {
    if ( !v.empty() ) {
        out << '[';
        std::copy (v.begin(), v.end(), std::ostream_iterator<T>(out, ","));
        out << "\b]";
    }
    return out;
}

using namespace std;
class Solution {
public :
    static vector<int> plusOne(vector<int> &digits) {
        for (int i = digits.size() - 1; i >= 0 ; i--) {
            digits[i]++;
            if (digits[i] != 10) {
                return digits;     // return directly if need not carry
            } else {
             digits[i] = 0;        // set zero if carry
            }
        }
        // still not return, all carry in this way.
        digits[0] = 1;             // set 1 at the left-most
        digits.push_back(0);       // append zero at the right-back
        return digits;
    }
};

int main() {
    Solution sol;
    vector<int> digits = {1,2,3};
    vector<int> result = Solution::plusOne(digits);
    cout << result << endl; // [1,2,4]
    digits = {4,3,2,2};
    result = Solution::plusOne(digits);
    cout << result << endl; // [4,3,2,3]
    digits = {0};
    result = Solution::plusOne(digits);
    cout << result << endl; // [1]
    digits = {9,9,9};
    result = Solution::plusOne(digits);
    cout << result << endl; // [1]
}

