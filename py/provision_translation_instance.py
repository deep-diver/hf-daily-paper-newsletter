import os
import sys
import argparse

from dstack.api import Client, Task, GPU, Client, Resources

def main(args):
  client = Client.from_config(
      project_name=args.dstack_project,
      server_url="https://cloud.dstack.ai",
      user_token=args.dstack_token
  )

  arxiv_ids_str = ' '.join(args.arxiv_ids)
  
  task = Task(
      python=args.dstack_python_version,
      commands=[
          'echo "Cloneing the arXiv translator repo..."',
          f'git clone https://{args.github_username}:{args.github_token}@github.com/deep-diver/arxiv-translator.git',
          f'git clone https://{args.github_username}:{args.github_token}@github.com/{args.target_archive_github_repo}.git',
  
          'echo "Installing gh CLI..."',
          'type -p curl >/dev/null || (apt update && apt install curl -y)',
          'curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null && apt update && apt install gh -y',
  
          'echo "Login to the GitHub account..."',
          f'echo {args.github_token} > GH_TOKEN.txt',
          'gh auth login --with-token < GH_TOKEN.txt',
  
          'echo "Installing requirements..."',
          'pip install nougat-ocr arxiv-dl kss',
  
          'echo "Creating temporary directory for holding the paper..."',
          f'for arxiv_id in {arxiv_ids_str}; do mkdir -p papers/$arxiv_id; done',
  
          f'echo "Downloading paper (IDs: {args.arxiv_ids})"...',
          f'for arxiv_id in {arxiv_ids_str}; do paper $arxiv_id -p -d papers/$arxiv_id/; done',

          'echo "OCR processing with nougat-ocr..."',
          f'for arxiv_id in {arxiv_ids_str}; do PDF_FN=$(ls -1 papers/$arxiv_id/*.pdf | head -n 1); nougat $PDF_FN -o papers/$arxiv_id; done',
  
          'echo "Translation processing"...',
          f'for arxiv_id in {arxiv_ids_str}; do PDF_FN=$(ls -1 papers/$arxiv_id/*.pdf | head -n 1); MMD_FN=$(echo $PDF_FN | sed "s/pdf/mmd/g"); python arxiv-translator/translate_mmd.py $MMD_FN; python arxiv-translator/ready_templates.py $MMD_FN arxiv-translator/assets/html_template.html papers/$arxiv_id; done',
  
          'echo "Git commit & push"...',
          f'cd {os.path.normpath(args.target_archive_github_repo).split(os.sep)[1]}',
          f'git pull',
          f'for arxiv_id in {arxiv_ids_str}; do mkdir -p {args.target_archive_dir}/$arxiv_id; cp ../papers/$arxiv_id/paper.ko.html ./{args.target_archive_dir}/$arxiv_id; done',
          f'git config --global user.name "{args.github_realname}"',
          f'git config --global user.email "{args.github_email}"',
          'git add . ',
          f'git commit -m "automatically added [{arxiv_ids_str}]-ko paper"',
          'git push origin main',
      ],
      ports=["6006"],
  )
  
  run = client.runs.submit(
      run_name=args.dstack_run_name,
      configuration=task,
      resources=Resources(
          gpu=GPU(memory="16GB")
      ),
      spot_policy="on-demand"
  )
  
  run.attach()
  for log in run.logs():
    print(log)
  run.detach()

if __name__ == "__main__":
  parser = argparse.ArgumentParser()
  parser.add_argument('--dstack-token', type=str, default="")
  parser.add_argument('--dstack-project', type=str, default="deep-diver-main")
  parser.add_argument('--dstack-run-name', type=str, default="paper-ko-translation-run")
  parser.add_argument('--dstack-python-version', type=str, default="3.10")
  
  parser.add_argument('--github-token', type=str, default="")
  parser.add_argument('--github-username', type=str, default="deep-diver")
  parser.add_argument('--github-realname', type=str, default="Chansung Park")
  parser.add_argument('--github-email', type=str, default="deep.diver.csp@gmail.com")

  parser.add_argument('--arxiv-ids', nargs='+', help='List of paper ids')
  parser.add_argument('--target-archive-github-repo', type=str, default="deep-diver/hf-daily-paper-newsletter")
  parser.add_argument('--target-archive-dir', type=str, default="translated-papers")
  args = parser.parse_args()
  print(args)
  main(args)
