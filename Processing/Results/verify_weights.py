from matplotlib import pyplot as plt
import numpy as np
from matplotlib import style
import pandas as pd
import sys


  
plt.style.use("ggplot")

COLOR1 = '#444444'
COLOR2 = '#888888'


df_weights = pd.read_csv("weighted_results.csv", index_col="operation")['weights']
df = pd.read_csv("results.csv")
df = df.groupby('operation').mean()
df['weights'] = df_weights


df['simple_score'] = df['amount']
df['advanced_score'] = df['amount']*df['weights']
df['simple_score_per_ms'] = df['simple_score'] / df['time']
df['advanced_score_per_ms'] = df['advanced_score'] / df['time']

df.to_csv("verified_weighted_results.csv")

df[["simple_score_per_ms"]].plot.bar(legend=None, color=COLOR1)
print(df[['simple_score_per_ms']])
plt.ylim(df['simple_score_per_ms'].min() - df['simple_score_per_ms'].min()*1/20,df['simple_score_per_ms'].max() + df['simple_score_per_ms'].max()*1/20)
#plt.ylim(0.5,2)
plt.title("Unweighted Operation Score Per Millisecond")
plt.xlabel("Operation")
plt.ylabel("Score")
plt.tight_layout()
plt.savefig("scores_simple_time.png")

df[["simple_score"]].plot.bar(legend=None, color=COLOR1)
plt.ylim(df['simple_score'].min() - df['simple_score'].min()*1/20,df['simple_score'].max() + df['simple_score'].max()*1/20)
plt.title("Unweighted Operation Score For " + str(int((df['amount'].mean()))) + " Operations")
plt.xlabel("Operation")
plt.ylabel("Score")
plt.tight_layout()
plt.savefig("scores_simple.png")


df[["advanced_score_per_ms"]].plot.bar(legend=None, color=COLOR2)
plt.ylim(df['advanced_score_per_ms'].min() - df['advanced_score_per_ms'].min()*1/20,df['advanced_score_per_ms'].max() + df['advanced_score_per_ms'].max()*1/20)
#plt.ylim(0.5,2)
plt.title("Weighted Operation Score Per Millisecond")
plt.xlabel("Operation")
plt.ylabel("Score")
plt.tight_layout()
plt.savefig("scores_advanced_time.png")

df[["advanced_score"]].plot.bar(legend=None, color=COLOR2)
plt.ylim(df['advanced_score'].min() - df['advanced_score'].min()*1/20,df['advanced_score'].max() + df['advanced_score'].max()*1/20)
plt.title("Weighted Operation Score For " + str(int((df['amount'].mean()))) + " Operations")
plt.xlabel("Operation")
plt.ylabel("Score")
plt.tight_layout()
plt.savefig("scores_advanced.png")

df[['advanced_score_per_ms']].plot(kind='box', title='Scores')
plt.ylim(0.7,1.7)
df[['simple_score_per_ms']].plot(kind='box', title='Scores')
plt.ylim(0.7,1.7)
plt.show()
