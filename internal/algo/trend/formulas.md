* Opt 1

package main

import (
    "fmt"
    "gonum.org/v1/gonum/floats"
    "gonum.org/v1/gonum/mat"
    "gonum.org/v1/gonum/stat"
)

// polyfit performs a polynomial fit of the given degree and returns the coefficients.
func polyfit(x, y []float64, degree int) []float64 {
    n := len(x)
    if n != len(y) {
        panic("x and y slices must have the same length")
    }

    // Create the Vandermonde matrix
    vandermonde := mat.NewDense(n, degree+1, nil)
    for i := 0; i < n; i++ {
        for j := 0; j <= degree; j++ {
            vandermonde.Set(i, j, floats.Pow(x[i], float64(j)))
        }
    }

    // Perform least squares fitting
    var coeff mat.VecDense
    stat.Regression(&coeff, vandermonde, y, nil)

    return coeff.RawVector().Data
}

func main() {
    close := []float64{100.0, 102.5, 101.0, 105.0} // Example data points
    x := make([]float64, len(close))
    for i := range x {
        x[i] = float64(i)
    }

    degree := 1 // For a linear fit
    coefs := polyfit(x, close, degree)

    fmt.Printf("Coefficients: %v\n", coefs)
}
 * opt 2

package main

import (
"fmt"
)

// linearRegression calculates the slope and intercept for a simple linear fit.
func linearRegression(x, y []float64) (float64, float64) {
if len(x) != len(y) {
panic("x and y slices must have the same length")
}

    n := float64(len(x))
    var sumX, sumY, sumXY, sumX2 float64

    for i := 0; i < len(x); i++ {
        sumX += x[i]
        sumY += y[i]
        sumXY += x[i] * y[i]
        sumX2 += x[i] * x[i]
    }

    // Calculating slope (m) and intercept (b)
    slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
    intercept := (sumY - slope*sumX) / n

    return slope, intercept
}

func main() {
close := []float64{100.0, 102.5, 101.0, 105.0} // Example close values
x := make([]float64, len(close))
for i := range x {
x[i] = float64(i)
}

    // Perform linear regression (1st-degree polynomial fit)
    slope, intercept := linearRegression(x, close)

    fmt.Printf("Slope: %f, Intercept: %f\n", slope, intercept)
}

* opt 3






