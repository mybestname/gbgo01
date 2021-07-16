#include <iostream>
#include <vector>
#include <queue>
#include <unordered_set>
#include <unordered_map>
#include "../../../base/algo_base.h"
using namespace std;
class Twitter {
public:
    /** Initialize your data structure here. */
    Twitter() : init_time(0) {
    }

    /** Compose a new tweet. */
    void postTweet(int userId, int tweetId) {
        tweets[userId].push_back(make_pair(init_time++, tweetId));
    }

    /** Follower follows a followee. If the operation is invalid, it should be a no-op. */
    void follow(int followerId, int followeeId) {
        // 入口条件，不能follow自己
        if (followerId == followeeId) return;
        auto& f = follows[followerId];
        // 不能重复follow
        if (f.find(followeeId) == f.end()) {
            f.insert(followeeId);
        }
    }

    /** Follower unfollows a followee. If the operation is invalid, it should be a no-op. */
    void unfollow(int followerId, int followeeId) {
        //入口条件, unfollow自己没有意义
        if (followerId == followeeId) return;
        //两次存在性检查
        if (follows.find(followerId) != follows.end()) {
            auto& f = follows[followerId];
            if (f.find(followeeId) != f.end()) {
               f.erase(followeeId);  // unfollow
            }
        }
    }

    /** Retrieve the 10 most recent tweet ids in the user's news feed. Each item in the news feed must be posted by users who the user followed or by the user herself. Tweets must be ordered from most recent to least recent. */
    vector<int> getNewsFeed(int userId) {
        vector<int> ans;
        //使用优先队列，对时间排序
        priority_queue<pair<timestamp,int>> q;
        //拿到该用户的所有tweet, 加入优先队列
        if (tweets.find(userId) != tweets.end()) {
           for (auto& t : tweets[userId]){
               q.push(t);
           }
        }
        // 拿到关注的用户的所有tweet，加入队列
        for (auto id : follows[userId]) {
            for(auto& t : tweets[id]){
                q.push(t);
            }
        }

        //拿到队列里面的前most_recent_size条 (10条），注意q的empty判断。
        while(!q.empty()&& ans.size()< most_recent_size) {
            auto t = q.top();  //注意这里传copy
            q.pop();
            ans.push_back(t.second);
        }
        return ans;
    }


private:
    typedef int timestamp;
    unordered_map<int, unordered_set<int>> follows;          // user_id -> (user_1, user_2)
    unordered_map<int, vector<pair<timestamp,int>>> tweets;  // user_id -> [<time_1, tweet_id>, <time_2, tweet_id_2> ...]
    // 最近为10条
    const int most_recent_size = 10;
    timestamp init_time;
};

int main() {

    Twitter twitter;

    // 用户1发送了一条新推文 (用户id = 1, 推文id = 5).
    twitter.postTweet(1, 5);

    // 用户1的获取推文应当返回一个列表，其中包含一个id为5的推文.
    auto result = twitter.getNewsFeed(1);
    cout << result << endl;


    // 用户1关注了用户2.
    twitter.follow(1, 2);

    // 用户2发送了一个新推文 (推文id = 6).
    twitter.postTweet(2, 6);

    // 用户1的获取推文应当返回一个列表，其中包含两个推文，id分别为 -> [6, 5].
    // 推文id6应当在推文id5之前，因为它是在5之后发送的.
    result = twitter.getNewsFeed(1);
    cout << result << endl;

    // 用户1取消关注了用户2.
    twitter.unfollow(1, 2);

    // 用户1的获取推文应当返回一个列表，其中包含一个id为5的推文.
    // 因为用户1已经不再关注用户2.
    result = twitter.getNewsFeed(1);
    cout << result << endl;

    return 0;

}
