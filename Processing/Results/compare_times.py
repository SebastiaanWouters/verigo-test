import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

plt.style.use("ggplot")

COLOR1 = '#333333'
COLOR2 = '#888888'
COLOR3 = '#BBBBBB'

# Read CSV files
df1 = pd.read_csv('results_middle.csv')
df2 = pd.read_csv('results.csv')
df3 = pd.read_csv('results_normal.csv')

# Get the common operations from all three dataframes
common_operations = set(df1['operation']).intersection(df2['operation'], df3['operation'])

# Filter the dataframes to include only the common operations
df1 = df1[df1['operation'].isin(common_operations)]
df2 = df2[df2['operation'].isin(common_operations)]
df3 = df3[df3['operation'].isin(common_operations)]


# Group by operation and calculate median time
df1 = df1.groupby('operation')['time'].median()
df2 = df2.groupby('operation')['time'].median()
df3 = df3.groupby('operation')['time'].median()


# Get list of operations
operations = df1.index.tolist()

# Get list of median times
median_time1 = df1.tolist()
median_time2 = df2.tolist()
median_time3 = (df3/10).tolist()



# Create figure and axis
fig, ax = plt.subplots()

# Define bar width
bar_width = 0.25

# Create index for operations
index = np.arange(len(operations))

# Create bars with different colors
bar1 = ax.bar(index, median_time1, bar_width, color=COLOR1, label='Verigo (Without Channels)')
bar2 = ax.bar(index+bar_width, median_time2, bar_width, color=COLOR2, label='Verigo (With Channels)')
bar3 = ax.bar(index+bar_width*2, median_time3, bar_width, color=COLOR3, label='Native Go')

# Set y-axis scale to logarithmic
ax.set_yscale('log')

# Add labels, title, and legend
ax.set_xlabel('Operation')
ax.set_ylabel('Median Time (Log Scale)')
ax.set_title('Median Time per Operation')
ax.set_xticks(index + 1.5*bar_width)
ax.set_xticklabels(operations, rotation=45, ha='right')

# Adjust bottom margin
plt.subplots_adjust(bottom=0.2)

plt.ylim(100,10**7+110000000)  # Adjust the limits as per your requirement
# Position the legend outside the plot
ax.legend()

# Save the plot
plt.tight_layout()
plt.savefig("time_comparison.png", bbox_inches='tight')