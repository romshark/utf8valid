# Go `unicode/utf8` Optimization Opportunity

This package demonstrates a standard library optimization opportunity in:

- [`utf8.Valid`](https://pkg.go.dev/unicode/utf8#Valid)
- [`utf8.ValidString`](https://pkg.go.dev/unicode/utf8#ValidString)

This is a copy of the original implementation dating
[Sep 15, 2024, 4:05 AM GMT+2](https://github.com/golang/go/blob/3d33437c450aa74014ea1d41cd986b6ee6266984/src/unicode/utf8/utf8.go):
The changes that seem to lead to a significant performance improvement of up to 41%
are minimal and located at lines 481 and 533.

The geomean for M1 ARM64 processors is around `-12.6%`.

## Benchmark Results

Execute `./bench.sh Fn 10` to run the benchmarks with 10 iterations.

### Apple M1 Max

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/utf8valid
                                          │ .std_Fn.txt │             .opt_Fn.txt             │
                                          │   sec/op    │   sec/op     vs base                │
FnValidString/empty-10                      2.174n ± 0%   2.129n ± 1%   -2.07% (p=0.000 n=10)
FnValidString/single_byte-10                2.484n ± 0%   2.174n ± 0%  -12.50% (p=0.000 n=10)
FnValidString/single_utf8_rune-10           3.416n ± 0%   3.131n ± 0%   -8.36% (p=0.000 n=10)
FnValidString/short_ascii-10                4.038n ± 0%   3.727n ± 0%   -7.69% (p=0.000 n=10)
FnValidString/short_utf8-10                 30.26n ± 0%   23.31n ± 0%  -22.98% (p=0.000 n=10)
FnValidString/all_utf8_ukranian_poetry-10   5.532µ ± 0%   3.255µ ± 0%  -41.16% (p=0.000 n=10)
FnValidString/long_ascii-10                 294.0n ± 0%   294.5n ± 0%   +0.17% (p=0.001 n=10)
FnValidString/wikipedia_diacritic_html-10   433.0µ ± 0%   431.1µ ± 0%   -0.44% (p=0.000 n=10)
FnValidString/wikipedia_japan_html-10       1.594m ± 0%   1.504m ± 0%   -5.66% (p=0.000 n=10)
FnValidString/invalid_surrogate_max-10      2.794n ± 0%   2.485n ± 0%  -11.08% (p=0.000 n=10)
FnValid/empty-10                            2.174n ± 0%   2.035n ± 0%   -6.39% (p=0.000 n=10)
FnValid/single_byte-10                      2.484n ± 0%   2.174n ± 0%  -12.50% (p=0.000 n=10)
FnValid/single_utf8_rune-10                 3.416n ± 0%   3.138n ± 0%   -8.14% (p=0.000 n=10)
FnValid/short_ascii-10                      4.036n ± 0%   3.727n ± 0%   -7.68% (p=0.000 n=10)
FnValid/short_utf8-10                       30.03n ± 0%   23.30n ± 0%  -22.41% (p=0.000 n=10)
FnValid/all_utf8_ukranian_poetry-10         5.532µ ± 0%   3.257µ ± 0%  -41.13% (p=0.000 n=10)
FnValid/long_ascii-10                       296.5n ± 0%   298.1n ± 0%   +0.54% (p=0.000 n=10)
FnValid/wikipedia_diacritic_html-10         433.1µ ± 0%   431.8µ ± 0%   -0.30% (p=0.000 n=10)
FnValid/wikipedia_japan_html-10             1.594m ± 0%   1.502m ± 0%   -5.81% (p=0.000 n=10)
FnValid/invalid_surrogate_max-10            2.795n ± 0%   2.485n ± 0%  -11.09% (p=0.000 n=10)
geomean                                     153.1n        134.3n       -12.25%

                                          │ .std_Fn.txt  │             .opt_Fn.txt             │
                                          │     B/op     │    B/op     vs base                 │
FnValidString/empty-10                      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/single_byte-10                0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/single_utf8_rune-10           0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/short_ascii-10                0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/short_utf8-10                 0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/all_utf8_ukranian_poetry-10   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/long_ascii-10                 0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/wikipedia_diacritic_html-10   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/wikipedia_japan_html-10       0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/invalid_surrogate_max-10      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/empty-10                            0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/single_byte-10                      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/single_utf8_rune-10                 0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/short_ascii-10                      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/short_utf8-10                       0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/all_utf8_ukranian_poetry-10         0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/long_ascii-10                       0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/wikipedia_diacritic_html-10         0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/wikipedia_japan_html-10             0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/invalid_surrogate_max-10            0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
geomean                                                ²               +0.00%                ²
¹ all samples are equal
² summaries must be >0 to compute geomean

                                          │ .std_Fn.txt  │             .opt_Fn.txt             │
                                          │  allocs/op   │ allocs/op   vs base                 │
FnValidString/empty-10                      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/single_byte-10                0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/single_utf8_rune-10           0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/short_ascii-10                0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/short_utf8-10                 0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/all_utf8_ukranian_poetry-10   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/long_ascii-10                 0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/wikipedia_diacritic_html-10   0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/wikipedia_japan_html-10       0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValidString/invalid_surrogate_max-10      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/empty-10                            0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/single_byte-10                      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/single_utf8_rune-10                 0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/short_ascii-10                      0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/short_utf8-10                       0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/all_utf8_ukranian_poetry-10         0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/long_ascii-10                       0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/wikipedia_diacritic_html-10         0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/wikipedia_japan_html-10             0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
FnValid/invalid_surrogate_max-10            0.000 ± 0%     0.000 ± 0%       ~ (p=1.000 n=10) ¹
geomean                                                ²               +0.00%                ²
¹ all samples are equal
² summaries must be >0 to compute geomean
```
