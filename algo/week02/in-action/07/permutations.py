class Solution:
    def permute(self, nums: list[int]) -> list[list[int]]:
        self.nums = nums
        self.ans = []
        self.per = []
        self.n = len(nums)
        self.used = [False] * self.n
        self.find(0)
        return self.ans

    # 依次考虑0,1,...,n-1位置放哪个数  
    # “从还没用过的”数中选一个放在当前位置
    def find(self, index):
        if index == self.n:
            self.ans.append(self.per[:])  # make a copy, 因为python是传引用的，如果不先拷贝，那么后续的per.pop()会影响到ans中的per。
            return
        for i in range(self.n):
            if not self.used[i]:
                self.used[i] = True
                self.per.append(self.nums[i])
                self.find(index + 1)
                self.per.pop()
                self.used[i] = False
def main():
    s = Solution()
    p = s.permute([1,2,3])
    print(p)

if __name__ == "__main__":
    main()