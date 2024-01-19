import os
import sys
import argparse
import concurrent.futures

from dstack.api import Client, Task, GPU, Client, Resources

def parallel_job(i, client, arxiv_id, args):
  port = 6006 + i
  
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
          f'mkdir -p papers/{arxiv_id}',
  
          f'echo "Downloading paper (ID: {arxiv_id})"...',
          f'paper {arxiv_id} -p -d papers/{arxiv_id}/',
          f'PDF_FN=$(ls -1 papers/{arxiv_id}/*.pdf | head -n 1)',

          'echo "OCR processing with nougat-ocr..."',
          f'nougat $PDF_FN -o papers/{arxiv_id}',
          'MMD_FN=$(echo $PDF_FN | sed "s/pdf/mmd/g")',
  
          'echo "Translation processing"...',
          'python arxiv-translator/translate_mmd.py $MMD_FN',
          f'python arxiv-translator/ready_templates.py $MMD_FN arxiv-translator/assets/html_template.html papers/{arxiv_id}',
  
          'echo "Git commit & push"...',
          f'cd {os.path.normpath(args.target_archive_github_repo).split(os.sep)[1]}',
          f'git pull',
          f'mkdir -p {args.target_archive_dir}/{arxiv_id}',
          f'cp ../papers/{arxiv_id}/paper.ko.html ./{args.target_archive_dir}/{arxiv_id}',
          f'git config --global user.name "{args.github_realname}"',
          f'git config --global user.email "{args.github_email}"',
          'git add . ',
          f'git commit -m "automatically added {arxiv_id}-ko paper"',
          'git push origin main',
      ],
      ports=[str(port)],
  )
  
  run = client.runs.submit(
      run_name=f'{args.dstack_run_name}({arxiv_id})',
      configuration=task,
      resources=Resources(
          gpu=GPU(memory="16GB")
      ),
      spot_policy="on-demand"
  )
  
  run.attach()
  for log in run.logs():
    print(log)
    sys.stdout.flush()
  run.detach()

  return f'{args.dstack_run_name}({arxiv_id}) is completed'

def main(args):
  client = Client.from_config(
      project_name=args.dstack_project,
      server_url="https://cloud.dstack.ai",
      user_token=args.dstack_token
  )
  
  num_jobs = len(args.arxiv_ids)
  
  clients = [client for _ in range(num_jobs)]
  arxiv_ids = [arxiv_id for arxiv_id in args.arxiv_ids]
  args_list = [args for _ in range(num_jobs)]
  
  with concurrent.futures.ThreadPoolExecutor(max_workers=num_jobs) as executor:
    futures = [
      executor.submit(parallel_job, i, clients[i], arxiv_ids[i], args_list[i]) for i in range(num_jobs-1)
    ]
    
    concurrent.futures.wait(futures)

  last_idx = num_jobs-1
  parallel_job(last_idx, clients[last_idx], arxiv_ids[last_idx], args_list[last_idx])

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
