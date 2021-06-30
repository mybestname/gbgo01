#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;

class NumMatrix {
public:
    NumMatrix(vector<vector<int>>& matrix) {
        sum.clear();
        for(int i = 0; i < matrix.size(); i++) {
            sum.push_back({});
            for (int j = 0; j < matrix[i].size(); j++) {
                // 公式1
                int n = get_sum(i-1,j) + get_sum(i,j-1) - get_sum(i-1,j-1) + matrix[i][j];
                sum[i].push_back(n);
            }
        }
        cout << "sum loaded "<<sum<<endl;
    }

    int sumRegion(int row1, int col1, int row2, int col2) {
        // 公式2
        // sum(p,q,i,j) = A[p]A[q]+...+A[i]A[j]
        //             = S[i][j]-S[i][q-1]-S[p-1][j]+S[p-1][q-1]
        return get_sum(row2,col2) - get_sum(row2,col1-1) - get_sum(row1-1,col2) + get_sum(row1-1,col1-1);
    }
private:
    int get_sum(int i, int j) {
        if (i >= 0 && j >=0) return sum[i][j];
        return 0; //小于零的情况。
    }
    vector<vector<int>> sum;
};

int main() {
    vector<vector<int>> data = {
            {3, 0, 1, 4, 2},
            {5, 6, 3, 2, 1},
            {1, 2, 0, 1, 5},
            {4, 1, 0, 1, 7},
            {1, 0, 3, 0, 5}};
    NumMatrix m = NumMatrix(data);
    int r1 = m.sumRegion(2, 1, 4, 3); // 8
    int r2 = m.sumRegion(1, 1, 2, 2); // 11
    int r3 = m.sumRegion(1, 2, 2, 4); // 12
    vector<int> result = {r1, r2 , r3};
    cout << "result = " << result << endl;
    return 0;
}
