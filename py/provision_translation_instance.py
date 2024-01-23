import os
import sys
import time
import argparse
import concurrent.futures

from dstack.api import Client, Task, GPU, Client, Resources

def parallel_job_split(i, client, arxiv_ids, args):
  port = 6006 + i

  arxiv_ids_str = ' '.join(arxiv_ids)
  
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
  
          f'echo "Downloading paper (IDs: {arxiv_ids_str})"...',
          f'for arxiv_id in {arxiv_ids_str}; do paper $arxiv_id -p -d papers/$arxiv_id/; done',

          'echo "OCR processing with nougat-ocr..."',
          f'for arxiv_id in {arxiv_ids_str}; do PDF_FN=$(ls -1 papers/$arxiv_id/*.pdf | head -n 1); nougat $PDF_FN -o papers/$arxiv_id; done',
  
          'echo "Translation processing"...',
          f'for arxiv_id in {arxiv_ids_str}; do PDF_FN=$(ls -1 papers/$arxiv_id/*.pdf | head -n 1); MMD_FN=$(echo $PDF_FN | sed "s/pdf/mmd/g"); python arxiv-translator/translate_mmd.py $MMD_FN; python arxiv-translator/ready_templates.py $MMD_FN arxiv-translator/assets/html_template.html papers/$arxiv_id; done',
  
          'echo "Git commit & push"...',
          f'cd {os.path.normpath(args.target_archive_github_repo).split(os.sep)[1]}',
          f'git pull',
          f'for arxiv_id in {arxiv_ids_str}; do mkdir -p {args.target_archive_dir}/$arxiv_id; cp ../papers/$arxiv_id/paper.ko.html ../papers/$arxiv_id/paper.en.html ./{args.target_archive_dir}/$arxiv_id; done',
          f'git config --global user.name "{args.github_realname}"',
          f'git config --global user.email "{args.github_email}"',
          'git add . ',
          f'git commit -m "automatically added [{arxiv_ids_str}]-ko paper"',
          'git push origin main',
      ],
      ports=[str(port)],
  )
  
  run = client.runs.submit(
      run_name=f"{args.dstack_run_name}_{i}",
      configuration=task,
      resources=Resources(
          gpu=GPU(memory="16GB")
      ),
      spot_policy="on-demand"
  )

  return run

def split_jobs(total_jobs, num_workers):
  # Ensure that the number of workers is not zero to avoid division by zero error
  if num_workers <= 0:
    return "Number of workers must be greater than 0"

  jobs_per_worker = total_jobs // num_workers
  remaining_jobs = total_jobs % num_workers

  result = []

  for i in range(num_workers):
    if i < remaining_jobs:
      # Assign an extra job to the first 'remaining_jobs' workers
      result.append(jobs_per_worker + 1)
    else:
      result.append(jobs_per_worker)

  return result

def distribute_parameters(parameters, job_distribution):
  distributed_params = []
  start_index = 0

  for jobs in job_distribution:
    # The end index for slicing is calculated by adding the number of jobs to the start index
    end_index = start_index + jobs
    # Slice the parameters list from start_index to end_index
    worker_params = parameters[start_index:end_index]
    # Append the sliced parameters to the distributed_params list
    distributed_params.append(worker_params)
    # Update the start_index for the next iteration
    start_index = end_index

  return distributed_params

def main(args):
  client = Client.from_config(
      project_name=args.dstack_project,
      server_url="https://cloud.dstack.ai",
      user_token=args.dstack_token
  )
  
  total_num_jobs = len(args.arxiv_ids)

  runs = {}
  done_status = ["done", "terminated", "failed"]
  work_splits = split_jobs(total_num_jobs, args.num_of_workers)
  arxiv_id_splits = distribute_parameters([arxiv_id for arxiv_id in args.arxiv_ids], work_splits)

  for i, arxiv_ids in enumerate(arxiv_id_splits):
    if arxiv_ids:
      run = parallel_job_split(i, client, arxiv_ids, args)
      run.arxiv_ids = arxiv_ids
      runs[i] = run

  while True:
    to_be_removed = []
    
    for index, (idx, run) in enumerate(runs.items()):
      if run.status not in done_status:
        run.refresh()

        if run.status in done_status:
          to_be_removed.append(idx)
          
    for to_be_removed_worker_idx in to_be_removed:
      del runs[to_be_removed_worker_idx]

    if runs:
      print(runs)
    if not runs:
      break
    
    time.sleep(10)

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

  parser.add_argument('--num-of-workers', type=int, default=4)
  args = parser.parse_args()
  print(args)
  main(args)
