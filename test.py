# Example snippet of your data
data = {
    "lines": [
        {
            "text": "If you visited the RVPR office today,",
            "voice": None
        },
        {
            "text": "you would have seen a lot of activity.",
            "voice": None
        },
        {
            "text": "People were busy preparing for the upcoming event.",
            "voice": None
        }
    ]
}

# Extracting the "text" values
lines_text = [line["text"] for line in data["lines"]]

# Printing the extracted text
for text in lines_text:
    print(text)