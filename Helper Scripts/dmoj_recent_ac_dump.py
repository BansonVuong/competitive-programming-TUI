#!/usr/bin/env python3
"""Download recent accepted DMOJ submissions for a user.

This script is intentionally explicit and modular so you can extend it.
"""

from __future__ import annotations

import argparse
import html
import json
import os
import re
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from dataclasses import dataclass
from typing import Any, Dict, Iterable, List, Optional, Tuple


LISTING_ENDPOINT_CANDIDATES = [
    "/api/v2/submissions",
    "/api/v2/submissions/{username}",
    "/api/submissions/{username}",
    "/submissions/{username}",
    "/submissions/user/{username}",
    "/submissions",
]

SUBMISSION_DETAIL_ENDPOINT_CANDIDATES = [
    "/api/v2/submission/{id}",
    "/api/submission/{id}",
    "/submission/{id}",
]

SOURCE_ENDPOINT_CANDIDATES = [
    "/api/v2/submission/{id}",
    "/api/submission/{id}",
    "/src/{id}",
    "/submission/{id}/source",
    "/submission/{id}/raw",
]

AC_MARKERS = {
    "ac",
    "accepted",
    "done",
    "correct",
}


@dataclass
class Submission:
    submission_id: int
    problem: str
    language: str
    verdict: str
    source: Optional[str]


class HttpClient:
    def __init__(
        self,
        base_url: str,
        api_token: Optional[str],
        sessionid: Optional[str],
        timeout: float,
        user_agent: str,
    ) -> None:
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout
        self.default_headers: Dict[str, str] = {
            "User-Agent": user_agent,
            "Accept": "application/json,text/html;q=0.9,*/*;q=0.8",
        }
        if api_token:
            self.default_headers["Authorization"] = f"Token {api_token}"
        cookie_parts = []
        if sessionid:
            cookie_parts.append(f"sessionid={sessionid}")
        if cookie_parts:
            self.default_headers["Cookie"] = "; ".join(cookie_parts)

    def make_url(self, path: str, query: Optional[Dict[str, Any]] = None) -> str:
        path = path if path.startswith("/") else f"/{path}"
        url = f"{self.base_url}{path}"
        if query:
            query = {k: v for k, v in query.items() if v is not None}
            if query:
                url = f"{url}?{urllib.parse.urlencode(query)}"
        return url

    def get(self, url: str) -> Tuple[int, str, str]:
        req = urllib.request.Request(url, headers=self.default_headers, method="GET")
        with urllib.request.urlopen(req, timeout=self.timeout) as resp:
            content_type = resp.headers.get("Content-Type", "")
            charset = resp.headers.get_content_charset() or "utf-8"
            body = resp.read().decode(charset, errors="replace")
            return int(resp.status), content_type, body


class DmojDownloader:
    def __init__(self, client: HttpClient, username: str, max_pages: int, pause_s: float, verbose: bool) -> None:
        self.client = client
        self.username = username
        self.max_pages = max_pages
        self.pause_s = pause_s
        self.verbose = verbose

    def log(self, msg: str) -> None:
        if self.verbose:
            print(msg, file=sys.stderr)

    def fetch_recent_submission_ids(self) -> List[int]:
        seen: set[int] = set()
        ordered_ids: List[int] = []

        for page in range(1, self.max_pages + 1):
            page_ids = self._fetch_ids_for_page(page)
            if not page_ids:
                self.log(f"page {page}: no submissions found")
                break
            for sid in page_ids:
                if sid not in seen:
                    seen.add(sid)
                    ordered_ids.append(sid)
            time.sleep(self.pause_s)

        return ordered_ids

    def _fetch_ids_for_page(self, page: int) -> List[int]:
        for template in LISTING_ENDPOINT_CANDIDATES:
            path = template.format(username=self.username)
            query = {"page": page}
            if template.endswith("/submissions"):
                query["user"] = self.username

            url = self.client.make_url(path, query)
            try:
                status, content_type, body = self.client.get(url)
            except urllib.error.HTTPError as e:
                self.log(f"{url} -> HTTP {e.code}")
                continue
            except urllib.error.URLError as e:
                self.log(f"{url} -> URL error: {e}")
                continue

            if status != 200:
                continue

            ids: List[int]
            if "json" in content_type.lower() or body.strip().startswith(("{", "[")):
                ids = self._extract_ids_from_json_body(body)
            else:
                ids = self._extract_ids_from_html(body)

            if ids:
                self.log(f"page {page}: found {len(ids)} ids via {path}")
                return ids

        return []

    def _extract_ids_from_json_body(self, body: str) -> List[int]:
        try:
            payload = json.loads(body)
        except json.JSONDecodeError:
            return []

        ids: List[int] = []
        for obj in iter_dicts(payload):
            sid = coerce_int(obj.get("id"))
            if sid is None:
                continue

            userish = (
                obj.get("user")
                or obj.get("username")
                or obj.get("author")
                or obj.get("creator")
            )
            if isinstance(userish, dict):
                userish = userish.get("username") or userish.get("id")

            if userish is None or str(userish) == self.username:
                ids.append(sid)

        return dedupe_keep_order(ids)

    def _extract_ids_from_html(self, body: str) -> List[int]:
        ids = [int(m.group(1)) for m in re.finditer(r"href=\"/submission/(\d+)\"", body)]
        return dedupe_keep_order(ids)

    def fetch_submission(self, sid: int) -> Optional[Submission]:
        meta = self._fetch_submission_meta(sid)
        problem = first_nonempty_str(meta, ["problem", "problem_code", "problem_id"]) or "unknown"
        language = first_nonempty_str(meta, ["language", "lang", "language_name"]) or "unknown"
        verdict = first_nonempty_str(meta, ["result", "status", "verdict"]) or "unknown"

        source = first_nonempty_str(meta, ["source", "source_code", "code"])
        if not source:
            source = self._fetch_source_text(sid)

        if not self._is_accepted(meta, verdict):
            return None

        if not source:
            self.log(f"submission {sid}: accepted but source unavailable")
            return None

        return Submission(
            submission_id=sid,
            problem=problem,
            language=language,
            verdict=verdict,
            source=source,
        )

    def _fetch_submission_meta(self, sid: int) -> Dict[str, Any]:
        for template in SUBMISSION_DETAIL_ENDPOINT_CANDIDATES:
            path = template.format(id=sid)
            url = self.client.make_url(path)

            try:
                status, content_type, body = self.client.get(url)
            except urllib.error.HTTPError as e:
                self.log(f"{url} -> HTTP {e.code}")
                continue
            except urllib.error.URLError as e:
                self.log(f"{url} -> URL error: {e}")
                continue

            if status != 200:
                continue

            if "json" in content_type.lower() or body.strip().startswith(("{", "[")):
                try:
                    payload = json.loads(body)
                except json.JSONDecodeError:
                    continue
                out = extract_submission_like_dict(payload, sid)
                if out:
                    return out
            else:
                # HTML fallback: detect obvious verdict/problem hints.
                out = {
                    "id": sid,
                    "verdict": extract_html_verdict(body),
                    "problem": extract_html_problem_code(body),
                }
                return out

        return {"id": sid}

    def _fetch_source_text(self, sid: int) -> Optional[str]:
        for template in SOURCE_ENDPOINT_CANDIDATES:
            path = template.format(id=sid)
            url = self.client.make_url(path)
            try:
                status, content_type, body = self.client.get(url)
            except urllib.error.HTTPError as e:
                self.log(f"{url} -> HTTP {e.code}")
                continue
            except urllib.error.URLError as e:
                self.log(f"{url} -> URL error: {e}")
                continue

            if status != 200:
                continue

            if "json" in content_type.lower() or body.strip().startswith(("{", "[")):
                try:
                    payload = json.loads(body)
                except json.JSONDecodeError:
                    continue
                source = first_nonempty_str(payload, ["source", "source_code", "code"])
                if not source:
                    for obj in iter_dicts(payload):
                        source = first_nonempty_str(obj, ["source", "source_code", "code"])
                        if source:
                            break
                if source:
                    return source
                continue

            source = extract_source_from_html_or_text(body)
            if source:
                return source

        return None

    def _is_accepted(self, meta: Dict[str, Any], verdict: str) -> bool:
        if isinstance(meta.get("is_accepted"), bool):
            return bool(meta["is_accepted"])

        score = coerce_float(meta.get("score") or meta.get("points"))
        max_score = coerce_float(meta.get("max_score") or meta.get("total") or meta.get("max_points"))
        if score is not None and max_score is not None and max_score > 0:
            if score >= max_score:
                return True

        v = (verdict or "").strip().lower()
        return any(marker in v for marker in AC_MARKERS)


def iter_dicts(node: Any) -> Iterable[Dict[str, Any]]:
    if isinstance(node, dict):
        yield node
        for v in node.values():
            yield from iter_dicts(v)
    elif isinstance(node, list):
        for item in node:
            yield from iter_dicts(item)


def extract_submission_like_dict(payload: Any, sid: int) -> Dict[str, Any]:
    for obj in iter_dicts(payload):
        obj_sid = coerce_int(obj.get("id"))
        if obj_sid == sid:
            return obj
    for obj in iter_dicts(payload):
        if any(k in obj for k in ("source", "source_code", "code", "status", "result", "verdict")):
            return obj
    return {}


def first_nonempty_str(node: Any, keys: List[str]) -> Optional[str]:
    if not isinstance(node, dict):
        return None
    for k in keys:
        v = node.get(k)
        if isinstance(v, str) and v.strip():
            return v.strip()
        if isinstance(v, dict):
            nested = first_nonempty_str(v, keys)
            if nested:
                return nested
    return None


def coerce_int(value: Any) -> Optional[int]:
    try:
        if value is None:
            return None
        return int(value)
    except (TypeError, ValueError):
        return None


def coerce_float(value: Any) -> Optional[float]:
    try:
        if value is None:
            return None
        return float(value)
    except (TypeError, ValueError):
        return None


def dedupe_keep_order(items: List[int]) -> List[int]:
    seen: set[int] = set()
    out: List[int] = []
    for x in items:
        if x in seen:
            continue
        seen.add(x)
        out.append(x)
    return out


def safe_filename_component(text: str) -> str:
    cleaned = re.sub(r"[^A-Za-z0-9._-]+", "_", text.strip())
    cleaned = cleaned.strip("._")
    return cleaned or "unknown"


def guess_extension(language: str) -> str:
    lang = language.lower()
    mapping = {
        "python": ".py",
        "py": ".py",
        "pypy": ".py",
        "c++": ".cpp",
        "cpp": ".cpp",
        "c": ".c",
        "java": ".java",
        "kotlin": ".kt",
        "go": ".go",
        "rust": ".rs",
        "javascript": ".js",
        "typescript": ".ts",
        "pascal": ".pas",
    }
    for key, ext in mapping.items():
        if key in lang:
            return ext
    return ".txt"


def extract_html_verdict(body: str) -> str:
    m = re.search(r"(Accepted|AC|Wrong Answer|WA|Time Limit Exceeded|TLE|Runtime Error|RTE)", body, re.IGNORECASE)
    return m.group(1) if m else "unknown"


def extract_html_problem_code(body: str) -> str:
    patterns = [
        r"/problem/([A-Za-z0-9_-]+)",
        r"problem-code\"\s*>\s*([A-Za-z0-9_-]+)",
    ]
    for pat in patterns:
        m = re.search(pat, body)
        if m:
            return m.group(1)
    return "unknown"


def extract_source_from_html_or_text(body: str) -> Optional[str]:
    block_patterns = [
        r"<pre[^>]*id=\"source\"[^>]*>(.*?)</pre>",
        r"<pre[^>]*id=\"submission-source\"[^>]*>(.*?)</pre>",
        r"<code[^>]*>(.*?)</code>",
    ]
    for pat in block_patterns:
        m = re.search(pat, body, re.IGNORECASE | re.DOTALL)
        if m:
            return html.unescape(strip_tags(m.group(1))).strip("\n")

    # If this is already plain-text code, keep it.
    if "<html" not in body.lower() and "</" not in body:
        return body.rstrip("\n")
    return None


def strip_tags(text: str) -> str:
    return re.sub(r"<[^>]+>", "", text)


def write_submission(out_dir: str, sub: Submission) -> str:
    problem = safe_filename_component(sub.problem)
    lang = safe_filename_component(sub.language)
    ext = guess_extension(sub.language)
    filename = f"{sub.submission_id}_{problem}_{lang}{ext}"
    path = os.path.join(out_dir, filename)
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(sub.source or "")
        if not (sub.source or "").endswith("\n"):
            f.write("\n")
    return path


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Download recent accepted DMOJ submissions for a user.")
    p.add_argument("username", help="DMOJ username")
    p.add_argument("--max-pages", type=int, default=3, help="How many recent submission pages to scan (default: 3)")
    p.add_argument("--base-url", default="https://dmoj.ca", help="Judge base URL (default: https://dmoj.ca)")
    p.add_argument("--output-dir", default="dmoj_submissions", help="Directory to store downloaded source files")
    p.add_argument("--api-token", default=os.environ.get("DMOJ_API_TOKEN"), help="DMOJ API token (or env DMOJ_API_TOKEN)")
    p.add_argument("--sessionid", default=os.environ.get("DMOJ_SESSIONID"), help="DMOJ sessionid cookie (or env DMOJ_SESSIONID)")
    p.add_argument("--timeout", type=float, default=20.0, help="HTTP timeout in seconds")
    p.add_argument("--sleep", type=float, default=0.25, help="Delay between requests in seconds")
    p.add_argument("--verbose", action="store_true", help="Print debug logs to stderr")
    return p.parse_args()


def main() -> int:
    args = parse_args()

    if args.max_pages <= 0:
        print("--max-pages must be > 0", file=sys.stderr)
        return 2

    os.makedirs(args.output_dir, exist_ok=True)

    client = HttpClient(
        base_url=args.base_url,
        api_token=args.api_token,
        sessionid=args.sessionid,
        timeout=args.timeout,
        user_agent="dmoj-recent-ac-dump/1.0",
    )
    downloader = DmojDownloader(
        client=client,
        username=args.username,
        max_pages=args.max_pages,
        pause_s=args.sleep,
        verbose=args.verbose,
    )

    ids = downloader.fetch_recent_submission_ids()
    if not ids:
        print("No submission IDs found. Try --verbose and verify auth/token, username, and endpoint compatibility.")
        return 1

    saved = 0
    for sid in ids:
        sub = downloader.fetch_submission(sid)
        if not sub:
            time.sleep(args.sleep)
            continue
        path = write_submission(args.output_dir, sub)
        print(f"saved {sub.submission_id} ({sub.problem}, {sub.language}, {sub.verdict}) -> {path}")
        saved += 1
        time.sleep(args.sleep)

    print(f"done: saved {saved} accepted submissions out of {len(ids)} discovered")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
