#include <iostream>
#include <vector>
#include <unordered_set>
#include "../../../base/algo_base.h"
using namespace std;

class Solution {
public:
    int minEatingSpeed(vector<int>& piles, int h) {
        int minSpeed = 1;
        int maxSpeed = maxPiles(piles)/1;
        // 暴力
        /*
        for (int speed = 1; speed < maxSpeed; speed++) {
            if(canFinish(piles, speed, h)) {
                break;
            }
        }
        */
        // 二分查找
        int left = minSpeed;
        int right = maxSpeed;
        while(left < right) {
            int mid = left + (right-left)/2;
            if(canFinish(piles,mid,h))
                right = mid;
            else {
                left = mid+1;
            }
        }
        return left;
    }

private:
    bool canFinish(vector<int>& piles, int speed, int h){
        int time = 0;
        for(auto p : piles) {
             time += p/speed;
            if (p%speed!=0) time++;
        }
        return time <= h;
    }
    int maxPiles(vector<int>& piles){
        int result = 0;
        for(auto p : piles) {
           result = max(result,p);
        }
        return result;
    }
};


int main() {
    struct Test {
        vector<int> piles;
        int H;
        int expect;
    };
    {
        vector<Test> tests = {
                {.piles = {3,6,7,11},       .H = 8, .expect = 4, },
                {.piles = {30,11,23,4,20},  .H = 5, .expect = 30, },
                {.piles = {30,11,23,4,20},  .H = 6, .expect = 23, },
        };

        Solution s;
        for (auto &test : tests) {
            cout << "in : nums=" << test.piles << ",H=" << test.H << endl;
            auto result = s.minEatingSpeed(test.piles,test.H);
            cout << "out: expect=" << test.expect << ", got=" << result << endl;
        }
    }
}