import re

with open("full_blog.html") as f:
    text = f.read()

# I will find all the code snippets, to map them correctly.
from bs4 import BeautifulSoup
soup = BeautifulSoup(text, 'html.parser')

snippets = []
for code in soup.find_all("code"):
    if "language-yaml" in code.get("class", []) or "language-bash" in code.get("class", []) or "language-toml" in code.get("class", []):
        snippets.append(code.get_text(separator="\n"))

for i, s in enumerate(snippets):
    print(f"--- SNIPPET {i} ---")
    lines = s.split('\n')
    clean_lines = []
    for l in lines:
        l = l.strip()
        if len(l) > 0 and l[0].isdigit():
            # remove line number
            parts = l.split(' ', 1)
            if len(parts) > 1:
                clean_lines.append(parts[1])
            else:
                clean_lines.append(l)
        else:
            clean_lines.append(l)

    print("\n".join(clean_lines[:5]))
