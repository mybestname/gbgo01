#include <iostream>
#include <vector>
#include <stack>
#include "../../../base/algo_base.h"

using namespace std;
class Solution {
public:
    static vector<int> twoSum(vector<int>& numbers, int target) {
        // 暴力法
        /*
        for (int i= 0; i < numbers.size()-1; i++) {        // i: 0 ~ n-2
            for (int j= i+1; j <= numbers.size()-1; j++) { // j: 1 ~ n-1
               if (numbers[i] + numbers[j] == target) {
                   return {i+1,j+1};                       // 1~n-1, 2~n
               }
            }
        }
        return {};
        */
        // 暴力2
        /*
        uint n = numbers.size();
        for (int j = 1; j < n ; j++) {       // j: 1 ~ n-1
            for (int i = 0; i < n-1; i++) {  // i: 0 ~ n-2
                if (numbers[i]+numbers[j] == target) {
                    return {i+1, j+1};       // {[1,n-1],[2,n]}
                }
            }
        }
        return {};
        */
        // 观察：
        // 1. 固定i，看j如何改变：找到j，使得nums[j] = target - nums[i]
        // 2. 移动i，看j如何改变：因为数组有序，所以i增加，代表nums[i]增加，而target固定不变，所以nums[j]必定减小。又因为数组有序，所以j必定减少。
        //    - 所以i，j二者应该是从两端开始向中间移动。i单调递增，j单调递减。
        //    - 所以不需要枚举，可以从上一个位置向下一个位置继续向先走（不需要掉头）
        //    - 这样i走一遍时候，j也走一遍。中间相遇。整体变为O(n)
        //    - 归纳为可以使用双指针，中间相遇的一类扫描题目。
        int j = numbers.size()-1;  // j从尾端向中间移动，设开始为n-1；
        for (int i = 0; i<numbers.size(); i++) { //i的外层枚举不变，目标是去掉内层循环
            while(i<j && numbers[i]+numbers[j] > target){ //没有相遇，且值还是大（注意这里对于单调性的使用。一个增加，一个减少）
                j--;  //注意，这个while条件只是单纯减小j，注意j是单调下降的。numbers[j]也是单调下降的。
                // 这里是说，针对任何一个固定的i，j需要不需要减小。i在外层枚举，是一个单调增的状态。
                // 换句话说，是一个两头向中间逼近的效果。外层代表左端，每个迭代向中心逼近一步。而while控制右端。
                // 如果两端加起来还小于target，那么右端是不动的，让外层迭代去增加左端的值。直到两端加起来大于target
                // 此时左端一定是固定的状态，所以把右端向中心逼近一步。
                // 这样左右逼近，直到找到左右顿之和==target的情况。
            }
            if (i<j && numbers[i]+numbers[j] == target) { // 注意 i<j 的判断
                return {i+1,j+1};
            }
        }
        return {};
        // 时间复杂度分析：
        // 虽然在for里面有一个while，但是j只减小不增大。所以整体最坏还是O(n)
    }
    // 写法2
    static vector<int> twoSum2(vector<int>& numbers, int target) {
        int j = numbers.size()-1;
        for (int i = 0; i<j; i++) {
            while(numbers[i]+numbers[j] > target) j--;
            if (numbers[i]+numbers[j] == target) return {i+1, j+1};
        }
        return {};
    }
    // 写法3：注意！这是一种可能会引起错误的写法！！！！只能算第一个找到的答案
    static vector<int> twoSum3(vector<int>& numbers, int target) {
        int i = 0;
        int j = numbers.size()-1;
        while(i<j) {
            if (numbers[i]+numbers[j] < target) i++;
            else if (numbers[i] + numbers[j] > target) j--;
            else return {i+1, j+1};
            //注意，当找到第一个答案就返回了。但是如果想找到所以的答案，这种写法是无法完成的！！！，因为无法保证有一个端点一定增加。
            //所以，会重复循环同样的答案！！！！！
            //而如果外层的for循环，则保证一个端点一定增加。就能继续寻找下一个答案。
        }
        return {};
    }

};

int main() {
    vector<int> nums = {5,25,75};
    int target = 100;
    vector<int> result = Solution::twoSum(nums, target);
    cout << "nums=" << nums << ", target=" << target <<", result=" << result << endl;
    // 输入：numbers = [2,7,11,15], target = 9
    // 输出：[1,2]
    return 0;
}
