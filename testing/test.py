def float_range(start, stop, step):
    result = []
    value = start
    while value < stop:
        result.append(value)
        value += step
    return result

# Example usage:
arr = float_range(0.0, 1.0, 0.02)
print(arr)  # [0.0, 0.2, 0.4, 0.6, 0.8]

import numpy as np

arr = np.arange(0.0, 1.0, 0.02) 
print(arr)
