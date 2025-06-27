# Map faster than go map

This prove of concept code is the result of a successful attempt to implement a map faster than go map from go 1.24.3. As the standard map, my map is a swiss table but it is using 8 bit top hashes and an exact match. It doesn’t move items and thus uses tombstones for deleted items. Beside the difference of using 8 bit top hashes and an exact match that reduces the false positive, it uses xxh3 as hash function instead of the maphash used by the go map. This maphash function is fast on intel processor, and not so fast on arm64 processors.

The go map and my table use an extensible hash table (directory), but the directory adds an overhead that may be removed when the map is used for a Cache as its size is constant at full regime. I kept the directory in this fast map to display some possible optimizations. The table, made of 256 groups of 8 items, split when reaching 90% load. This size and load threshold may be modified by the constant parameters `tableSizeLog2` and `maxUsed`. A table rehash is also triggered when the number of tombstones reach `maxTombstones` set to 15% of the capacity in this implementation.

The hash table is named Cache as it was initially designed to be used for a cache.

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

The following are the benchmarks of the exacte same map, but with int keys. It uses an xxh3 hash for integers which is currently not provided in the cheebo package.

```text
goos: darwin
goarch: arm64
pkg: fastmap/map
cpu: Apple M2
                      │ altmapint/stats_arm64.txt │      stdmapint/stats_arm64.txt       │
                      │          sec/op           │    sec/op     vs base                │
Cache2Hit/_______1-8                  3.515n ± 1%   2.080n ±  0%  -40.81% (p=0.000 n=10)
Cache2Hit/______10-8                  4.330n ± 1%   5.107n ±  0%  +17.94% (p=0.000 n=10)
Cache2Hit/_____100-8                  4.427n ± 3%   6.139n ±  7%  +38.68% (p=0.000 n=10)
Cache2Hit/____1000-8                  4.401n ± 6%   5.949n ±  1%  +35.17% (p=0.000 n=10)
Cache2Hit/___10000-8                  5.139n ± 1%   7.624n ±  1%  +48.35% (p=0.000 n=10)
Cache2Hit/__100000-8                  6.805n ± 1%   9.424n ±  1%  +38.49% (p=0.000 n=10)
Cache2Hit/_1000000-8                  18.26n ± 1%   26.37n ±  0%  +44.39% (p=0.000 n=10)
Cache2Hit/10000000-8                  28.26n ± 3%   38.13n ±  1%  +34.95% (p=0.000 n=10)
Cache2Miss/_______1-8                 3.013n ± 1%   2.361n ±  1%  -21.64% (p=0.000 n=10)
Cache2Miss/______10-8                 3.751n ± 0%   4.660n ±  1%  +24.23% (p=0.001 n=10)
Cache2Miss/_____100-8                 3.880n ± 4%   6.081n ± 10%  +56.69% (p=0.000 n=10)
Cache2Miss/____1000-8                 4.575n ± 7%   6.033n ±  2%  +31.86% (p=0.000 n=10)
Cache2Miss/___10000-8                 6.577n ± 3%   9.341n ±  2%  +42.03% (p=0.000 n=10)
Cache2Miss/__100000-8                 15.74n ± 1%   20.22n ±  1%  +28.50% (p=0.000 n=10)
Cache2Miss/_1000000-8                 14.85n ± 0%   21.24n ±  0%  +42.98% (p=0.000 n=10)
Cache2Miss/10000000-8                 24.61n ± 1%   33.97n ±  0%  +38.06% (p=0.000 n=10)
geomean                               7.088n        8.897n        +25.52%
```

```text
goarch: amd64
pkg: fastmap/map
cpu: 11th Gen Intel(R) Core(TM) i5-11400 @ 2.60GHz
                       │ altmapint/stats_amd64.txt │       stdmapint/stats_amd64.txt       │
                       │          sec/op           │    sec/op      vs base                │
Cache2Hit/_______1-12                  6.288n ± 0%    2.988n ±  1%  -52.47% (p=0.000 n=10)
Cache2Hit/______10-12                  6.276n ± 1%    5.782n ±  2%   -7.87% (p=0.000 n=10)
Cache2Hit/_____100-12                  6.393n ± 2%    6.440n ±  5%        ~ (p=0.579 n=10)
Cache2Hit/____1000-12                  6.556n ± 2%    6.960n ±  3%   +6.17% (p=0.000 n=10)
Cache2Hit/___10000-12                  7.303n ± 1%    8.101n ±  1%  +10.92% (p=0.000 n=10)
Cache2Hit/__100000-12                  10.45n ± 2%    12.21n ±  1%  +16.95% (p=0.000 n=10)
Cache2Hit/_1000000-12                  31.72n ± 2%    34.54n ±  2%   +8.87% (p=0.000 n=10)
Cache2Hit/10000000-12                  44.18n ± 1%    47.08n ±  3%   +6.55% (p=0.000 n=10)
Cache2Miss/_______1-12                 5.562n ± 0%    3.214n ±  0%  -42.22% (p=0.000 n=10)
Cache2Miss/______10-12                 5.562n ± 0%    5.748n ±  2%   +3.34% (p=0.000 n=10)
Cache2Miss/_____100-12                 5.679n ± 2%    7.580n ± 13%  +33.47% (p=0.000 n=10)
Cache2Miss/____1000-12                 6.204n ± 2%    7.530n ±  3%  +21.39% (p=0.000 n=10)
Cache2Miss/___10000-12                 8.021n ± 2%   10.210n ±  1%  +27.28% (p=0.000 n=10)
Cache2Miss/__100000-12                 18.99n ± 1%    23.98n ±  1%  +26.24% (p=0.000 n=10)
Cache2Miss/_1000000-12                 21.39n ± 1%    30.21n ±  2%  +41.27% (p=0.000 n=10)
Cache2Miss/10000000-12                 34.72n ± 2%    47.26n ±  4%  +36.14% (p=0.000 n=10)
geomean                                10.50n         11.00n         +4.77%
```

## Contributions

Special thanks to Claude AI for its assistance throughout this project.
