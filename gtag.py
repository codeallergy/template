#!/bin/python

import sys
import os
from fnmatch import fnmatch

template = """
    <!-- Google tag (gtag.js) -->
    <script async src="https://www.googletagmanager.com/gtag/js?id=MYGTAG"></script>
    <script>
      window.dataLayer = window.dataLayer || [];
      function gtag(){dataLayer.push(arguments);}
      gtag('js', new Date());

      gtag('config', 'MYGTAG');
    </script>
"""
endingHead = "</head>"
gtag = ""

def list_files(dir, pattern, cb):
    for path, subdirs, files in os.walk(dir):
        for name in files:
            if fnmatch(name, pattern):
                cb(os.path.join(path, name))

def read_fully(file_path):
    with open(file_path, 'r') as file:
        return file.read()

def save_fully(file_path, content):
    with open(file_path, 'w') as file:
        file.write(content)

def patch_file(file_path):
    file_content = read_fully(file_path)

    if "googletagmanager" in file_content:
        print(f"Error: File {file_path} already has gtag")
        return

    idx = file_content.find(endingHead)
    if idx == -1:
        print(f"Error: File {file_path} does not have closing </head> tag")
        return

    new_content = file_content[:idx] + gtag + file_content[idx:]
    save_fully(file_path, new_content)

    print(f"Patched html file : {file_path}")

def main():
    global gtag
    if len(sys.argv) < 3:
        print("Usage: python3 gtag.py gtag folder")
        return
    if "MYGTAG" == sys.argv[1]:
        print("Warning: replace MYGTAG in main.go file")
        return
    gtag = template.replace("MYGTAG", sys.argv[1])
    folder = sys.argv[2]
    list_files(folder, "*.html", patch_file)

if __name__ == "__main__":
    main()