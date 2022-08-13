#!/usr/bin/python3
import os
import sys
from json import load
from base64 import b64decode
from urllib.parse import unquote
from gzip import decompress

def main():
  ign_filename = sys.argv[1]
  with open(ign_filename, 'r') as fh: ign_json = load(fh)

  ign_directories = ign_json.get('storage', {}).get('directories', {})
  ign_files = ign_json.get('storage', {}).get('files', {})
  ign_units = ign_json.get('systemd', {}).get('units', {})

  for dir in ign_directories:
    path = dir['path']
    mode = int(dir['mode'])
    if not os.path.exists(path): os.makedirs(path, mode)
    os.chmod(path, mode)

  for file in ign_files:
    path = file['path']
    mode = int(file['mode'])
    dict_contents = file.get('contents', {})
    content_compression = dict_contents.get('compression', '')
    content_source = dict_contents.get('source', '')

    if content_compression == "gzip":
      base64encoded_string = content_source.replace('data:;base64,', '')
      gzipped_string = b64decode(base64encoded_string)
      decompressed_bytestring = decompress(gzipped_string)
      the_content = str(decompressed_bytestring)[2:-1]

      # \n !!!!
      with open(path, 'w') as fh: fh.write(the_content)
    else:
      quoted_string = content_source.replace('data:,', '')
      the_content = unquote(quoted_string)
      with open(path, 'w') as fh: fh.write(the_content)

if __name__ == "__main__":
  main()
