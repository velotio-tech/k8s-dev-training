import os

os.system("kubectl delete deploy --all --force")
os.system("kubectl delete pods --all --force")

os.system("make install")
os.system("make run")