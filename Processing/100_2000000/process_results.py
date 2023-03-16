from matplotlib import pyplot as plt
import numpy as np
from matplotlib import style
import pandas as pd
import sys
  
plt.style.use("ggplot")


df = pd.read_csv("results.csv")
df = df.groupby('operation').median()
#df = df.sort_values(by=['time'])
df['weights'] = df['time'].map(lambda time: round(time / df['time'].min(), 4))

df.to_csv("weighted_results.csv")

df[["time"]].plot.bar(legend=None)
plt.title("Operation Benchmark")
plt.xlabel("Operation")
plt.ylabel("Time Per " + str(int((df['amount'].mean()))) + " Operations (ms)")

plt.tight_layout()
plt.savefig("benchmark.png")

df[["weights"]].plot.bar(legend=None)
plt.title("Operation Weights")
plt.xlabel("Operation")
plt.ylabel("Weight")

plt.tight_layout()
plt.savefig("weights.png")
