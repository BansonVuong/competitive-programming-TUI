from pathlib import Path
import argparse

def split_by_h1(md_text: str) -> dict[str, str]:
    sections = {}
    title = None
    lines = []

    for line in md_text.splitlines():
        if line.startswith("# "):
            if title is not None:
                sections[title] = "\n".join(lines).strip() + "\n"
            title = line[2:].strip()
            lines = [line]
        elif title is not None:
            lines.append(line)
    
    if title is not None:
        sections[title] = "\n".join(lines).strip() + "\n"
    
    return sections

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("input_file", help="Path to the input markdown file")
    args = parser.parse_args()
    
    src = Path(args.input_file)

    sections = split_by_h1(src.read_text())

    out_dir = Path("./problems").mkdir(parents=True, exist_ok=True)
    for title, content in sections.items():
        with open(f"./problems/{title}.md", "w") as f:
            f.write(content)

if __name__ == "__main__":
    main()