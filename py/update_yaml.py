import os
import yaml
import argparse

def main(args):
  entries = os.listdir(args.translated_output_dir)
  arxiv_ids_in_dir = [
      entry for entry in entries 
      if os.path.isdir(os.path.join(args.translated_output_dir, entry))
  ]

  for filename in os.listdir(args.yaml_dir):
      if filename.endswith('.yaml') or filename.endswith('.yml'):
          file_path = os.path.join(args.yaml_dir, filename)

          with open(file_path, 'r') as file:
              parsed_yaml = yaml.safe_load(file)

          arxiv_id_in_yaml = os.path.normpath(parsed_yaml['link']).split(os.sep)[-1]
          if arxiv_id_in_yaml in arxiv_ids_in_dir:
              parsed_yaml.update(
                  {
                      'translated_paths': {
                          "INT": f"https://raw.githack.com/{args.translated_output_github_repo}/main/{args.translated_output_dir}/{arxiv_id_in_yaml}/{args.non_translated_filename}",
                          "KR": f"https://raw.githack.com/{args.translated_output_github_repo}/main/{args.translated_output_dir}/{arxiv_id_in_yaml}/{args.translated_filename}"
                      }
                  }
              )

              with open(file_path, 'w') as file:
                  yaml.dump(parsed_yaml, file, default_flow_style=False)

if __name__ == "__main__":
  parser = argparse.ArgumentParser()
  parser.add_argument('--yaml-dir', type=str, default="current")
  parser.add_argument('--translated-output-github-repo', type=str, default="codingpot/hf-daily-paper-newsletter-tester")
  parser.add_argument('--translated-output-dir', type=str, default="translated-papers")
  parser.add_argument('--non-translated-filename', type=str, default="paper.en.html")
  parser.add_argument('--translated-filename', type=str, default="paper.ko.html")

  args = parser.parse_args()
  main(args)
