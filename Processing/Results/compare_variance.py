import matplotlib.pyplot as plt
import numpy as np

# Data for the first violin
data1 = [4.757218, 3.683231, 4.698836, 5.137418, 2.370159, 4.862087, 3.408404, 3.809639, 4.796626, 5.142529, 2.803204, 1.554430, 3.497523, 4.967762, 3.893909, 4.772597, 3.436606]

# Data for the second violin
data2 = [5.118291, 5.115271, 5.069104, 5.137418, 5.147749, 5.107137, 5.128966, 5.143394, 5.139105, 5.156414, 5.052215, 5.157911, 5.123522, 5.144117, 5.141128, 5.144383, 5.131541]

# Create a figure and axes
fig, ax1 = plt.subplots(figsize=(10, 5))
ax2 = ax1.twinx()

# Create the violin plots for both sets of data
violin_parts1 = ax1.violinplot(data1, positions=[0], showmedians=True)
violin_parts2 = ax2.violinplot(data2, positions=[1], showmedians=True)

# Set the color for the first violin plot
violin_parts1['bodies'][0].set(facecolor='#888888')
violin_parts1['cmedians'].set(color='#888888')

# Set the color for the second violin plot
violin_parts2['bodies'][0].set(facecolor='#888888')
violin_parts2['cmedians'].set(color='#888888')

# Set labels and title
ax1.set_ylabel('Values', color='#000000')
ax1.set_title('Variation Between Scores (Weighted vs. Unweighted))')

# Set x-axis tick labels
ax1.set_xticks([0, 1])
ax1.set_xticklabels(['Weighted', 'Unweighted'])

# Calculate margins for the y-axis limits
margin1 = 0.05 * (max(data1) - min(data1))
margin2 = 0.05 * (max(data2) - min(data2))

# Adjust y-axis limits based on the data with margins
ax1.set_ylim([min(data1) - margin1, max(data1) + margin1])
ax2.set_ylim([min(data2) - margin2, max(data2) + margin2])

# Show the plot
plt.savefig('violin_plot.png')
