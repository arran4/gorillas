from bs4 import BeautifulSoup

with open("full_blog.html") as f:
    soup = BeautifulSoup(f, "html.parser")

out = []
for code in soup.find_all("code"):
    classes = code.get("class", [])
    if any(c in classes for c in ["language-yaml", "language-toml", "language-json", "language-bash"]):
        block_text = ""
        for line in code.find_all("span", class_="line"):
            cl = line.find("span", class_="cl")
            if cl:
                block_text += cl.get_text() + "\n"
        out.append(block_text)

for i, block in enumerate(out):
    print(f"================================= BLOCK {i} =================================")
    print(block.strip())
