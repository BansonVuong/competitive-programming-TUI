import pathlib
import re
import time
from typing import List, Tuple

import requests
from bs4 import BeautifulSoup, Tag


class RateLimitedSession:
    def __init__(self, min_interval_seconds: float = 1.0) -> None:
        self.session = requests.Session()
        self.session.headers.update(
            {
                "User-Agent": (
                    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
                    "AppleWebKit/537.36 (KHTML, like Gecko) "
                    "Chrome/122.0.0.0 Safari/537.36"
                )
            }
        )
        self.min_interval_seconds = min_interval_seconds
        self._last_request_time = 0.0

    def get(self, url: str, timeout: int = 20) -> requests.Response:
        elapsed = time.time() - self._last_request_time
        if elapsed < self.min_interval_seconds:
            time.sleep(self.min_interval_seconds - elapsed)
        resp = self.session.get(url, timeout=timeout)
        self._last_request_time = time.time()
        return resp


def extract_problem_codes(statements_dir: pathlib.Path) -> List[str]:
    codes: List[str] = []
    for cpp_file in sorted(statements_dir.glob("*.cpp")):
        code = cpp_file.stem.strip()
        lower = code.lower()
        if lower.startswith("codeforces") or lower.startswith("cses"):
            continue
        codes.append(code)
    return codes


def _clean_text(text: str) -> str:
    text = text.replace("\r\n", "\n").replace("\r", "\n")
    return text.strip("\n")


def extract_samples_from_headers(soup: BeautifulSoup) -> Tuple[List[str], List[str]]:
    inputs: List[str] = []
    outputs: List[str] = []
    header_tags = ("h1", "h2", "h3", "h4", "h5", "h6", "strong", "b")

    for tag in soup.find_all(header_tags):
        if not isinstance(tag, Tag):
            continue
        label = tag.get_text(" ", strip=True).lower()
        if label.startswith("sample input"):
            pre = tag.find_next("pre")
            if pre:
                inputs.append(_clean_text(pre.get_text("\n")))
        elif label.startswith("sample output"):
            pre = tag.find_next("pre")
            if pre:
                outputs.append(_clean_text(pre.get_text("\n")))
    return inputs, outputs


def extract_samples_fallback(soup: BeautifulSoup) -> Tuple[List[str], List[str]]:
    pres = [_clean_text(pre.get_text("\n")) for pre in soup.find_all("pre")]
    inputs: List[str] = []
    outputs: List[str] = []
    for idx, text in enumerate(pres):
        if idx % 2 == 0:
            inputs.append(text)
        else:
            outputs.append(text)
    return inputs, outputs


def extract_sample_pairs(html: str) -> List[Tuple[str, str]]:
    soup = BeautifulSoup(html, "html.parser")
    inputs, outputs = extract_samples_from_headers(soup)

    if not inputs or not outputs:
        inputs, outputs = extract_samples_fallback(soup)

    count = min(len(inputs), len(outputs))
    return [(inputs[i], outputs[i]) for i in range(count)]


def save_sample_pairs(output_dir: pathlib.Path, code: str, pairs: List[Tuple[str, str]]) -> None:
    for idx, (sample_in, sample_out) in enumerate(pairs, start=1):
        in_path = output_dir / f"{code}.{idx}.in"
        out_path = output_dir / f"{code}.{idx}.out"
        in_path.write_text(sample_in + "\n", encoding="utf-8", newline="\n")
        out_path.write_text(sample_out + "\n", encoding="utf-8", newline="\n")


def main() -> None:
    statements_dir = pathlib.Path("./problems")
    output_dir = pathlib.Path("./problems")
    output_dir.mkdir(parents=True, exist_ok=True)

    codes = extract_problem_codes(statements_dir)
    if not codes:
        print("No eligible statement files found.")
        return

    client = RateLimitedSession(min_interval_seconds=1.0)

    for code in codes:
        if re.match(r"^(codeforces|cses)", code, flags=re.IGNORECASE):
            continue

        url = f"https://dmoj.ca/problem/{code}"
        print(f"Fetching {url}")
        try:
            resp = client.get(url)
            if resp.status_code != 200:
                print(f"  skipped ({resp.status_code})")
                continue

            pairs = extract_sample_pairs(resp.text)
            if not pairs:
                print("  no sample pairs found")
                continue

            save_sample_pairs(output_dir, code, pairs)
            print(f"  wrote {len(pairs)} sample pair(s)")
        except requests.RequestException as exc:
            print(f"  request failed: {exc}")
        except Exception as exc:
            print(f"  parse/save failed: {exc}")


if __name__ == "__main__":
    main()
