# Redis 

## Install

```
git clone git@github.com:redis/redis.git
cd redis
git checkout 6.2.4
make distclean
make
make test
```

## Prepare
### Server 
```
$ ./src/redis-server
45899:C 11 Jun 2021 06:00:42.943 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
45899:C 11 Jun 2021 06:00:42.943 # Redis version=6.2.4, bits=64, commit=f9470c9a, modified=0, pid=45899, just started
45899:C 11 Jun 2021 06:00:42.943 # Warning: no config file specified, using the default config. In order to specify a config file use ./src/redis-server /path/to/redis.conf
45899:M 11 Jun 2021 06:00:42.944 * monotonic clock: POSIX clock_gettime
                _._
           _.-``__ ''-._
      _.-``    `.  `_.  ''-._           Redis 6.2.4 (f9470c9a/0) 64 bit
  .-`` .-```.  ```\/    _.,_ ''-._
 (    '      ,       .-`  | `,    )     Running in standalone mode
 |`-._`-...-` __...-.``-._|'` _.-'|     Port: 6379
 |    `-._   `._    /     _.-'    |     PID: 45899
  `-._    `-._  `-./  _.-'    _.-'
 |`-._`-._    `-.__.-'    _.-'_.-'|
 |    `-._`-._        _.-'_.-'    |           https://redis.io
  `-._    `-._`-.__.-'_.-'    _.-'
 |`-._`-._    `-.__.-'    _.-'_.-'|
 |    `-._`-._        _.-'_.-'    |
  `-._    `-._`-.__.-'_.-'    _.-'
      `-._    `-.__.-'    _.-'
          `-._        _.-'
              `-.__.-'

45899:M 11 Jun 2021 06:00:42.947 # Server initialized
45899:M 11 Jun 2021 06:00:42.948 * Ready to accept connections
```
### Check Server is working
```
$ ./src/redis-cli ping
PONG
```

## Benchmark

### Benchmark Task
- 1、使用 redis benchmark 工具
   - 测试 10 20 50 100 200 1k 5k 字节 value 大小
   - redis get set 性能。

- 2、写入一定量的 kv 数据, 
   - 根据数据大小 1w-50w 自己评估, 
   - 结合写入前后的 info memory 信息  
   - 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

### System 

```
Processor Name:         	Quad-Core Intel Core i7
Processor Speed:        	3.1 GHz
Number of Processors:   	1
Total Number of Cores:    	4
L2 Cache (per Core):    	256 KB
L3 Cache:               	8 MB
Hyper-Threading Technology:	Enabled
Memory:                 	16 GB (2133 MHz LPDDR3)
```
### Task 1
```
input:
    -c client : use default 50
    -d 10 20 50 100 200 1000 5000
    -t SET/GET
    -n requests, default 100000, use 1000000 
    -r random key, keyspace same with request: 1m
```
```
$ for d in 10 20 50 100 200 1000 5000; do ./src/redis-benchmark -d $d -r 1000000 -n 1000000 -t set,get -q; done
SET: 58343.06 requests per second, p50=0.623 msec
GET: 80327.73 requests per second, p50=0.487 msec

SET: 78094.49 requests per second, p50=0.511 msec
GET: 91457.84 requests per second, p50=0.487 msec

SET: 91390.97 requests per second, p50=0.487 msec
GET: 90735.87 requests per second, p50=0.487 msec

SET: 89597.70 requests per second, p50=0.487 msec
GET: 92293.49 requests per second, p50=0.487 msec

SET: 84203.43 requests per second, p50=0.479 msec
GET: 91482.94 requests per second, p50=0.479 msec

SET: 84688.35 requests per second, p50=0.495 msec
GET: 92859.13 requests per second, p50=0.487 msec

SET: 88873.09 requests per second, p50=0.495 msec
GET: 91365.92 requests per second, p50=0.487 msec
```

```
$ for d in 10 20 50 100 200 1000 5000; do ./src/redis-benchmark -d $d -r 1000000 -n 1000000 -t set,get -q; done
SET: 96880.45 requests per second, p50=0.455 msec
GET: 87039.77 requests per second, p50=0.503 msec

SET: 92276.46 requests per second, p50=0.479 msec
GET: 92336.10 requests per second, p50=0.495 msec

SET: 92182.89 requests per second, p50=0.479 msec
GET: 92764.38 requests per second, p50=0.487 msec

SET: 91157.70 requests per second, p50=0.479 msec
GET: 91157.70 requests per second, p50=0.495 msec

SET: 92781.59 requests per second, p50=0.479 msec
GET: 90383.23 requests per second, p50=0.495 msec

SET: 90843.02 requests per second, p50=0.503 msec
GET: 91116.17 requests per second, p50=0.495 msec

SET: 69628.18 requests per second, p50=0.591 msec
GET: 79936.05 requests per second, p50=0.527 msec
```



#### with pipeline
```
$ for d in 10 20 50 100 200 1000 5000  10 20 50 100 200 1000 5000 ; do echo -d=$d ; ./src/redis-benchmark -d $d -r 1000000 -n 1000000 -t set,get -q -P 16; done
-d=10
SET: 97389.95 requests per second, p50=5.815 msec
GET: 312109.88 requests per second, p50=2.127 msec

-d=20
SET: 325309.06 requests per second, p50=2.231 msec
GET: 518134.72 requests per second, p50=1.311 msec

-d=50
SET: 357015.34 requests per second, p50=2.023 msec
GET: 558035.69 requests per second, p50=1.223 msec

-d=100
SET: 246974.58 requests per second, p50=2.127 msec
GET: 548245.62 requests per second, p50=1.247 msec

-d=200
SET: 353481.78 requests per second, p50=1.991 msec
GET: 575705.25 requests per second, p50=1.207 msec

-d=1000
SET: 250626.56 requests per second, p50=1.959 msec
GET: 527426.12 requests per second, p50=1.311 msec

-d=5000
SET: 135189.94 requests per second, p50=1.519 msec
GET: 255297.42 requests per second, p50=2.759 msec

-d=10
SET: 253485.42 requests per second, p50=2.879 msec
GET: 469924.81 requests per second, p50=1.503 msec

-d=20
SET: 312304.81 requests per second, p50=2.359 msec
GET: 579374.31 requests per second, p50=1.199 msec

-d=50
SET: 342935.53 requests per second, p50=2.135 msec
GET: 595238.12 requests per second, p50=1.167 msec

-d=100
SET: 364564.34 requests per second, p50=1.999 msec
GET: 590318.75 requests per second, p50=1.183 msec

-d=200
SET: 347101.69 requests per second, p50=2.103 msec
GET: 579038.81 requests per second, p50=1.207 msec

-d=1000
SET: 263296.47 requests per second, p50=1.807 msec
GET: 510725.25 requests per second, p50=1.359 msec

-d=5000
SET: 158403.30 requests per second, p50=1.303 msec
GET: 272628.12 requests per second, p50=2.551 msec

```

### Task2

```
$ for d in 10 20 50 100 200 1000 5000  10 20 50 100 200 1000 5000 ; do echo -d=$d ; ./src/redis-benchmark -d $d -r 1000000 -n 1000000 -t set,get -q -P 16; ./src/redis-cli info memory|grep human; done
-d=10
SET: 518403.31 requests per second, p50=1.303 msec
GET: 711237.56 requests per second, p50=0.983 msec

used_memory_human:66.91M
used_memory_rss_human:70.12M
used_memory_peak_human:69.51M
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=20
SET: 531067.44 requests per second, p50=1.343 msec
GET: 661813.38 requests per second, p50=1.055 msec

used_memory_human:97.87M
used_memory_rss_human:103.34M
used_memory_peak_human:100.43M
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=50
SET: 446827.53 requests per second, p50=1.623 msec
GET: 633713.56 requests per second, p50=1.103 msec

used_memory_human:120.71M
used_memory_rss_human:134.96M
used_memory_peak_human:123.27M
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=100
SET: 410509.00 requests per second, p50=1.767 msec
GET: 604960.69 requests per second, p50=1.159 msec

used_memory_human:165.88M
used_memory_rss_human:184.21M
used_memory_peak_human:168.44M
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=200
SET: 376081.25 requests per second, p50=1.935 msec
GET: 587889.50 requests per second, p50=1.183 msec

used_memory_human:240.32M
used_memory_rss_human:281.00M
used_memory_peak_human:242.88M
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=1000
SET: 256147.53 requests per second, p50=2.039 msec
GET: 511770.72 requests per second, p50=1.359 msec

used_memory_human:750.35M
used_memory_rss_human:882.49M
used_memory_peak_human:752.91M
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=5000
SET: 132240.16 requests per second, p50=1.591 msec
GET: 262123.19 requests per second, p50=2.655 msec

used_memory_human:3.34G
used_memory_rss_human:3.92G
used_memory_peak_human:3.34G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=10
SET: 266311.56 requests per second, p50=2.719 msec
GET: 474833.81 requests per second, p50=1.487 msec

used_memory_human:1.29G
used_memory_rss_human:3.92G
used_memory_peak_human:3.34G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=20
SET: 330033.00 requests per second, p50=2.223 msec
GET: 588235.25 requests per second, p50=1.183 msec

used_memory_human:559.15M
used_memory_rss_human:3.92G
used_memory_peak_human:3.34G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=50
SET: 355239.78 requests per second, p50=2.063 msec
GET: 597371.56 requests per second, p50=1.167 msec

used_memory_human:297.57M
used_memory_rss_human:3.92G
used_memory_peak_human:3.34G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=100
SET: 400480.59 requests per second, p50=1.799 msec
GET: 578703.69 requests per second, p50=1.207 msec

used_memory_human:333.60M
used_memory_rss_human:3.92G
used_memory_peak_human:3.34G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=200
SET: 372856.09 requests per second, p50=1.943 msec
GET: 587199.06 requests per second, p50=1.183 msec

used_memory_human:302.19M
used_memory_rss_human:3.92G
used_memory_peak_human:3.34G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=1000
SET: 256673.52 requests per second, p50=1.799 msec
GET: 466200.47 requests per second, p50=1.487 msec

used_memory_human:772.60M
used_memory_rss_human:3.92G
used_memory_peak_human:3.34G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
-d=5000
SET: 149097.97 requests per second, p50=1.375 msec
GET: 249750.23 requests per second, p50=2.791 msec

used_memory_human:3.34G
used_memory_rss_human:3.95G
used_memory_peak_human:3.35G
total_system_memory_human:16.00G
used_memory_lua_human:37.00K
used_memory_scripts_human:0B
maxmemory_human:0B
```


#### key size estimated by `used_memory`

1. n = 1000000 (100w) 
```
$ for d in 10 20 50 100 200 1000 5000 ; do n=1000000; ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); k=$(echo "($m/$n)-$d"|bc); echo -e "d = $d \nk = $m/$n - $d = $k";done;

SET: 513347.03 requests per second, p50=1.311 msec

d = 10
k = 70166064/1000000 - 10 = 60
SET: 531067.44 requests per second, p50=1.335 msec

d = 20
k = 102554448/1000000 - 20 = 82
SET: 436109.91 requests per second, p50=1.663 msec

d = 50
k = 134765664/1000000 - 50 = 84
SET: 399361.03 requests per second, p50=1.799 msec

d = 100
k = 176876000/1000000 - 100 = 76
SET: 312109.88 requests per second, p50=2.063 msec

d = 200
k = 253150976/1000000 - 200 = 53
SET: 248200.56 requests per second, p50=1.919 msec

d = 1000
k = 786749280/1000000 - 1000 = -214
SET: 124906.33 requests per second, p50=1.687 msec

d = 5000
k = 3581632496/1000000 - 5000 = -1419
```

2. n = 5000000 (50w)
```
$ for d in 10 20 50 100 200 1000 5000 ; do n=500000; ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); k=$(echo "($m/$n)-$d"|bc); echo -e "d = $d \nk = $m/$n - $d = $k";done;
SET: 531914.88 requests per second, p50=1.271 msec

d = 10
k = 35588816/500000 - 10 = 61
SET: 533617.94 requests per second, p50=1.335 msec

d = 20
k = 51808832/500000 - 20 = 83
SET: 432152.12 requests per second, p50=1.663 msec

d = 50
k = 67882144/500000 - 50 = 85
SET: 422297.28 requests per second, p50=1.719 msec

d = 100
k = 89012192/500000 - 100 = 78
SET: 378214.81 requests per second, p50=1.911 msec

d = 200
k = 127118576/500000 - 200 = 54
SET: 273822.56 requests per second, p50=1.767 msec

d = 1000
k = 394208944/500000 - 1000 = -212
SET: 137779.00 requests per second, p50=1.543 msec

d = 5000
k = 1792397824/500000 - 5000 = -1416
```
3. n = 100000  (10w)
```
$ for d in 10 20 50 100 200 1000 5000 ; do n=100000; ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); k=$(echo "($m/$n)-$d"|bc); echo -e "d = $d \nk = $m/$n - $d = $k";done;
SET: 588235.31 requests per second, p50=1.167 msec

d = 10
k = 7663536/100000 - 10 = 66
SET: 555555.56 requests per second, p50=1.247 msec

d = 20
k = 11423008/100000 - 20 = 94
SET: 480769.22 requests per second, p50=1.495 msec

d = 50
k = 14645744/100000 - 50 = 96
SET: 448430.47 requests per second, p50=1.591 msec

d = 100
k = 18869184/100000 - 100 = 88
SET: 386100.38 requests per second, p50=1.759 msec

d = 200
k = 26494192/100000 - 200 = 64
SET: 259740.27 requests per second, p50=1.799 msec

d = 1000
k = 79784976/100000 - 1000 = -203
SET: 147492.62 requests per second, p50=1.447 msec

d = 5000
k = 359648864/100000 - 5000 = -1404
```

4. n = 10000   (1w)

```
$ for d in 10 20 50 100 200 1000 5000 ; do n=10000; ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); k=$(echo "($m/$n)-$d"|bc); echo -e "d = $d \nk = $m/$n - $d = $k";done;
SET: 416666.66 requests per second, p50=1.575 msec

d = 10
k = 1737392/10000 - 10 = 163
SET: 434782.59 requests per second, p50=1.479 msec

d = 20
k = 2132480/10000 - 20 = 193
SET: 384615.38 requests per second, p50=1.631 msec

d = 50
k = 2450160/10000 - 50 = 195
SET: 370370.38 requests per second, p50=1.783 msec

d = 100
k = 2867712/10000 - 100 = 186
SET: 416666.66 requests per second, p50=1.671 msec

d = 200
k = 3631456/10000 - 200 = 163
SET: 270270.28 requests per second, p50=1.735 msec

d = 1000
k = 9014608/10000 - 1000 = -99
SET: 147058.81 requests per second, p50=1.431 msec

d = 5000
k = 36811344/10000 - 5000 = -1319

```

#### key size estimated by `used_memory_peak` 

n=100w

```shell
$ for d in 10 20 50 100 200 1000 5000 ; do n=1000000; ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory_peak:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); k=$(echo "($m/$n)-$d"|bc); echo -e "d = $d \nk = $m/$n - $d = $k";done;
SET: 513610.69 requests per second, p50=1.311 msec

d = 10
k = 72907776/1000000 - 10 = 62
SET: 524658.94 requests per second, p50=1.359 msec

d = 20
k = 105271920/1000000 - 20 = 85
SET: 434027.78 requests per second, p50=1.671 msec

d = 50
k = 137410272/1000000 - 50 = 87
SET: 410509.00 requests per second, p50=1.767 msec

d = 100
k = 179609408/1000000 - 100 = 79
SET: 373273.62 requests per second, p50=1.959 msec

d = 200
k = 255774032/1000000 - 200 = 55
SET: 257400.27 requests per second, p50=1.895 msec

d = 1000
k = 789628096/1000000 - 1000 = -211
SET: 133209.00 requests per second, p50=1.567 msec

d = 5000
k = 3584437552/1000000 - 5000 = -1416
```

### Explain

题意表述不清楚，统计的不是key的size，因为key的size是输入决定的。而是每个object的占用空间。

#### Estimated object size

1.d=10
```bash
d=10 n=100000 && ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory_peak:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); o=$(echo "($m/$n)"|bc); echo -e "d = $d \no = $m/$n  = $o";
SET: 578034.69 requests per second, p50=1.183 msec

d = 10
o = 10398304/100000  = 103
```

2. d=20
```shell
$ d=20 n=100000 && ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory_peak:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); o=$(echo "($m/$n)"|bc); echo -e "d = $d \no = $m/$n  = $o";
SET: 584795.31 requests per second, p50=1.183 msec

d = 20
o = 11378160/100000  = 113
```

3. d=50
```shell
$ d=50 n=100000 && ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory_peak:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); o=$(echo "($m/$n)"|bc); echo -e "d = $d \no = $m/$n  = $o";
SET: 476190.50 requests per second, p50=1.455 msec

d = 50
o = 13444272/100000  = 134
```

4. d=100
```shell
d = 100
o = 16423040/100000  = 164
```

5. d=200
```shell
$ d=200 n=100000 && ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory_peak:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); o=$(echo "($m/$n)"|bc); echo -e "d = $d \no = $m/$n  = $o";
SET: 458715.59 requests per second, p50=1.543 msec
d = 200
o = 22546752/100000  = 225
```

6. d=1000
```shell
$ d=1000 n=100000 && ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory_peak:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); o=$(echo "($m/$n)"|bc); echo -e "d = $d \no = $m/$n  = $o";
SET: 333333.31 requests per second, p50=1.391 msec

d = 1000
o = 73050976/100000  = 730
```

7. d=5000
```shell
$ d=5000 n=100000 && ./src/redis-benchmark -d $d -r $n -n $n -t set -q -P 16; used=$(./src/redis-cli info memory|grep used_memory_peak:); m=$(echo $used|tr '\r' ':'|cut -d':' -f2); o=$(echo "($m/$n)"|bc); echo -e "d = $d \no = $m/$n  = $o";
SET: 168067.22 requests per second, p50=1.215 msec
d = 5000
o = 333154144/100000  = 3331
```


## Footnotes

### Redis persistence 

RDB (Redis Database): The RDB persistence performs point-in-time snapshots of your dataset at specified intervals.
- https://redis.io/topics/persistence
```
$ du -sh ./temp-50352.rdb
3.1G    ./temp-50352.rdb
$ du -sh ./temp-50352.rdb
du: ./temp-50352.rdb: No such file or directory
$ du -sh ./dump.rdb
3.2G    ./dump.rdb
```
### How to disable all persistence

`redis-server --save "" --appendonly no`

```c 
/* Return true if this instance has persistence completely turned off:
 * both RDB and AOF are disabled. */
int allPersistenceDisabled(void) {
    return server.saveparamslen == 0 && server.aof_state == AOF_OFF;
}
```

### Some useful client commands for checking redis status
```shell
$./src/redis-cli stat
$./src/redis-cli info memory
$./src/redis-cli info keyspace
$./src/redis-cli --bigkeys
$./src/redis-cli --memkeys
#  --bigkeys          Sample Redis keys looking for keys with many elements (complexity).
#  --memkeys          Sample Redis keys looking for keys consuming a lot of memory.
$./src/redis-cli --scan | head -10 
$./src/redis-cli --scan  --pattern 'key:00000000000[123]'
"key:000000000001"
"key:000000000003"
"key:000000000002"
$ ./src/redis-cli --quoted-input get "key:000000000001"
"VXKeHogKgJ"
$ ./src/redis-cli --quoted-input get "key:000000000002"
"VXKeHogKgJ=[5V9_X^b?"
```


## Reference
- https://redis.io/topics/quickstart
- https://redis.io/topics/benchmarks
- https://github.com/redis/redis-doc/blob/master/topics/introduction.md
  - https://github.com/redis/redis-doc/blob/master/topics/data-types.md
  - https://github.com/redis/redis-doc/blob/master/topics/data-types-intro.md
