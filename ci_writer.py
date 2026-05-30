import sys
from bs4 import BeautifulSoup

def clean_code(html):
    soup = BeautifulSoup(html, "html.parser")
    for s in soup.find_all("span", class_="ln"):
        s.extract()
    return soup.get_text()

with open("full_blog.html") as f:
    html = f.read()

soup = BeautifulSoup(html, "html.parser")
codes = soup.find_all("code")

def get_code(index):
    return clean_code(str(codes[index])).strip()

print(get_code(2)) # name: CI/CD
print("===========")
print(get_code(3)) # prepare-release-tag
print("===========")
print(get_code(4)) # concurrency
print("===========")
print(get_code(6)) # discover
print("===========")
print(get_code(15)) # golangci
print("===========")
print(get_code(22)) # autofix
print("===========")
print(get_code(27)) # goreleaser
print("===========")
print(get_code(38)) # publish-draft
print("===========")
