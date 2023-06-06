from matplotlib import pyplot as plt
import numpy as np
from matplotlib import style
import pandas as pd
import sys
  
plt.style.use("ggplot")

##Process results from benchmark (verigo)
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

df[["weights"]].plot.bar(legend=None, color='#333333')
plt.title("Operation Weights")
plt.xlabel("Operation")
plt.ylabel("Weight")

plt.tight_layout()
plt.savefig("weights.png")

##Process results from middle benchmark (verigo without channel)
df = pd.read_csv("results_middle.csv")
df = df.groupby('operation').median()
#df = df.sort_values(by=['time'])
df['weights'] = df['time'].map(lambda time: round(time / df['time'].min(), 4))

df.to_csv("weighted_results_middle.csv")

df[["time"]].plot.bar(legend=None)
plt.title("Operation Benchmark")
plt.xlabel("Operation")
plt.ylabel("Time Per " + str(int((df['amount'].mean()))) + " Operations (ms)")

plt.tight_layout()
plt.savefig("benchmark_middle.png")

df[["weights"]].plot.bar(legend=None, color='#FF8888')
plt.title("Operation Weights")
plt.xlabel("Operation")
plt.ylabel("Weight")

plt.tight_layout()
plt.savefig("weights_middle.png")


##Process results from simple benchmark (verigo without counter)
df = pd.read_csv("results_simple.csv")
df = df.groupby('operation').median()
#df = df.sort_values(by=['time'])
df['weights'] = df['time'].map(lambda time: round(time / df['time'].min(), 4))

df.to_csv("weighted_results_simple.csv")

df[["time"]].plot.bar(legend=None)
plt.title("Operation Benchmark")
plt.xlabel("Operation")
plt.ylabel("Time Per " + str(int((df['amount'].mean()))) + " Operations (ms)")

plt.tight_layout()
plt.savefig("benchmark_simple.png")

df[["weights"]].plot.bar(legend=None)
plt.title("Operation Weights")
plt.xlabel("Operation")
plt.ylabel("Weight")

plt.tight_layout()
plt.savefig("weights_simple.png")

##Process results from normal benchmark (basic arithmetic)
df = pd.read_csv("results_normal.csv")
df = df.groupby('operation').median()
#df = df.sort_values(by=['time'])
df['weights'] = df['time'].map(lambda time: round(time / df['time'].min(), 4))

df.to_csv("weighted_results_normal.csv")

df[["time"]].plot.bar(legend=None)
plt.title("Operation Benchmark")
plt.xlabel("Operation")
plt.ylabel("Time Per " + str(int((df['amount'].mean()))) + " Operations (ms)")

plt.tight_layout()
plt.savefig("benchmark_normal.png")

df[["weights"]].plot.bar(legend=None)
plt.title("Operation Weights")
plt.xlabel("Operation")
plt.ylabel("Weight")

plt.tight_layout()
plt.savefig("weights_normal.png")
