#include <iostream>
#include <vector>
#include <unordered_set>
#include "../../../base/algo_base.h"
using namespace std;

class Solution {
public:
    int shipWithinDays(vector<int>& weights, int days) {
        int maxCap = sumWeight(weights);
        // 暴力
        /*
        for (int cap=1; cap < maxCap; cap++ ) {
            if (canDeliveryInDays(weights,cap,days)) return cap;
        }
        return maxCap;
        */
        // 二分
        int left = 1;
        int right = maxCap;
        while (left < right) {
           int mid = left + (right - left)/2 ;
           if (canDeliveryInDays(weights,mid,days)) {
               right = mid;
           }else {
               left = mid + 1;
           }
        }
        return left;
    }

private:
    bool canDeliveryInDays(vector<int>& weights, int capacity, int days){
        int takes = 1;
        int sw = 0;
        int cur = 0;
        for (int i= 0; i< weights.size(); i++) {
           if (weights[i] > capacity) return false; //不可能完成
           sw += weights[i];
           if (sw < capacity)  continue;
           else {
               int load = 0;
               while (sw >= 0 && load + weights[cur] <= capacity && cur<=i) {
                   sw -= weights[cur];
                   load += weights[cur];
                   cur++;
               }
               takes++;
           }
        }
        if ( sumWeight(weights) %  capacity == 0) takes--;
        return takes <= days;
    }
    int sumWeight(vector<int>& weights){
        int sw = 0;
        for (auto w : weights) {
            sw += w;
        }
        return sw;
    }
};

int main() {
    struct Test {
        vector<int> weights;
        int D;
        int expect;
    };
    {
        vector<Test> tests = {
                {.weights = {1,2,3,4,5,6,7,8,9,10}, .D = 5, .expect = 15, },
                {.weights = {3,2,2,4,1,4},  .D = 3, .expect = 6, },
                {.weights = {1,2,3,1,1},  .D = 4, .expect = 3, },
                {.weights = {1,2,3,4,5,6,7,8,9,10}, .D = 10, .expect = 10,},
                {.weights = {3,3,3,3,3,3}, .D = 2, .expect = 9,},
        };
        Solution s;
        for (auto &test : tests) {
            cout << "nums=" << test.weights << ",D=" << test.D << endl;
            auto result = s.shipWithinDays(test.weights,test.D);
            cout << "==>> expect=" << test.expect << ", got=" << result << endl;
        }
    }
}

