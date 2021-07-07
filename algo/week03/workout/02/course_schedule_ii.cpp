#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    vector<int> findOrder(int numCourses, vector<vector<int>>& prerequisites) {
        n = numCourses;
        // 出边数组的初始化, n个空数组
        edges = vector<vector<int>>(n, vector<int>());
        in_degree = vector<int>(n, 0);
        // 加边
        for (auto& pre: prerequisites) {
            int ai = pre[0];
            int bi = pre[1];
            addEdge(bi,ai); // b要先学，才能学a,  b->a

        }
        // 拓扑排序
        auto ans = topsort();
        if (ans.size() < n) return {}; // 无法完成所有
        return ans;
    }

private:
    int n ;
    vector<vector<int>> edges;  //存图
    vector<int> in_degree; // n个点的入度数
    void addEdge(int x, int y) {
       edges[x].push_back(y);
       in_degree[y]++;   //x到y有一条边 x->y
    }
    // 返回学的课程
    vector<int> topsort() {
        vector<int> course;
        // 基于BFS，使用队列
        queue<int> q;
        // 从所有的零入度点出发
        for (int i =0; i<n ; i++) {
            if (in_degree[i]==0) {
                q.push(i);
            }
        }
        // BFS
        while(!q.empty()) {
            int x = q.front(); //队头出队，这门课已经学了。
            q.pop();
            course.push_back(x);
            // 对于x的所有出边
            for (int y : edges[x]) {
                in_degree[y]--; //去掉约束关系
                if (in_degree[y] == 0 ) {
                    q.push(y); // y可学
                }
            }
        }
        return course;
    }
};


int main() {
    int n = 4;
    vector<vector<int>> pereq = {{1,0},{2,0},{3,1},{3,2}};
    Solution s;
    auto result =  s.findOrder(4,pereq);
    // 输入: 4, [[1,0],[2,0],[3,1],[3,2]]
    // 输出: [0,1,2,3] or [0,2,1,3]
    cout << result << endl;
    return 0;
}

