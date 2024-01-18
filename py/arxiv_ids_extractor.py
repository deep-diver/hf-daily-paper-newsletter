import os
import yaml
import argparse

def main(args):
  arxiv_ids = ''

  for filename in os.listdir(args.yaml_dir):
      if filename.endswith('.yaml') or filename.endswith('.yml'):
          file_path = os.path.join(args.yaml_dir, filename)

          with open(file_path, 'r') as file:
              parsed_yaml = yaml.safe_load(file)
              arxiv_ids = arxiv_ids + '"' + os.path.normpath(parsed_yaml['link']).split(os.sep)[-1] + '"' + " "

  arxiv_ids = arxiv_ids.strip()

  if not arxiv_ids:
    print("no arxiv ids to write")
  else:
    with open(args.output_filename, 'w') as file:
        file.write(arxiv_ids)

if __name__ == "__main__":
  parser = argparse.ArgumentParser()
  parser.add_argument('--yaml-dir', type=str, default="current")
  parser.add_argument('--output-filename', type=str, default="arxiv_ids.txt")

  args = parser.parse_args()
  main(args)
