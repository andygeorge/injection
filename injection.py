#!/usr/bin/python3
from os import system,chmod,makedirs
from os.path import exists
from sys import exit
from json import load
from base64 import b64decode
from urllib.parse import unquote
from gzip import decompress
from argparse import ArgumentParser

SYSTEMD_UNIT_PATH = '/etc/systemd/system'

def main():
  parser = ArgumentParser(description='Process an Ignition file on an already-running system.')
  parser.add_argument('path', help='Path to a valid Ignition file')
  args = parser.parse_args()
  ign_filename = args.path

  if not exists(ign_filename):
    print(f'{ign_filename} does not exist, please enter a valid Ignition file.')
    exit()
  with open(ign_filename, 'r') as fh: ign_json = load(fh)

  ign_directories = ign_json.get('storage', {}).get('directories', {})
  ign_files = ign_json.get('storage', {}).get('files', {})
  ign_units = ign_json.get('systemd', {}).get('units', {})

  for dir in ign_directories:
    path = dir['path']
    mode = int(dir['mode'])
    if not exists(path): makedirs(path, mode)
    chmod(path, mode)

  for file in ign_files:
    path = file['path']
    mode = int(file['mode'])
    dict_contents = file.get('contents', {})
    content_compression = dict_contents.get('compression', '')
    content_source = dict_contents.get('source', '')

    if content_compression == 'gzip':
      base64encoded_string = content_source.replace('data:;base64,', '')
      gzipped_string = b64decode(base64encoded_string)
      decompressed_bytestring = decompress(gzipped_string)
      raw_string = str(decompressed_bytestring)[2:-1]
      the_content = raw_string.replace('\\n', '\n').replace('\\\\', '\\')

      with open(path, 'w') as fh: fh.write(the_content)
    else:
      quoted_string = content_source.replace('data:,', '')
      the_content = unquote(quoted_string)
      with open(path, 'w') as fh: fh.write(the_content)

  for unit in ign_units:
    name = unit['name']
    path = f'{SYSTEMD_UNIT_PATH}/{name}'
    enabled = bool(unit['enabled'])
    content = unit['contents']

    with open(path, 'w') as fh: fh.write(content)
    if enabled:
      system(f'systemctl enable {name}')
      system('systemctl daemon-reload')
      system(f'systemctl start {name}')

if __name__ == "__main__":
  main()
