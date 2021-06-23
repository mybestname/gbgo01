#include <iostream>
#include <vector>

using namespace std;
class Solution {
private:
    static void printDp(const int i, const int j, const vector<vector<int>> &dp, const vector<vector<int>> &grid ){
        char buffer[50];
        snprintf(buffer, 36, "i=%d,j=%d, dp[%d][%d]=%d, grid[%d][%d]=%d\n", i, j, i, j, dp[i][j], i, j, grid[i][j]);
        cout << buffer << endl;
    }
public :
    static int minPathSum(const vector<vector<int>> &grid) {
        int m = grid.size();
        int n = grid[0].size();
        // m x n matrix
        vector<vector<int>> dp(m, vector<int>(n, 0));
        for (int i = 0 ; i < m ; i++) {
            for (int j = 0; j < n; j++) {
                if (i == 0  && j == 0 ) {
                    dp [i][j] = grid[i][j];
                } else if (i == 0) {
                    dp[i][j] = dp[i][j-1] + grid[i][j];
                } else if (j == 0) {
                    dp[i][j] = dp[i-1][j] + grid[i][j];
                } else {
                    dp[i][j] = min(dp[i - 1][j], dp[i][j - 1]) + grid[i][j];
                }
                printDp(i,j,dp,grid);
            }
        }
        return dp[m-1][n-1];
    }
};

int main() {
    Solution sol;
    vector<vector<int>> grid = {{1,3,1},{1,5,1},{4,2,1}};
    int min = Solution::minPathSum(grid);
    cout << min << endl; // 7

    //input 2
    grid = {{1,2,3},{4,5,6}};
    min = Solution::minPathSum(grid);
    cout << min << endl; // 12
    return 0;
}
