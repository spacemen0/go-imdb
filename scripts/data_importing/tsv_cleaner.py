import pandas as pd

def filter_bad_lines(input_file, output_file):
    try:
        # Read the TSV file and ignore bad lines automatically
        df = pd.read_csv(input_file, sep='\t', dtype=str, on_bad_lines='skip')


        # Additional criteria can be added here as needed
        # Example: Removing rows where a specific column is not a valid integer
        # filtered_df = filtered_df[filtered_df['some_column'].apply(lambda x: x.isdigit())]

        # Write the filtered data to a new TSV file
        df.to_csv(output_file, sep='\t', index=False)

        print(f"Filtered file saved to {output_file}")

        print(f"Error processing the file: {e}")

# Usage example
input_file = 'title.basics.tsv'
output_file = 'title.tsv'
filter_bad_lines(input_file, output_file)
