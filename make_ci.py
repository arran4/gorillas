import sys

with open("extracted_full_blocks.txt") as f:
    text = f.read()

blocks = text.split("================================= BLOCK")
blocks = [b.split("=================================\n", 1)[1] if "=================================\n" in b else "" for b in blocks]
blocks = [b.strip() for b in blocks if b.strip()]

def get(idx): return blocks[idx]

# We need to construct the CI file.
# The template is block 41. But it has # ...
# So we'll replace the # ... with the actual implementations.
# For languages we don't have, we can probably skip them or put them as # ... but the blog post says:
# "The target outcome: One workflow file handles... supports mixed repos..."
# "Project-type decisions should mostly be install/template-time (human comments + toggles), with lightweight runtime detection as a safety net."
# Actually, the user asked: "Add a single CI/CD workflow at .github/workflows/ci.yml combining everything required from the blog post."
# I should put everything in to make it a fully featured workflow as the post outlines.

template = blocks[41]

# Let's inspect block 41
import re
print("Jobs to replace in block 41:")
for match in re.findall(r"([a-z0-9_-]+):\n\s+needs:.*?\n\s+# \.\.\.", template, re.DOTALL):
    print(match)
