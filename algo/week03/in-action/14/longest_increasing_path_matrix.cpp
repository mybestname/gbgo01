#include <iostream>
#include <vector>
#include <queue>
#include <unordered_map>
#include "../../../base/algo_base.h"
using namespace std;

// DFS
class Solution1 {
public:
    int longestIncreasingPath(vector<vector<int>>& matrix) {
        M = matrix.size();
        N = matrix[0].size();
        answers = vector<vector<int>>(M,vector<int>(N,-1));
        int longest = 0;
        for (int i =0 ; i < M; i++) {
            for (int j=0; j < N; j++) {
               longest = max(longest,how_far(matrix, i, j));
            }
        }
        return longest;
    }

private:
    // 子问题：从x,y出发能走多远
    int how_far(vector<vector<int>>& matrix, int x, int y) {
        // 避免重算
        if (answers[x][y] != -1) return answers[x][y];
        answers[x][y] = 1; // 至少是1
        // 可以往哪里走？四个方向
        for (int i = 0; i< 4; i++) {
            int nx = x+dx[i];
            int ny = y+dy[i];
            // 首先，合法性判断
            if (nx < 0 || ny < 0 || nx >= M || ny >= N) continue;
            // 下一个要长
            if (matrix[nx][ny] > matrix[x][y]) {
                // x,y答案 为 下一步+1 和 本身 之间的大者
                answers[x][y] = max(answers[x][y], how_far(matrix,nx,ny)+1);
            }
        }
        return answers[x][y]; //注意可能的bug：30行忘记设定初始值，导致返回-1，答案永远不应该返回-1
    }
    // 矩阵长度
    size_t M{},N{};
    // 方向数组
    const int dx[4] = { -1,  0, 0, 1 };
    const int dy[4] = {  0, -1, 1, 0 };
    // 答案缓存（避免重算）
    vector<vector<int>> answers;
};

// BFS
class Solution2 {
public:
    int longestIncreasingPath(vector<vector<int>> &matrix) {
        M = matrix.size();
        N = matrix[0].size();
        dist = vector<int>(M*N);
        // M*N个点的出边
        edges = vector<vector<int>>(M*N, vector<int>());
        in_degree = vector<int>(M*N);

        // 对每一个点, 构造出边数组
        for (int i =0 ; i < M; i++) {
            for (int j = 0; j < N; j++) {
                // 四种方向
                for (int k = 0; k < 4; k++) {
                    int ni = i + dx[k];
                    int nj = j + dy[k];
                    if (valid_next(ni,nj) && matrix[ni][nj] > matrix[i][j]) {
                        // (i,j) -> (ni,nj)
                        add_edge(num(i,j), num(ni,nj));
                    }
                }
            }
        }
        // 拓扑排序
        topsort();
        int ans = 0;
        for (int i = 0; i< M*N; i++){
            ans = max(ans, dist[i]);
        }
        return ans;
    }
private:
    void topsort(){
        // 拓扑排序需要使用队列
        queue<int> q;
        // 对所有入度为0的点
        for (int i=0; i< M*N; i++) {
            if(in_degree[i] == 0) {
                q.push(i);
                dist[i] = 1;
            }
        }
        while (!q.empty()) {
            //取队头
            int x = q.front();
            q.pop();
            // cout << "S2 q.front=" << x << ",edges=" <<edges[x] << endl;
            // 考虑所有的出边
            for (auto y : edges[x]) {
                in_degree[y]--;
                dist[y] = max(dist[y], dist[x]+1);
                // 当入度为0时候，可以push
                if (in_degree[y] == 0) {
                    q.push(y);
                }
            }
        }
    }
    // 到M*N个点的长度
    vector<int> dist;
    // (col,row) to number in [0,M*N)
    // 0 <= x <= M-1, 0<= y <= N-1
    int num(int x, int y) {
        return x*N + y;  //possible bug, x * N is ulong
    }

    // 下一个点的合法性检查
    bool valid_next(int nx, int ny) const {
        return !(nx < 0 || ny < 0 || nx >= M || ny >= N);
    }
    size_t M{},N{};
    const int dx[4] = { -1,  0, 0, 1 };
    const int dy[4] = {  0, -1, 1, 0 };
    // 出边数组
    vector<vector<int>> edges;
    // 入度数
    vector<int> in_degree;
    void add_edge(int x, int y) {
        edges[x].push_back(y);
        in_degree[y]++;  // x -> y
    }
};

// BFS S2 ：
// 1. 不用出边数组，只用点的入度数就够了
//    - 如果不用出边数组，那么当top sort时候
//      - 每一轮搜索都要遍历当前层的所有单元格，更新其余单元格的出度，并将出度变为 0 的单元格加入下一层搜索。
// 2. 使用pair表示点坐标，这样不用在考虑多一层转化，代码简洁。
// 3. 没有必要存所有的dist，在每次pop时候，累加即可。
class Solution2_2 {
public:
    int longestIncreasingPath(vector<vector<int>> &matrix) {
        M = matrix.size();
        N = matrix[0].size();
        in_degree = vector<vector<int>>(M, vector<int>(N,0));
        // 对每一个点, 计算入度数
        for (int i =0 ; i < M; i++) {
            for (int j = 0; j < N; j++) {
                // 四种方向
                for (int k = 0; k < 4; k++) {
                    int ni = i + dx[k];
                    int nj = j + dy[k];
                    if (valid(ni, nj) && matrix[ni][nj] > matrix[i][j]) {
                        in_degree[ni][nj]++;
                    }
                }
            }
        }
        ans = 0;
        // 拓扑排序
        topsort(matrix);
        return ans;
    }
private:
    int ans=0;
    void topsort(vector<vector<int>> &matrix){
        // 拓扑排序需要使用队列
        queue<pair<int,int>> q;
        // 对所有入度为0的点
        for (int i =0 ; i < M; i++) {
            for (int j = 0; j < N; j++) {
                if (in_degree[i][j] == 0) {
                    q.push({i,j});
                }
            }
        }
        while (!q.empty()) {
            // cout << "s2_2: q.front=" << q.front() << ",q.size=" << q.size() << endl;
            ans++;
            auto size = q.size();
            //取所有的cell
            for (int i = 0; i < size; ++i) {
                int x = q.front().first;
                int y = q.front().second;
                q.pop();
                // 考虑所有方向
                for (int k = 0; k < 4; k++) {
                    int nx = x + dx[k];
                    int ny = y + dy[k];
                    if (valid(nx, ny) && matrix[nx][ny] > matrix[x][y]) {
                        in_degree[nx][ny]--;
                        // 当入度为0时候，可以push
                        if (in_degree[nx][ny] == 0) {
                            q.push({nx, ny});
                        }
                    }
                }
            }
        }
    }
    // 下一个点的合法性检查
    bool valid(int nx, int ny){
        return !(nx < 0 || ny < 0 || nx >= M || ny >= N);
    }
    size_t M{},N{};
    const int dx[4] = { -1,  0, 0, 1 };
    const int dy[4] = {  0, -1, 1, 0 };
    // 点(x,y)的入度数
    vector<vector<int>> in_degree;
};

int main(){
    struct Test {
        vector<vector<int>> matrix;
        int expect ;
    };
    vector<Test> tests = {
            {
                    .matrix = {{9,9,4},{6,6,8},{2,1,1}},
                    .expect = 4,
            },
            {
                    .matrix = {{3,4,5},{3,2,6},{2,2,1}},
                    .expect = 4,
            },
            {
                    .matrix = {{1}},
                    .expect = 1,
            },
            {
                    .matrix = {{0},{1},{5},{5}},
                    .expect = 3,
            }
    };

    {
        Solution1 s;
        for (auto &test : tests) {
            auto result = s.longestIncreasingPath(test.matrix);
            cout << "S1: matrix=" << test.matrix
                 << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    {
        Solution2 s;
        for (auto &test : tests) {
            auto result = s.longestIncreasingPath(test.matrix);
            cout << "S2: matrix=" << test.matrix
                 << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    {
        Solution2_2 s;
        for (auto &test : tests) {
            auto result = s.longestIncreasingPath(test.matrix);
            cout << "S2_2: matrix=" << test.matrix
                 << ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}