import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

# Read CSV files
df1 = pd.read_csv('results_middle.csv')
df2 = pd.read_csv('results.csv')

# Group by operation and calculate median time
df1 = df1.groupby('operation')['time'].median()
df2 = df2.groupby('operation')['time'].median()

# Get list of operations
operations = df1.index.tolist()

# Get list of median times
median_time1 = df1.tolist()
median_time2 = df2.tolist()

# Create figure and axis
fig, ax = plt.subplots()

# Define bar width
bar_width = 0.35


# Create index for operations
index = np.arange(len(operations))

# Create bars with different colors
bar1 = ax.bar(index, median_time1, bar_width, color='b', label='unchanneled verigo')
bar2 = ax.bar(index+bar_width, median_time2, bar_width, color='r', label='channeled verigo')

# Add labels, title, and legend
ax.set_xlabel('Operations')
ax.set_ylabel('Median Time')
ax.set_title('Median Time per Operation')
ax.set_xticks(index + bar_width / 2)
ax.set_xticklabels(operations)
ax.legend()

# Show the plot
plt.show()