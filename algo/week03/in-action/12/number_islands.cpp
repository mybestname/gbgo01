#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"

using namespace std;
// 使用DFS
class Solution1 {
public:
    int numIslands(vector<vector<char>>& grid) {
        int ans = 0;
        // M x N grid
        this->M = grid.size();
        this->N = grid[0].size();
        this->visit = vector<vector<bool>>(M, vector<bool>(N,false));
        //
        for (int i =0; i< M; i++) {
            for (int j =0; j<N; j++) {
                if (grid[i][j] == '1' && !visit[i][j]) {
                    dfs(grid, i, j);
                    ans++ ;
                    // 递归要么从主函数进入，要么从递归进入。
                    // 有多少次从主函数进入递归，就有多少块儿。
                }
            }
        }
        return ans;
    }

private:
    void dfs(vector<vector<char>>& grid, int x, int y) {
        // 终止条件 (不需要：都visit了结束）
        // 标记已经访问
        visit[x][y] = true;
        // 考虑所有的出边 (上下左右四个方向，四条出边）
        for (int i=0 ; i< 4; i++) {
           int nextX = x + dx[i];
           int nextY = y + dy[i];
           // nextX,nextY的合法性检查 （访问数组前，一定记得检查合法性）
            if (nextX < 0 || nextY < 0 || nextX >= M || nextY >= N) continue;
           // 什么时候可以走？1 且 没有走过
           if (grid[nextX][nextY] == '1' && !visit[nextX][nextY]) {
              dfs(grid, nextX,nextY);
           }
        }
    }
    int M{};
    int N{};
    vector<vector<bool>> visit;
    // 方向数组
    // dir=0 N; dir=1 E; dir=2 S; dir=3 W
    const int dx[4] = { -1,  0,  0, 1 };
    const int dy[4] = {  0, -1,  1, 0 };
};


// 使用bfs
class Solution2 {
public:
    int numIslands(vector<vector<char>>& grid) {
        int ans = 0;
        // M x N grid
        this->M = grid.size();
        this->N = grid[0].size();
        this->visit = vector<vector<bool>>(M, vector<bool>(N,false));
        for (int i =0; i< M; i++) {
            for (int j =0; j<N; j++) {
                if (grid[i][j] == '1' && !visit[i][j]) {
                    bfs(grid, i, j);
                    ans++ ;
                }
            }
        }
        return ans;
    }
private:
    void bfs(vector<vector<char>>& grid, int x, int y) {
        // 广搜需要队列, 二维坐标
        queue<pair<int,int>> q;
        // 首先，push起点
        q.push(make_pair(x,y));
        visit[x][y] = true;
        while(!q.empty()) {
            int _x = q.front().first;
            int _y = q.front().second;
            q.pop();
            // 扩展所有的出边, 上下左右四个方向
            for (int i=0 ; i< 4; i++) {
                int nextX = _x + dx[i];
                int nextY = _y + dy[i];
                if (nextX < 0 || nextY < 0 || nextX >= M || nextY >= N) continue;
                // 什么可以被push入队列？
                if (grid[nextX][nextY] == '1' && !visit[nextX][nextY]) {
                   q.push(make_pair(nextX,nextY));
                   visit[nextX][nextY] = true; //一定记得入队时候要标记visit
                }
            }
        }
    }
    size_t M{};
    size_t N{};
    vector<vector<bool>> visit;
    const int dx[4] = { -1,  0,  0, 1 };
    const int dy[4] = {  0, -1,  1, 0 };
};



int main(){
    struct Test {
        vector<vector<char>> grid;
        int expect;
    };
    vector<Test> tests = {
            {.grid =
                    {
                            {'1','1','1','1','0'},
                            {'1','1','0','1','0'},
                            {'1','1','0','0','0'},
                            {'0','0','0','0','0'}
                    }, .expect=1 },

            {.grid =
                    {
                            {'1','1','0','0','0'},
                            {'1','1','0','0','0'},
                            {'0','0','1','0','0'},
                            {'0','0','0','1','1'}
                    }, .expect=3 },
    };
    {
        Solution1 s;
        for (auto &test : tests) {
            auto result = s.numIslands(test.grid);
            cout << " grid=" << test.grid<< ",expect=" << test.expect << ",got=" << result << endl;
        }
    }

    {
        Solution2 s;
        for (auto &test : tests) {
            auto result = s.numIslands(test.grid);
            cout << " grid=" << test.grid<< ",expect=" << test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}

