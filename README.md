Package intset implements set theory methods for sets in the ℤ
domain. The operations can handle sets in the range -∞:∞, but the
minimum and maximum values of the integers are limited by the
underlying 32-bit or a 64-bit machine platform.

Installation
------------

The recommended way to install intset

```
go get github.com/stianwa/intset
```

Examples
--------

```go

package main
 
import (
       "github.com/stianwa/intset"
       "fmt"
)

func main() {
       a := intset.New(intset.Range(-300, -30), intset.NegInf(-500), intset.PosInf(500))
       b := intset.New(intset.Range(-47, 23))

       fmt.Printf("%s union %s = %s\n", a, b, a.Union(b))
       fmt.Printf("%s intersect %s = %s\n", a, b, a.Intersect(b))
       fmt.Printf("complement of %s = %s\n", a, a.Complement())

       if inf, c := b.Cardinality(); !inf {
            fmt.Printf("cardinality of %s is: %d\n", b, c)
       } else {
            fmt.Println("cardinality of %s is infinite")
       }
}
```

State
-------
The intset module is currently under development. Do not use for production.


License
-------

GPLv3, see [LICENSE.md](LICENSE.md)
