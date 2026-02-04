# Example 1:

## go run . examples/example00.txt 
- Number of ants: 4
- time elapsed: 64.101µs
 
## go run . examples/example00.txt 
- Number of ants: 40
- time elapsed: 183.264µs

## go run . examples/example00.txt 
- Number of ants: 400
- time elapsed: 2242.369µs

## go run . examples/example00.txt 
- Number of ants: 4000
- time elapsed: 114699.803µs

### Analysis:
Given the analysis on the time complexity, the time complexity at first seems like it starts with logarithmic time, then slowly transitions into linear time, and then, finally, to quadratic time towards the end. This implies that there is some an efficient-cache application in the beginning, but soon becomes bottle-necked somewhere.

# Example 2:
## go run . examples/example01.txt 
- Number of ants: 10
- time elapsed: 95.06µs

## go run . examples/example01.txt 
- Number of ants: 100
- time elapsed: 345.011µs

## go run . examples/example01.txt 
- Number of ants: 1000
- time elapsed: 9.38227ms

## go run . examples/example01.txt 
- Number of ants: 10000
- time elapsed: 343.827609ms

### Analysis:
As we can see, the same trend holds true here as well. The time complexity seemed logarithmic at first as it went from 95µs to 365µs instead of ~950µs. However that trend doesn't hold true for larger values. Meaning some bottle-neck is for sure taking place. 