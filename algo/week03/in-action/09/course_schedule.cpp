#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"

using namespace std;
// 思路1：DFS解法
class Solution1 {
public:
    bool canFinish(int numCourses, vector<vector<int>>& prerequisites) {
        if (prerequisites.empty()) return true;
        n = numCourses;
        edge = vector<vector<int>> (n, vector<int>());
        for (auto& p : prerequisites) {
            auto x = p[1];
            auto y = p[0];
            add_edge(x,y);
        }
        visit = vector<bool>(n,false);
        processed = vector<bool>(n, false);
        for (int i = 0; i < numCourses; i++) {
            hasCycle = false;
            dfs(i, -1);
            if (hasCycle) return false;
        }
        return true;
    }

private:
    // 图的深度优先遍历需要visit数组，避免重复访问。
    // 普通的fa法判断，适用于无向图
    // 如果是有向图，判断有环，如果用DFS，需要3色标记（其实也就是内存垃圾回收的做法）
    // 这里在无向图算法中添加了一个是否处理中的状态位。
    // 进一步优化可以让
    // vector<bool> visit -> vector<int> visit，
    // 里面存一个三色状态（0，1，2）：未搜索，搜索中，已完成，
    // 搜索中（visit==1）则表示有环。（已经搜索过这个节点，但还没有回溯到该节点，即该节点还没有入栈，还有相邻的节点没有搜索完成）
    void dfs(int x, int fa) {
        // 首先标记访问
        visit[x] = true;
        // 第二：遍历所有的出边
        for (auto y : edge[x]) {
            if ( x == fa ) continue;
            if (!visit[y]) {
                if (!hasCycle) dfs(y, x); // x 是 y 的 father
            }
            else if (visit[y] && !processed[y]) {
                hasCycle=true ;
                return;
            }
        }
        processed[x] = true;
    }
    int n = 0;
    void add_edge(int x, int y) {
        edge[x].push_back(y);
    }
    vector<vector<int>> edge;
    vector<bool> visit;
    vector<bool> processed;
    bool hasCycle = false;
};

// 思路2 ：BFS+topsort
class Solution2 {
public:
    bool canFinish(int numCourses, vector<vector<int>> &prerequisites) {
        n = numCourses;
        in_degree = vector<int>(n, 0);
        edge = vector<vector<int>>(n, vector<int>());
        for (auto& e : prerequisites) {
            auto x = e[0];
            auto y = e[1];
            add_edge(y, x);  // 按题意 e[1] -> e[0]
        }
        return top_sort() == n;
    }

private:
    int top_sort() {
        int learned = 0;
        queue<int> q;
        //1. 从所有的零入度点出发
        for (int i=0; i<n; i++) {
            if(in_degree[i] == 0) {
               q.push(i);
            }
        }
        while (!q.empty()) {
            int x = q.front(); //取队头
            q.pop();
            learned++; // x已学
            // 考虑x的所有出边
            for (auto y : edge[x]) {
                // 去掉 x 对 y 的约束关系。
                in_degree[y] -- ;  // 即入度数减1；
                if (in_degree[y] == 0) { // 说明y可以学
                    q.push(y);
                }
            }
        }
        // 最后看是不是每个点都被学过。
        // 返回学过课的总数
        return learned;
    }
    int n = 0; // n points
    vector<vector<int>> edge;
    void add_edge(int x, int y) {
        edge[x].push_back(y);
        in_degree[y]++;  // x->y, y++入度
    }
    vector<int> in_degree;
};


int main() {
    struct Test {
        int numCourses;
        vector<vector<int>> prerequisites;
        bool expect;
    };
    vector<Test> tests = {
           { .numCourses = 1, .prerequisites={}, .expect=true},
           { .numCourses = 2, .prerequisites={}, .expect=true},
           { .numCourses = 2, .prerequisites={{1,0}}, .expect = true},
           { .numCourses = 2, .prerequisites={{1,0},{0,1}}, .expect = false},
           { .numCourses = 3, .prerequisites={{1,0}}, .expect = true},
           { .numCourses = 3, .prerequisites={{1,0},{2,0},{0,2}}, .expect = false},
            //              1  ->   0
            //                   v/ /^
            //                    2
            { .numCourses = 4, .prerequisites= {{2,0},{1,0},{3,1},{3,2},{1,3}}, .expect = false},
            //
            { .numCourses = 5, .prerequisites={{1,4},{2,4},{3,1},{3,2}}, .expect=true},
            //     1->4 2->4 3->1 3->2
            //
            //                 3->  1 ->  4
            //                  \       /
            //                   \--2--/
            //
            { .numCourses = 20,.prerequisites={{0,10},{3,18},{5,5},{6,11},{11,14},{13,1},{15,1},{17,4}}, .expect= false},
            //
            //
            //                         0->10 3->18 5->5 6->11->14 13->1 15->1 17->4
    };
    {
        Solution1 s;
        for (auto &test : tests) {
            auto result = s.canFinish(test.numCourses, test.prerequisites);
            cout << "S1: numCourses=" << test.numCourses << ",prerequisites=" << test.prerequisites
                 << ", can_finish expect="<< test.expect << ",got=" << result << endl;
        }
    }

    {
        Solution2 s;
        for (auto &test : tests) {
            auto result = s.canFinish(test.numCourses, test.prerequisites);
            cout << "S2: numCourses=" << test.numCourses << ",prerequisites=" << test.prerequisites
                 << ", can_finish expect="<< test.expect << ",got=" << result << endl;
        }
    }
    return 0;
}

