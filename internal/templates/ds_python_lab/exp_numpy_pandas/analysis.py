import pandas as pd
import matplotlib.pyplot as plt


df = pd.read_csv("dataset.csv")
print(df.describe())

df["value"].hist()
plt.savefig("output_plot.png")
