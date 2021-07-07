#include <iostream>
#include <vector>
#include <unordered_map>
#include "../../../base/algo_base.h"

using namespace std;
class LRUCache {
public:
    LRUCache(int capacity) {
        this->capacity = capacity;
        // 建立保护节点
        this->head = Node{};
        this->tail = Node{};
        this->head.next = &tail;
        this->tail.next = &head;
    }

    int get(int key) {
        if (map.find(key)== map.end()) return -1;
        // 从链表中删除
        Node* node = map[key];
        removeFromList(node);
        // 并重新插入到头部
        insertToHead(node->key, node->val);
        return node->val;
    }

    void put(int key, int value) {
        if (map.find(key)!=map.end()) {
            Node* node = map[key];
            removeFromList(node);
            insertToHead(key, value);
        }else {
            Node* node = insertToHead(key, value);
        }
        if (map.size() > capacity) {
            removeFromList(tail.pre);
        }
    }

private:
    struct Node {
        int key;
        int val;
        Node* pre;
        Node* next;
    } head, tail; //保护节点
    int capacity;
    unordered_map<int,Node*> map;

    void removeFromList(Node *node){
        node->pre->next = node->next;
        node->next->pre = node->pre;
        map.erase(node->key);
        delete node;
    }
    Node* insertToHead(int key, int value) {
        Node* node = new Node();
        node->key = key;
        node->val = value;
        node->next = head.next;
        head.next->pre = node;
        node->pre = &head;
        head.next = node;
        map[node->key] = node;
        return node;
    }
};

int main() {
    LRUCache lRUCache(2);
    lRUCache.put(1, 1); // 缓存是 {1=1}
    lRUCache.put(2, 2); // 缓存是 {1=1, 2=2}
    cout << "expect  1, got " << lRUCache.get(1) << endl;    // 返回 1
    lRUCache.put(3, 3); // 该操作会使得关键字 2 作废，缓存是 {1=1, 3=3}
    cout << "expect -1, got " << lRUCache.get(2) << endl;;    // 返回 -1 (未找到)
    lRUCache.put(4, 4); // 该操作会使得关键字 1 作废，缓存是 {4=4, 3=3}
    cout << "expect -1, got " << lRUCache.get(1) << endl;    // 返回 -1 (未找到)
    cout << "expect  3, got " << lRUCache.get(3) << endl;    // 返回 3
    cout << "expect  4, got " << lRUCache.get(4) << endl;    // 返回 4
    return 0;
}

