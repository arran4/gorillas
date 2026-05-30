import sys

with open("extracted_full_blocks.txt") as f:
    text = f.read()

blocks = text.split("================================= BLOCK")
blocks = [b.split("=================================\n", 1)[1] if "=================================\n" in b else "" for b in blocks]
blocks = [b.strip() for b in blocks if b.strip()]

def get_block(idx):
    if idx < len(blocks):
        return blocks[idx]
    return ""

print(f"Total blocks extracted: {len(blocks)}")

def print_block(i):
    print(f"\n--- BLOCK {i} ---")
    print("\n".join(blocks[i].split("\n")[:10]) + "...")

# The skeleton is block 41.
# I need to find the blocks that define jobs:
# route: block 3
# discover: block 6
# gitleaks: block 11
# golangci: block 15
# go-test: block 15
# go-vet: block 15
# go-fmt-pr: block 15
# autofix: block 23
# cleanup-autofix-prs: block 24
# goreleaser: block 28 (wait, goreleaser is block 27) Let's inspect them!

for i, b in enumerate(blocks):
    if "name: CI/CD" in b:
        print(f"BLOCK {i}: name: CI/CD")
    if "  route:" in b:
        print(f"BLOCK {i}: route")
    if "  discover:" in b:
        print(f"BLOCK {i}: discover")
    if "  gitleaks:" in b:
        print(f"BLOCK {i}: gitleaks")
    if "  golangci:" in b:
        print(f"BLOCK {i}: golangci")
    if "  autofix:" in b:
        print(f"BLOCK {i}: autofix")
    if "  goreleaser:" in b:
        print(f"BLOCK {i}: goreleaser")
    if "  publish-draft:" in b:
        print(f"BLOCK {i}: publish-draft")

print("----")
print(blocks[41][:500])
