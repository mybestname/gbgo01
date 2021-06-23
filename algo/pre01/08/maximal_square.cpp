#include <iostream>
#include <vector>

using namespace std;
class Solution {
private:
    static void printDp(const int i, const int j, const vector<vector<int>> &dp, const vector<vector<char>> &matrix ){
        char buffer[50];
        snprintf(buffer, 36, "i=%d,j=%d, dp[%d][%d]=%d, matrix[%d][%d]=%c\n", i, j, i, j, dp[i][j], i-1, j-1, matrix[i-1][j-1]);
        cout << buffer << endl;
    }
public :
    static int maximalSquare(const vector<vector<char>> &matrix) {
        if ( matrix.empty() || matrix[0].empty()) return 0; //check empty
        uint m = matrix.size();
        uint n = matrix[0].size();
        int max_width = 0;
        // m x n matrix
        // 注意dp始终表示右下脚，所以使用一个(m+1)x(n+1)矩阵，同时index从1开始。
        vector<vector<int>> dp(m+1, vector<int>(n+1, 0));
        for (int i = 1 ; i <= m ; i++) {
            for (int j = 1; j <= n; j++) {
                if (matrix[i-1][j-1] == '1') { //而matrix从(0,0)开始。
                    dp[i][j] = min(dp[i-1][j-1],min(dp[i-1][j],dp[i][j-1])) + 1;
                }
                max_width = max(max_width, dp[i][j]);
                printDp(i,j,dp,matrix);
            }
        }
        return max_width*max_width;
    }
};

int main() {
    Solution sol;
    vector<vector<char>> matrix = {{'1','0','1','0','0'},{'1','0','1','1','1'},{'1','1','1','1','1'},{'1','0','0','1','0'}};
    int max = Solution::maximalSquare(matrix);
    cout << max << endl; //4

    matrix = {{'0','1'},{'1','0'}};
    max = Solution::maximalSquare(matrix);
    cout << max << endl; // 1

    matrix = {{'0'}};
    max = Solution::maximalSquare(matrix);
    cout << max << endl; // 0
    return 0;
}

