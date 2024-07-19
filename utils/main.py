import pandas as pd

# Load the TSV file, specifying how to handle bad lines
df = pd.read_csv("data.tsv", sep="\t", on_bad_lines="skip")

# Save the cleaned file if needed
df.to_csv("cleaned_data.tsv", sep="\t", index=False)
