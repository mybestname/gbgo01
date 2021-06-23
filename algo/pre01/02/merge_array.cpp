#include <iterator> // needed for std::ostram_iterator
#include <iostream>
#include <vector>
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
    void merge(vector<int>& num1, int m, vector<int>& num2, int n) {
        // m+n游标
        int index = m + n;
        // num1 (m) 和 num2 (n) 自减，num1 做完时候，num2已经做完的情况
        while ( m > 0 && n > 0 ) {
            // 如果 num1 大， index记录，使用 num1
            if (num1[m-1] > num2[n-1]) {
                num1[index] = num1[m-1];
                --m; //下一个num1
            } else {
            // 如果 num2 大， index记录，使用 num2
                num1[index] = num2[n-1];
                --n; //下一个num2
            }
            --index; //m+n游标减1
        }
        // num1 做完，但是 num2还是没有做完，直接回填即可
        while (n > 0) {
            num1[index--] = num1[n--];
        }
    }
};

int main() {
    Solution sol;

    vector<int> nums1 = {1,2,3,0,0,0};
    int m = 3;
    vector<int> nums2 = {2,5,6};
    int n = 3;
    sol.merge(nums1,m,nums2, n);
    cout << nums1 << endl;

    nums1 = {1}; m = 1;
    nums2 = {}; n = 0;
    sol.merge(nums1,m,nums2, n);
    cout << nums1 << endl;

    return 0;
}
