# Map faster than go map

This code is the result of a successful attempt to implement a map faster than go map of go 1.24.3. It is a map of string keys with int values. I have also a fast map with int keys but it isn’t published. I focussed on string keys as it was closer to my use case.

My hash table is in altmap. stdmap is just a place holder for the standard go map benchmark. The struct type is called Cache as it was initially designed to be used for a cache.

My map is a swiss table like the go map. It is an extensible hash table of 256 group tables where each group has room for 8 items as it is implemented in pure go without simd instructions. Use of simd instructions should speed things up.

The hash function I use is xxh3. On amd64 xxh3 uses assembly with simd instructions for long strings. On arm64 there is currently no assembly, and it doesn’t use simd. It is fast for small strings and 4 times slower than maphash for long strings. The benchmarks are performed with short (8 byte) strings.

The main difference between my map and Go map is that go map uses 7 bit top hashes and I use 8 bit top hashes. It have also additional minor optimizations. On amd64 architectures, go map uses inlined simd instructions which I can't.

My tables split when they reach 90% load. The table size can be adjusted with the constant parameter `tableSizeLog2`, and the maximum load is adjusted with the `maxUsed` constant. A table rehash is also triggered when the number of tombstones reach 15%. This may be adjusted with the `maxTombstones` constant. The table rehash removes all tombstones.

Item deletion replaces them with tombstones. Items aren't thus moved in the groups or table. The hash table should then be compatible with the go map iteration policy. Iteration code is not exposed as it seam to be a tricky operation.

## Benchmarking

To regenerate the benchmark data, go into the stdmap and altmap directories and execute the basj script `./bench.sh Cache2` in each one of them. This will generate a file named `stats_arm64.txt` or `stats_amd64.txt` depending on your current architecture. When done, call `benchstat altmap/stats_arm64.txt stdmap/stats_arm64.txt` to view the stats. The `benchstat` command may be installed by executing `go install golang.org/x/perf/cmd/benchstat@latest`.

You should the see something like this:

```text
goos: darwin
goarch: arm64
pkg: fastmap/map
cpu: Apple M2
                      │ altmap/stats_arm64.txt │         stdmap/stats_arm64.txt         │
                      │         sec/op         │    sec/op      vs base                 │
Cache2Hit/_______1-8              6.377n ±  1%    7.081n ±  0%   +11.04% (p=0.000 n=10)
Cache2Hit/______10-8              8.811n ±  0%    8.114n ± 30%         ~ (p=0.481 n=10)
Cache2Hit/_____100-8              8.591n ±  4%   15.720n ±  8%   +82.99% (p=0.000 n=10)
Cache2Hit/____1000-8              8.698n ±  2%   13.715n ± 23%   +57.68% (p=0.000 n=10)
Cache2Hit/___10000-8              10.84n ±  1%    20.38n ±  6%   +87.96% (p=0.000 n=10)
Cache2Hit/__100000-8              12.96n ±  1%    17.62n ± 17%   +35.96% (p=0.000 n=10)
Cache2Hit/_1000000-8              42.46n ±  0%    55.45n ±  0%   +30.61% (p=0.000 n=10)
Cache2Hit/10000000-8              57.38n ±  4%    68.12n ±  3%   +18.73% (p=0.000 n=10)
Cache2Miss/_______1-8             5.011n ±  0%    8.243n ±  0%   +64.50% (p=0.000 n=10)
Cache2Miss/______10-8             5.892n ± 28%    9.238n ± 11%   +56.79% (p=0.000 n=10)
Cache2Miss/_____100-8             6.028n ±  6%   12.600n ±  9%  +109.02% (p=0.000 n=10)
Cache2Miss/____1000-8             7.182n ±  6%   12.015n ± 17%   +67.29% (p=0.000 n=10)
Cache2Miss/___10000-8             10.50n ±  2%    17.66n ±  9%   +68.16% (p=0.000 n=10)
Cache2Miss/__100000-8             20.96n ±  1%    26.56n ±  2%   +26.69% (p=0.000 n=10)
Cache2Miss/_1000000-8             33.62n ±  1%    36.83n ±  1%    +9.56% (p=0.000 n=10)
Cache2Miss/10000000-8             48.65n ±  2%    50.84n ±  1%    +4.50% (p=0.000 n=10)
geomean                           13.00n          18.38n         +41.46%
```

On amd64, the benchmarks will be on par with go map. But keep in mind that on amd64, go
map uses simd instructions and altmap is in pure go. Here is an example output.

```text
goos: linux
goarch: amd64
pkg: fastmap/map
cpu: 11th Gen Intel(R) Core(TM) i5-11400 @ 2.60GHz
                       │ altmap/stats_amd64.txt │        stdmap/stats_amd64.txt         │
                       │         sec/op         │    sec/op      vs base                │
Cache2Hit/_______1-12               9.325n ± 0%    7.486n ±  0%  -19.72% (p=0.000 n=10)
Cache2Hit/______10-12               9.313n ± 1%    8.594n ±  0%   -7.73% (p=0.000 n=10)
Cache2Hit/_____100-12               9.427n ± 2%   10.290n ±  4%   +9.16% (p=0.002 n=10)
Cache2Hit/____1000-12               10.74n ± 1%    11.06n ±  1%   +2.93% (p=0.000 n=10)
Cache2Hit/___10000-12               15.62n ± 2%    15.99n ±  2%   +2.40% (p=0.001 n=10)
Cache2Hit/__100000-12               21.34n ± 3%    21.91n ±  3%        ~ (p=0.066 n=10)
Cache2Hit/_1000000-12               77.37n ± 2%    76.78n ±  2%   -0.76% (p=0.015 n=10)
Cache2Hit/10000000-12               95.26n ± 4%    92.19n ±  2%   -3.22% (p=0.000 n=10)
Cache2Miss/_______1-12              7.041n ± 1%    8.268n ±  0%  +17.43% (p=0.000 n=10)
Cache2Miss/______10-12              7.045n ± 1%    6.989n ±  8%        ~ (p=0.481 n=10)
Cache2Miss/_____100-12              7.118n ± 1%    8.934n ± 19%  +25.52% (p=0.000 n=10)
Cache2Miss/____1000-12              8.251n ± 3%    9.651n ±  4%  +16.97% (p=0.000 n=10)
Cache2Miss/___10000-12              13.00n ± 1%    14.74n ±  2%  +13.42% (p=0.000 n=10)
Cache2Miss/__100000-12              26.08n ± 2%    30.82n ±  3%  +18.17% (p=0.000 n=10)
Cache2Miss/_1000000-12              62.06n ± 1%    64.69n ±  4%   +4.23% (p=0.000 n=10)
Cache2Miss/10000000-12              81.04n ± 1%    84.81n ±  2%   +4.66% (p=0.000 n=10)
geomean                             17.86n         18.71n         +4.75%
```

## Contributions

Special thanks to Claude AI for its assistance throughout this project.
