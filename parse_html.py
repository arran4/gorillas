from bs4 import BeautifulSoup

with open("full_blog.html") as f:
    soup = BeautifulSoup(f, "html.parser")

for i, code in enumerate(soup.find_all("code")):
    classes = code.get("class", [])
    if "language-yaml" in classes or "language-toml" in classes or "language-json" in classes or "language-bash" in classes:
        print(f"--- BLOCK {i} ({classes}) ---")
        lines = code.find_all("span", class_="line")
        for line in lines:
            cl = line.find("span", class_="cl")
            if cl:
                print(cl.get_text())
