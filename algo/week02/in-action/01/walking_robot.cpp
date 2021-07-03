#include <iostream>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static int robotSim(vector<int> &commands, vector<vector<int>> &obstacles, int choice) {
        if (choice==1) {
           return _robotSim<string>(commands,obstacles, stringConvXY);
        }
        return _robotSim<int64_t>(commands,obstacles, intConvXY);
    }
    template<typename T>
    static int _robotSim(vector<int> &commands, vector<vector<int>> &obstacles, function<T(int, int)> convXY) {
        // setup blockers, use "x,y" as (x,y)
        unordered_set<T> blockers;
        for (auto& o : obstacles) {
            blockers.insert(convXY(o[0],o[1]));
        }
        int x = 0, y= 0; // start from (0,0)
        int ans = 0;
        // 使用方向数组技巧
        //           N, E,  S,  W
        int dx[4] = {0, 1,  0, -1 };
        int dy[4] = {1, 0, -1,  0 };
        //  dir=0 -> N
        //  nextX = x + dx[dir];
        //  nextY = y + dy[dir];
        int dir = 0; // N(0) E(1) S(2) W(3)
        for (int cmd : commands) {
            if (cmd > 0) {
                for (int i = 0 ; i< cmd ; i++) {
                    // try next step
                    int next_x = x + dx[dir];
                    int next_y = y + dy[dir];
                    // find obstacles, stop
                    if (blockers.find(convXY(next_x,next_y))!= blockers.end()){
                        break;
                    }
                    x = next_x;
                    y = next_y;
                }
                ans = max(ans, x*x+y*y);
            } else if (cmd == -1){
             //  - -1 ：向右转 90 度
             // 右转90 = N->E, E->S, S->W, W->N = (dir+1)%4
                dir = (dir + 1) % 4;

            } else if (cmd == -2) {
             //  - -2 ：向左转 90 度
             // N->W, E->N, S->E, W->S (0->3, 1->0, 2->1, 3->2)
             // 左转 -1，避免负数，加mod数，-1+4=3 =>  (dir-1+4)%d = (dir+3)%4
                dir = (dir +3) % 4;
            }

        }
        return ans;
    }

private:
    static string stringConvXY(int x, int y) {
        return to_string(x) + "," + to_string(y);
    }
    static int64_t intConvXY(int x, int y) {
        // return (x+30000) * 60000ll + y + 30000;  // wrong!
        return (x+30000) * 60001ll + y + 30000;     // correct!
        // return x*60001ll + y;   //
        // x,y -> 因为x和y的取值范围为[-3e4,3e4]，
        // 为了保证为正数。把坐标平移(30000,30000)
        // (x,y)设为一个60001进制数 (因为两边都包括，再加0）。则转为10进制数为
        // 60001(x+30000)+(y+30000)
    }
};

// test for small number
// bug
int intConvXY(int x, int y) {
    return (x+300)*600 + y + 300 ;
    // ( -300<=x,y<=300)
}
// bug
long long intConvXY2(int x, int y) {
    return (x+30000) * 60000ll + y + 30000;
}
// correct
long long intConvXY3(int x, int y) {
    return x * 60001ll + y;
}
//
// x = 102
// y = -2
// x = -2
// y = 102

int main() {
    vector<int> commands = {4,-1,4,-2,4};
    vector<vector<int>> obstacles = {{2,4}};
    {
        int result = Solution::robotSim(commands, obstacles, 1);
        cout << "commands=" << commands << ",obstacles=" << obstacles << ",result=" << result << endl;
        // 输入：commands = [4,-1,4,-2,4], obstacles = [[2,4]]
        // 输出：65
    }
    {
        int result = Solution::robotSim(commands, obstacles, 2);
        cout << "commands=" << commands << ",obstacles=" << obstacles << ",result=" << result << endl;
        // 输入：commands = [4,-1,4,-2,4], obstacles = [[2,4]]
        // 输出：65
    }

    /*
    unordered_map<int,string> a;
    for (int i=-30000; i<=30000; i++) {
        for(int j=-30000; j<=30000; j++) {
            int c = intConvXY2(i,j);
            string xy = to_string(i) + "," + to_string(j);
            if (a.find(c)!= a.end()) {
                cout << "(" << i <<"," << j <<") ==" << a.find(c)->first << "== (" << a.find(c)->second  << ")" << endl;
                break;
            }
            a[c]=xy;
        }
    }
    */
    // bug
    cout << intConvXY2(-29999,-30000) << "==" << intConvXY2(-30000,30000) <<endl;

    return 0;
}







