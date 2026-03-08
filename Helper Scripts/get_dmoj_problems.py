import argparse
import requests
from bs4 import BeautifulSoup
import pathlib
import re
import time
import json
from http.cookiejar import MozillaCookieJar

def session_from_cookies(cookie_file: str):
    jar = MozillaCookieJar(cookie_file)
    jar.load(ignore_discard=True, ignore_expires=True)

    s= requests.Session()
    s.headers.update({"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"})
    s.cookies.update(jar)
    return s

def fetch_submission_ids(username: str, max_pages: int, problem_status: list[str], sess):
    url = f"https://dmoj.ca/submissions/user/{username}"
    items = []

    for page in range(1, max_pages + 1):
        page_url = f"{url}/{page}"
        r = sess.get(page_url, timeout=10)
        print("status code:", r.status_code)
        r.raise_for_status()

        html = r.text

        soup = BeautifulSoup(html, "html.parser")


        for row in soup.select("#submissions-table .submission-row"):
            result = row.select_one(".sub-result .status")
            if not result:
                continue
            verdict = result.get_text(strip=True)
            if verdict not in problem_status:
                continue
            
            submission_id = row.get("id")

            prob_a = row.select_one(".sub-main .sub-info .name a")
            problem_href = prob_a["href"] if prob_a and prob_a.has_attr("href") else None

            items.append({
                "submission_id": submission_id,
                "verdict": verdict,
                "problem_code": problem_href[9:] if problem_href else None,
            })

    return items

def get_raw_files(submission_id: str, sess):
    url = f"https://dmoj.ca/src/{submission_id}/raw"
    r = sess.get(url, timeout=10)
    r.raise_for_status()

    return r.text


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("username", help="DMOJ username")
    parser.add_argument("--max_pages", help="Maximum number of pages to scrape", type=int, default=3)
    parser.add_argument("--cookie_file", help="Path to cookie file", default="cookies.txt")
    args = parser.parse_args()

    print("username:", args.username)
    print("max_pages:", args.max_pages)

    sess = session_from_cookies(args.cookie_file)
    ids = fetch_submission_ids(args.username, args.max_pages, ["AC"], sess)
    # print(ids)
    problemids = set()
    deduped_items = []
    for item in ids:
        #dedupe by problem
        if item["problem_code"] in problemids:
            continue
        problemids.add(item["problem_code"])
        deduped_items.append(item)

    for item in deduped_items:
        # Move files with problem codes not in deduped_items
        problem_codes = {item["problem_code"] for item in deduped_items}
        pathlib.Path("./problems/notfound").mkdir(parents=True, exist_ok=True)
        
        for file in pathlib.Path("./problems").glob("*.cpp"):
            problem_code = file.stem
            if problem_code not in problem_codes:
                file.rename(f"./problems/notfound/{file.name}")
                print(f"Moved {file.name} to ./problems/notfound")

    # for item in deduped_items:
    #     print(f"Fetching raw code for submission {item['submission_id']} with verdict {item['verdict']}...")
    #     try:
    #         raw_code = get_raw_files(item["submission_id"], sess)
    #         item["raw_code"] = raw_code
    #         pathlib.Path("./problems").mkdir(parents=True, exist_ok=True)
    #         with open(f"./problems/{item['problem_code']}.cpp", "w", newline="") as f:
    #             f.write(raw_code.rstrip())
    #         print(f"Saved to ./problems/{item['problem_code']}.cpp")
    #     except Exception as e:
    #         print(f"Error fetching raw code for submission {item['submission_id']}: {e}")
    #         item["raw_code"] = None


        

if __name__ == "__main__":
    main()
