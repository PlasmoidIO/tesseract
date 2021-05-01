import os
import shutil
import time
import subprocess


def create_workspace(target_dir: str):
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    os.mkdir(target_dir)
    if not target_dir.startswith("/"):
        target_dir = "/" + target_dir
    p = subprocess.Popen(
        ["go", "build", "-o", "../tests" + target_dir], cwd="../peer")
    p.wait()


def remove_workspace(target_dir: str):
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    shutil.rmtree(target_dir)


create_workspace("hi")
time.sleep(1)
remove_workspace("hi")
