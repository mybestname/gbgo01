#include <iterator> // needed for std::ostream iterator
#include <iostream>
#include <vector>
template <typename T>
std::ostream& operator<< (std::ostream& out, const std::vector<T>& v) {
    if ( !v.empty() ) {
        out << '[';
        std::copy (v.begin(), v.end(), std::ostream_iterator<T>(out, ","));
        out << "\b]";
    }
    return out;
}
template <typename T>
std::ostream& operator<< (std::ostream& out, const std::vector<std::vector<T>>& v) {
    if ( !v.empty() ) {
        out << "[";
        for (int i= 0; i< v.size(); i++ ){
            std::vector<T> in = v[i];
            if (!in.empty()) {
                out << '[';
                std::copy(in.begin(), in.end(), std::ostream_iterator<T>(out, ","));
                out << "\b],";
            }
            if (i < v.size()-1) {
                out << "\n";
            }
        }
        out << "\n]";
    }
    return out;
}
using namespace std;

class Solution {
public :
    static vector<vector<int>> combine(int n, int k) {
        vector<vector<int>> result = {{}};  // 存放结果
        vector<int> path = {};              // 存放路径

        // 组合问题抽象树形结构
        // n = 4, k = 2
        //
        // 每一个成功的路径，相当于一个从起点到终点，距离为k的路径。起点为 （因为是组合）
        // (4,2) n=4 k=2
        // 起点为1：1->2, 1->3, 1->4
        // 起点为2：2->3, 2->4,
        // 起点为3：3->4
        // 起点为4：4->x
        // - 只需要把达到节点符合条件的路径收集起来可，即 size = k 的路径
        back_tracing(n,k, 1, path, result);
        return result;
    }

    // 优化的方案
    static vector<vector<int>> combine_o(int n, int k) {
        vector<vector<int>> result = {{}};  // 存放结果
        vector<int> path = {};              // 存放路径
        back_tracing_o(n,k, 1, path, result);
        return result;
    }
private:
    static void back_tracing(int n, int k, int start, vector<int> &path, vector<vector<int>> &result) {
        if (path.size() == k) {
            result.push_back(path); // 路径符合条件，加入结果集
            return;                 // 本路径/回溯 结束
        }
        for (int i = start ; i <= n; i++ ){                    // 这里没有优化
            path.push_back(i);                                 // 处理节点，把元素加入path
            back_tracing(n,k,i+1,path, result);    // 递归
            path.pop_back();                                   // 回溯如果结束，那么该路径已经处理，弹出该元素
        }
        // 全部元素处理完，整个结束。
    };

    // 优化的回溯
    static void back_tracing_o(int n, int k, int start, vector<int> &path, vector<vector<int>> &result) {
        cout << "back_tracing s="<< start <<", path=" << path << std::endl;
        if (path.size() == k) {
            result.push_back(path);
            return;
        }
        // 优化：
        // i 没有必要到n
        // 举例 (5,3)
        // 起点为1：1->2 2->3 => 1,2,3
        //             2->4 => 1,2,4
        //             2->5 => 1,2,5
        //        1->3 3->4 => 1,3,4
        //             3->5 => 1,3,5
        //        1->4 4->5 => 1,4,5
        //        1->5->x  (没必要1->5)
        // 起点为2：2->3 3->4 => 2,3,4
        //             3->5 => 2,3,5
        //        2->4 4->5 => 2,4,5
        //        2->5->x  (没必要2->5)
        // 起点为3：3->4 4->5 => 3,4,5
        //        3->5->x  (没必要3->5)
        // 起点为4：4->5->x (没必要4->)
        // 起点为5：5->x    (没必要5->)
        // 所以：path=1时候，起点没有必要是4，5。
        //      path=2时候，起点没有必要是5。
        // 总结为：
        //     i <= n - (k - path.size()) + 1
//      printf("for i=%d to %d, path_size=%lu, why=%lu\n", start, n, path.size(), n - (k - path.size()) + 1);
        for (int i = start;  i <= n - (k - path.size()) + 1; i++) {
//      for (int i = start ; i <= n; i++) {
            path.push_back(i);
            back_tracing_o(n, k, i + 1, path, result);
            path.pop_back();
//------------------------------------------------------      ----------------------------------------------------------
//                                              i  <= n   vs. i <=  n - (k - path.size()) + 1
//------------------------------------------------------      ----------------------------------------------------------
// back_tracing s=1, path=                                 |   back_tracing s=1, path=
// for i=1 to 5, path_size=0, why=3                        |   for i=1 to 5, path_size=0, why=3
//     back_tracing s=2, path=[1]                          |       back_tracing s=2, path=[1]
//       for i=2 to 5, path_size=1, why=4                  |           for i=2 to 5, path_size=1, why=4
//           back_tracing s=3, path=[1,2]                  |               back_tracing s=3, path=[1,2]
//               for i=3 to 5, path_size=2, why=5          |                   for i=3 to 5, path_size=2, why=5
//                    back_tracing s=4, path=[1,2,3]       |                       back_tracing s=4, path=[1,2,3]
//                    back_tracing s=5, path=[1,2,4]       |                       back_tracing s=5, path=[1,2,4]
//                    back_tracing s=6, path=[1,2,5]       |                       back_tracing s=6, path=[1,2,5]
//           back_tracing s=4, path=[1,3]                  |               back_tracing s=4, path=[1,3]
//               for i=4 to 5, path_size=2, why=5          |                   for i=4 to 5, path_size=2, why=5
//                   back_tracing s=5, path=[1,3,4]        |                       back_tracing s=5, path=[1,3,4]
//                   back_tracing s=6, path=[1,3,5]        |                       back_tracing s=6, path=[1,3,5]
//           back_tracing s=5, path=[1,4]                  |               back_tracing s=5, path=[1,4]
//               for i=5 to 5, path_size=2, why=5          |                   for i=5 to 5, path_size=2, why=5
//                   back_tracing s=6, path=[1,4,5]        |                       back_tracing s=6, path=[1,4,5]
//           back_tracing s=6, path=[1,5]                  |       //      path=2, 1->5 没有必要再走
//               for i=6 to 5, path_size=2, why=5          |
//     back_tracing s=3, path=[2]                          |       back_tracing s=3, path=[2]
//       for i=3 to 5, path_size=1, why=4                  |           for i=3 to 5, path_size=1, why=4
//           back_tracing s=4, path=[2,3]                  |               back_tracing s=4, path=[2,3]
//               for i=4 to 5, path_size=2, why=5          |                   for i=4 to 5, path_size=2, why=5
//                   back_tracing s=5, path=[2,3,4]        |                       back_tracing s=5, path=[2,3,4]
//                   back_tracing s=6, path=[2,3,5]        |                       back_tracing s=6, path=[2,3,5]
//           back_tracing s=5, path=[2,4]                  |               back_tracing s=5, path=[2,4]
//               for i=5 to 5, path_size=2, why=5          |                   for i=5 to 5, path_size=2, why=5
//                   back_tracing s=6, path=[2,4,5]        |                       back_tracing s=6, path=[2,4,5]
//           back_tracing s=6, path=[2,5]                  |       //      path=2, 2->5 没有必要再走
//               for i=6 to 5, path_size=2, why=5          |
//     back_tracing s=4, path=[3]                          |       back_tracing s=4, path=[3]
//       for i=4 to 5, path_size=1, why=4                  |           for i=4 to 5, path_size=1, why=4
//           back_tracing s=5, path=[3,4]                  |               back_tracing s=5, path=[3,4]
//               for i=5 to 5, path_size=2, why=5          |                   for i=5 to 5, path_size=2, why=5
//                   back_tracing s=6, path=[3,4,5]        |                       back_tracing s=6, path=[3,4,5]
//           back_tracing s=6, path=[3,5]                  |       //        path=2, 3->5 没有必要再走
//               for i=6 to 5, path_size=2, why=5          |       // path=1, 4->  没有必要再走
//     back_tracing s=5, path=[4]                          |
//       for i=5 to 5, path_size=1, why=4                  |
//           back_tracing s=6, path=[4,5]                  |
//               for i=6 to 5, path_size=2, why=5          |
//     back_tracing s=6, path=[5]                          |       // path=1, 5->  没有必要在走
//       for i=6 to 5, path_size=1, why=4                  |
//                                                         |
      }
        // 全部元素处理完，整个结束。
    };
};

int main() {
    Solution sol;
    int n = 5;
    int k = 3;
    vector<vector<int>> result = Solution::combine(n, k);
    cout << result << endl;
    vector<vector<int>> result_o = Solution::combine_o(n, k);
    cout << result_o << endl;

    return 0;
}
