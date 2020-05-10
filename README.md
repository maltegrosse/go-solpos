Go-Solar Position and Intensity (SOLPOS)
=======================================
[![Go Report Card](https://goreportcard.com/badge/github.com/maltegrosse/go-solpos)](https://goreportcard.com/report/github.com/maltegrosse/go-solpos)
[![GoDoc](https://godoc.org/github.com/maltegrosse/go-solpos?status.svg)](https://pkg.go.dev/github.com/maltegrosse/go-solpos)
![Go](https://github.com/maltegrosse/go-solpos/workflows/Go/badge.svg) 

NREL's Solar Position and Intensity (SOLPOS 2.0) C function (adopted in go lang) calculates the apparent solar position and intensity (theoretical maximum solar energy) based on date, time, and location on Earth for the years *1950–2050*.
## Installation

This packages requires Go 1.13. If you installed it and set up your GOPATH, just run:

`go get -u github.com/maltegrosse/go-solpos`

## Usage

You can find some examples in the [examples](examples) directory.

Please visit https://www.nrel.gov/grid/solar-resource/solpos.html for additional information.

Some additional helper functions have been added to the original application logic.
## Notes
Note that your final decimal place values may vary based on your computer's floating-point storage and your compiler's mathematical algorithms.  If you agree with NREL's values for at least 5 significant digits, assume it works.

If the difference is smaller than 0.000001,it is shown as 0.

|       | NREL   | SOLPOS GO    |  Difference |
|---------------|-------|-------|-------|
| day     | 22  | 22  | 0 |
| month          | 7  | 7  | 0 |
| year      | 1999  | 1999  | 0  |
| daynum           | 203  | 203  | 0 |
| ampress |  1.326522  | 1.3265254683395031  | 0 |
| cosinc         | 0.912569  | 0.9125697491698622  | 0 |
| prime        | 1.037040  | 1.0370400228238181  | 0 |
| elevref          | 48.409931 | 48.40974986370258  | 0.00018113629742089188 |
| etrtilt     | 1207.547363 | 1207.5486465939166 | 0 |
| sbcf     | 1.201910  | 1.2019108852729568 | 0 |
| sunrise      | 347.173431  | 347.1746053157332  | 0 |
| sunset         | 1181.111206  | 1181.1101930505447 | 0.0010129494553439145 |
| unprime           | 0.964283  | 0.9642829379690094 | 0 |
| zenref        | 41.590069  | 41.59025013629742  | 0 |
| amass           | 1.335752  | 1.3357557648388834  | 0 |
| azim          | 97.032875  | 97.03331438627943  | 0 |
| etr          | 989.668518  | 989.6657077767151  | 0.0028102232848823405 |
| etrn          | 1323.239868  | 1323.2398374944908  | 3.0505509130307473e-05 |



## License
**[ NREL data disclaimer](https://www.nrel.gov/disclaimer.html)**

Adoption in Golang under **[MIT license](http://opensource.org/licenses/mit-license.php)** 2020 © Malte Grosse.

