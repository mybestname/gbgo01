#include <iostream>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<int> findRedundantConnection(vector<vector<int>>& edges) {
        int n = edges.size();
        edge = vector<vector<int>> (n+1, vector<int>());
        visit = vector<bool>(n+1, false);
        hasCycle = false;
        for (auto& e : edges) {
            auto x = e[0];
            auto y = e[1];
            add_edge(x,y);
            add_edge(y,x);
            // 每新加一条边，dfs一下看没有环
            for(int i = 0 ; i <= n; i++) visit[i] = false;  // 每次调用前清除visit状态
            dfs(1, -1);    //-1表示，1没有父亲。
            if (hasCycle) return e;  // 如果有环，那么这条边就是寻找边。
        }
        return {};
    }

private:
    void dfs(int x, int father) {
        // 首先，标记已经访问
        visit[x] = true;
        // 第二，遍历所有出边
        for(auto& y : edge[x]) {
            if (y == father) continue; // 如果返回了父亲，这个边不用再看，这个不是环。
            // 不是父亲，但是这个点又访问过，说明是环
            if (visit[y]) hasCycle = true;
            else dfs(y, x);  //继续dfs，y的父亲是x，或y源于x
        }
    }
    bool hasCycle = false;
    void add_edge(int x, int y) {
       edge[x].push_back(y);
    }
    vector<vector<int>> edge;
    vector<bool> visit;


};

// 小改进：调用前visit状态可以在递归调用中恢复。
class Solution1_2 {
public:
    vector<int> findRedundantConnection(vector<vector<int>>& edges) {
        int n = edges.size();
        edge = vector<vector<int>> (n+1, vector<int>());
        visit = vector<bool>(n+1, false);
        hasCycle = false;
        for (auto& e : edges) {
            auto x = e[0];
            auto y = e[1];
            add_edge(x,y);
            add_edge(y,x);
            dfs(1, -1);
            if (hasCycle) return e;
        }
        return {};
    }
private:
    void dfs(int x, int father) {
        visit[x] = true;
        for(auto& y : edge[x]) {
            if (y == father) continue;
            if (visit[y]) hasCycle = true;
            else dfs(y, x);
        }
        visit[x] = false; //恢复visit状态
    }
    bool hasCycle = false;
    void add_edge(int x, int y) {
        edge[x].push_back(y);
    }
    vector<vector<int>> edge;
    vector<bool> visit;
};


int main() {
    // 输入: [[1,2], [2,3], [3,4], [1,4], [1,5]]
    // 输出: [1,4]
    vector<vector<int>> edges = {{1,2}, {2,3}, {3,4}, {1,4}, {1,5}};
    {
        Solution s;
        auto redundant = s.findRedundantConnection(edges);
        cout << "edges=" << edges << ",redundant_edge=" << redundant << endl;
    }
    {
        Solution1_2 s;
        auto redundant = s.findRedundantConnection(edges);
        cout << "S1_2 : edges=" << edges << ",redundant_edge=" << redundant << endl;
    }
    return 0;
}
