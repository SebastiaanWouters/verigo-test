import pandas as pd

# Read the CSV file
df = pd.read_csv('results.csv')

# Group by operation and calculate the median time
median_time_per_operation = df.groupby('operation')['time'].median()

# Print the median time per operation
print(median_time_per_operation)