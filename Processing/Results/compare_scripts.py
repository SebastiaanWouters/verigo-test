import pandas as pd
import matplotlib.pyplot as plt

plt.style.use("ggplot")

COLOR1 = '#333333'
COLOR2 = '#999999'

# Read the CSV file
df = pd.read_csv('results_scripts.csv')

# Calculate median time for attacker and useful scripts
median_attacker = df[df['operation'] == 'attacker']['time'].median()
median_useful = df[df['operation'] == 'useful']['time'].median()

# Create a bar plot
categories = ['Attacker', 'Useful']
median_times = [median_attacker, median_useful]
colors = [COLOR1, COLOR2]  # Less popping, more subtle colors

plt.bar(categories, median_times, color=colors)
plt.xlabel('Script Category')
plt.ylabel('Median Time (ms)')
plt.title('Median Time for Attacker vs Useful Script')

# Save the figure
plt.savefig('script_comparison.png')
