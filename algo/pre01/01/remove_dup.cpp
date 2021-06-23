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
using namespace std;

// 注意：
//  - i++ 和 ++i 都是将 i 加 1，
//    - 但是 i++ 返回值为 i
//    -  ++i 返回值为 i+1
//  - 如果只是希望增加i的值，而不需要返回值
//    - 则推荐使用 ++i，其运行速度会略快一些。

class Solution {
public :
    int removeDuplicates(vector<int>& nums) {
        if (nums.empty()) return 0; // check input
        int index = 0; // 游标位置
        for (int i=0; i < nums.size()-1 ; i++){ // 最多检查n-1次
            if (nums[i] != nums [i+1])  {
                //前后不一,游标加一，并移动数据
               ++index;
               nums[index] = nums [i+1];
            }
            // 前后一致，游标不动，检查下一个即可。
        }
        //最后的游标位置，即为长度，因为游标从0开始，所以长度需要加1
        return index+1;
    }
};

int main() {
    Solution sol;
    vector<int> nums = {1,1,2};
    int len = sol.removeDuplicates(nums);
    cout << len << ", nums = " << nums << endl;

    nums = {0,0,1,1,1,2,2,3,3,4};
    len = sol.removeDuplicates(nums);
    cout << len << ", nums = " << nums << endl;
    return 0;
}