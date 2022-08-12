#!/usr/bin/python3
import os
import sys
import json

def main():
  ign_filename = sys.argv[1]
  with open(ign_filename, 'r') as fh:
    ign_json = json.load(fh)

  ign_directories = ign_json['storage']['directories']
  ign_files = ign_json['storage']['files']
  ign_units = ign_json['systemd']['units']

  for dir in ign_directories:
    path = dir['path']
    mode = int(dir['mode'])
    if not os.path.exists(path):
      os.makedirs(path, mode)
    else:
      os.chmod(path, mode)

if __name__ == "__main__":
  main()
