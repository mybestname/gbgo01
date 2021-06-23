#include <iostream>
#include <vector>

using namespace std;
class Solution {
public :
    static int maxProfit(const vector<int> &price) {
        int result = 0;
        for (int i=0 ; i< price.size() - 1; i++) {
            if (price[i+1] > price[i]) {
               result += price[i+1] - price[i];
            }
        }
        return result;
    }
};

int main() {
    Solution sol;
    // [7,1,5,3,6,4]
    //
    //   7  1  5  3  6  4
    //      -6 4  2  -3 2
    vector<int> price = {7,1,5,3,6,4};
    int profit = Solution::maxProfit(price);
    cout << profit << endl;  // 7

    // [1,2,3,4,5]
    //   1  2  3  4  5
    //      1  1  1  1
    price = {1,2,3,4,5};
    profit = Solution::maxProfit(price);
    cout << profit << endl;  // 4

    // [7,6,4,3,1]
    //     7  6  4  3  1
    //       -1 -2 -1 -2
    price = {7,6,4,3,1};
    profit = Solution::maxProfit(price);
    cout << profit << endl;  // 0
    return 0;
}