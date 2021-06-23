#include <iterator> // needed for std::ostream iterator
#include <iostream>
#include <vector>
#include <unordered_map>
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

class Solution {
public :
    vector<int> twoSum(vector<int>& num, int target) {
        // std::unordered_map  key无序 底层实现为哈希表，适合这里使用。
        unordered_map <int,int> map;
        for (int i = 0; i < num.size(); i++) {
            // 在map中查找是否有合适的匹配,
            auto iter = map.find(target - num[i]);
            if (iter != map.end()) {
               // 如果查到，返回 map的v=one的位置，和i，two的位置。
               return {iter->second, i};
            }
            // 如果没有map中记录 key=元素，v=位置
            map[num[i]] = i;
        };
        return {0,0};  //not found
    }
};

int main() {
    Solution sol;

    vector<int> nums = {2,7,11,15};
    int target = 9;
    vector<int> sum = sol.twoSum(nums, target);
    cout << sum << endl;


    nums = {3,2,4};
    target = 6;
    sum = sol.twoSum(nums, target);
    cout << sum << endl;

    nums = {3,3};
    target = 6;
    sum = sol.twoSum(nums, target);
    cout << sum << endl;

    return 0;
}

