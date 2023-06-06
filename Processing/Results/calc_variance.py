import pandas as pd
import matplotlib.pyplot as plt

plt.style.use("ggplot")

COLOR1 = '#333333'
COLOR2 = '#999999'

# Read the CSV file
df = pd.read_csv('results.csv')

# Define the four operations
operations = ['addition', 'substraction', 'multiplication', 'division']

# Create a figure with a two by two grid
fig, axes = plt.subplots(2, 2, figsize=(12, 8))

# Iterate over operations and plot boxplots
for i, operation in enumerate(operations):
    # Filter data for the current operation
    operation_data = df[df['operation'] == operation]['time']

    # Calculate subplot position in the grid
    row = i // 2
    col = i % 2

    # Create a boxplot for the current operation
    boxplot = axes[row, col].boxplot(operation_data, patch_artist=True, labels=[operation], flierprops=dict(marker='D', markerfacecolor=COLOR1, markersize=8))

    # Change the color of the boxplot
    boxplot['boxes'][0].set(facecolor=COLOR2)

    axes[row, col].set_ylabel('Time (ms)')
    axes[row, col].set_title(f'Execution Times for {operation.capitalize()}')

# Remove empty subplot if the number of operations is odd
if len(operations) % 2 != 0:
    fig.delaxes(axes[1, 1])

# Adjust spacing between subplots
plt.tight_layout()

# Save the figure
plt.savefig('execution_times_boxplots.png')
