import os
import yaml
import requests
import argparse

import magic
from moviepy.editor import VideoFileClip

def is_video(filename):
  file_type = magic.from_file(filename, mime=True)
  return file_type.startswith('video/')

def main(args):
  for filename in os.listdir(args.yaml_dir):
    if filename.endswith('.yaml') or filename.endswith('.yml'):
      file_path = os.path.join(args.yaml_dir, filename)

      with open(file_path, 'r') as file:
        parsed_yaml = yaml.safe_load(file)

      thumbnail_url = parsed_yaml['thumbnail']
      tmp_filename = os.path.normpath(thumbnail_url).split(os.sep)[-1]
      response = requests.get(thumbnail_url)
      response.raise_for_status()

      with open(tmp_filename, 'wb') as f:
        f.write(response.content)
      print(f"File downloaded and saved as {tmp_filename}")
      
      if is_video(tmp_filename):
        print(f"video type found: {thumbnail_url}")
        arxiv_id = os.path.normpath(parsed_yaml['link']).split(os.sep)[-1]
        gif_filename = f"{args.gif_output_dir}/{arxiv_id}.gif"
        
        print(f"converting {tmp_filename} => {gif_filename}")
        videoClip = VideoFileClip(tmp_filename)
        videoClip = videoClip.resize(args.resize_factor)
        videoClip = videoClip.subclip(0, args.clip_to)
        videoClip.speedx(args.speedx).to_gif(gif_filename)

        parsed_yaml['thumbnail'] = f"https://github.com/{args.gif_output_github_repo}/blob/main/{gif_filename}?raw=true"
        with open(file_path, 'w') as file:
            yaml.dump(parsed_yaml, file, default_flow_style=False)

      print(f"remove {tmp_filename}")
      os.remove(tmp_filename)
        
if __name__ == "__main__":
  parser = argparse.ArgumentParser()
  parser.add_argument('--yaml-dir', type=str, default="current")
  parser.add_argument('--gif-output-github-repo', type=str, default="codingpot/hf-daily-paper-newsletter-tester")
  parser.add_argument('--gif-output-dir', type=str, default="assets")
  parser.add_argument('--speedx', type=int, default=6)
  parser.add_argument('--clip-to', type=int, default=15)
  parser.add_argument('--resize-factor', type=float, default=0.4)

  args = parser.parse_args()
  main(args)
