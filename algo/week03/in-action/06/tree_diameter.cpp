#include <iostream>
#include <vector>
#include <queue>
#include "../../../base/algo_base.h"

using namespace std;

// 思路1：2次DFS解法
// https://oi-wiki.org/graph/tree-diameter/
class Solution1 {
public:

    int treeDiameter(vector<vector<int>> edges) {
        size_t N = edges.size();
        // 加载Edge数据
        Edge = vector<vector<int>>(N+1, vector<int>());
        for (auto& e : edges) {
            int x = e[0];
            int y = e[1];
            addEdge(x,y);
            addEdge(y,x);
        }
        P = 0;
        Dist = vector<int>(N+1);
        //cout << "1st dfs" << endl;
        dfs(2,3);  //首先对任意一个结点做 DFS 求出最远的结点p
        //Dist = vector<int>(N+1);
        Dist[P] = 0;     //清空Dist[P], 为什么其它的点的距离不用清空？?

        //cout << "2nd dfs" << endl;
        dfs(P,3);     //然后以这个结点为根结点再做 DFS 到达另一个最远结点。
        return Dist[P];  //得到直径
    }

private:
    vector<int> Dist;
    int P = 0 ;   //Dist[P]为直径
    void dfs(int x, int y) {
        //cout << "dfs (" << x << "," << y <<")" << endl;
        //cout << "  P=" << P << ",Dist=" << Dist << ",Dist[P]=" << Dist[P] <<  endl;
        for (int i : Edge[x]) {  // x点的所有边
            if ( i == y ) continue; //剔除掉直接指向y的路径。x-y最小路径
            Dist[i] = Dist[x] +1 ;  //x到i点的距离，应该为x点本身距离加1
            if (Dist[i] > Dist[P]) P = i; //如果i点距离更大，更新为最大距离
            dfs(i, x); //递归下一层的距离。
        }
    }
    void addEdge(int x, int y) {
        Edge[x].push_back(y);
    }
    vector<vector<int>> Edge;
};

// 思路1：版本2 2次DFS变体
// https://zhuanlan.zhihu.com/p/115966044
// 区别主要在于dfs函数的设计。
class Solution1_2 {
public:
    int treeDiameter(vector<vector<int>> edges) {
        // 加载Edge数据
        size_t N = edges.size();
        Edge = vector<vector<int>>(N+1, vector<int>());
        for (auto& e : edges) {
            int x = e[0];
            int y = e[1];
            addEdge(x,y);
            addEdge(y,x);
        }
        Dist = vector<int>(N+1);
        Visit = vector<bool>(N+1,false);
        int P = 0;
        int maxDist = 0;
        //cout << "s1_2: 1st dfs" << endl;
        dfs(0); // start from 0;
        for (int i=0; i <= N; i++) {
            if (Dist[i] > maxDist) P=i;
        }
        Dist = vector<int>(N+1);
        Visit = vector<bool>(N+1,false);
        //cout << "s1_2: 2nd dfs" << endl;
        dfs(P);       // start from P;
        maxDist = 0;
        for (int i=0; i <= N; i++) {
            maxDist = max(maxDist,Dist[i]);
        }
        return maxDist;
    }

private:
    vector<bool> Visit;
    vector<int> Dist;
    void addEdge(int x, int y) {
        Edge[x].push_back(y);
    }
    vector<vector<int>> Edge;
    void dfs(int from){
        //cout << "s1_2: dfs (" << from << ")" << endl;
        //cout << "     Dist=" << Dist << ",Visit=" << Visit <<  endl;
        for (int i=0 ; i < Edge[from].size(); i++) {
           int to =  Edge[from][i];
           if (!Visit[to]){
               Visit[to] = true;
               Dist[to] = Dist[from] + 1;
               dfs(to);
           }
        }
    }
};

// 思路2：2次BFS解法,
// https://www.acwing.com/solution/LeetCode/content/5800/
class Solution2 {
public:
    int treeDiameter(vector<vector<int>> edges) {
        // 加载Edge数据
        size_t N = edges.size();
        Edge = vector<vector<int>>(N+1, vector<int>());
        for (auto& e : edges) {
            int x = e[0];
            int y = e[1];
            addEdge(x,y);
            addEdge(y,x);
        }
        pair<int,int> p;
        p = bfs(0);
        p = bfs(p.first);
        return p.second;
    }
private:
    pair<int,int> bfs(int start){
        queue<int> Q;    //bfs需要使用队列
        vector<int> Dist(Edge.size(), -1);
        Q.push(start);
        Dist[start] = 0;
        pair<int,int> ret;
        while (!Q.empty()){
            int x = Q.front(); //
            Q.pop();
            for (auto y : Edge[x]) {
                if (Dist[y] == -1) {
                    Dist[y] = Dist[x] + 1;
                    Q.push(y);
                }
            }
        }
        ret.first = 0, ret.second = 0;
        for (int i = 0 ; i < Dist.size(); i++) {
            if (Dist[i] > ret.second) {
                ret.first = i;
                ret.second = Dist[i];
            }
        }
        return ret;
    }
    void addEdge(int x, int y) {
       Edge[x].push_back(y);
    }
    vector<vector<int>> Edge;
};

// 思路2的优化，因为已经是bfs，所以最后部分求最大dist属于冗余。
class Solution2_2 {
public:
    int treeDiameter(vector<vector<int>> edges) {
        // 加载Edge数据
        size_t N = edges.size();
        Edge = vector<vector<int>>(N+1, vector<int>());
        for (auto& e : edges) {
            int x = e[0];
            int y = e[1];
            addEdge(x,y);
            addEdge(y,x);
        }
        pair<int,int> p;
        // cout << "S2_2 : 1st bfs from 0" << endl;
        p = bfs(0);
        // cout << "S2_2 : p=" << p << endl;
        // cout << "S2_2 : 2nd bfs from " << p.first << endl;
        p = bfs(p.first);
        return p.second;
    }
private:
    pair<int,int> bfs(int start){
        queue<int> Q;    //bfs需要使用队列
        vector<int> Dist(Edge.size(), -1);
        Q.push(start);
        Dist[start] = 0;
        pair<int,int> ret;
        while (!Q.empty()){
            int x = Q.front(); //
            Q.pop();
            ret.first = x;
            ret.second = Dist[x];   // bfs本身，导致最终的ret就是max
            // cout << "Dist="<< Dist << ", ret is " << ret << endl;
            for (auto y : Edge[x]) {
                if (Dist[y] == -1) {
                    Dist[y] = Dist[x] + 1;
                    Q.push(y);
                }
            }
        }
        return ret;
    }
    void addEdge(int x, int y) {
        Edge[x].push_back(y);
    }
    vector<vector<int>> Edge;
};

// 解法3：基于树的动态规划，未完成
// 记录当1为树的根时，每个节点作为子树的根向下，所能延伸的最远距离d1，和次远距离d2，那么直径就是所有d1+d2的最大值。
// https://oi-wiki.org/graph/tree-diameter/
// https://blog.csdn.net/pfdvnah/article/details/102922383
//
class Solution3 {
public:
    int treeDiameter(vector<vector<int>> edges) {
        // 加载Edge数据
        size_t N = edges.size();
        Edge = vector<vector<int>>(N + 1, vector<int>());
        for (auto &e : edges) {
            int x = e[0];
            int y = e[1];
            addEdge(x, y);
            addEdge(y, x);

        }
        Dist1 = vector<int>(N+1);
        Dist2 = vector<int>(N+1);
        D=0;
        dfs(0, -1);  //注意，初始y点不能在图中，否则有可能是错误答案。
        return D;
    }
private:
    void dfs(int x, int y) {
        //cout << "S3: dfs ("<<x <<","<<y<<")" << endl;
        Dist1[x]=0, Dist2[y]=0;
        for (int i : Edge[x] ) {
            if ( i == y ) continue;
            dfs(i,x);  // 一直到最深
            int cur_d = Dist1[i] + 1;
            if ( cur_d > Dist1 [x]) {
                Dist2[x] = Dist1[x];
                Dist1[x] = cur_d;
            }
            else if ( cur_d > Dist2[x]) {
                Dist2[x] = cur_d;
            }
        }
        D = max(D, Dist1[x]+Dist2[x]);
    }
    int D = 0 ;
    // 最远距离D1， 次远距离D2
    vector<int> Dist1, Dist2;
    void addEdge(int x, int y) {
        Edge[x].push_back(y);
    }
    vector<vector<int>> Edge;

};
int main(){
    vector<vector<int>> edges = {{0,1},{1,2},{2,3},{1,4},{4,5}};
    {
        Solution1 s;
        auto result = s.treeDiameter(edges);
        cout << "S1: edges=" << edges << ",diameter=" << result << endl;
    }
    {
        Solution1_2 s;
        auto result = s.treeDiameter(edges);
        cout << "S1_2: edges=" << edges << ",diameter=" << result << endl;
    }
    {
        Solution2 s;
        auto result = s.treeDiameter(edges);
        cout << "S2: edges=" << edges << ",diameter=" << result << endl;
    }
    {
        Solution2_2 s;
        auto result = s.treeDiameter(edges);
        cout << "S2_2: edges=" << edges << ",diameter=" << result << endl;
    }
    {
        Solution3 s;
        auto result = s.treeDiameter(edges);
        cout << "S3: edges=" << edges << ",diameter=" << result << endl;
    }
    return 0;
}