import pandas as pd

# Read the CSV file
df = pd.read_csv('results.csv')

# Group by operation and calculate the median time and amount
result = df.groupby('operation').agg({'amount': 'first', 'time': 'median'})

# Print the median time and amount per operation
print(result)
