import os

#====================Var declaration=======================#
pod_name = "pod-asgn"
#==========================================================#

def prune_docker_images():
    os.system("docker login")

def create_docker_images():
    os.system("docker build -t psdudeman39/pod-image:latest .")
    os.system("docker push psdudeman39/pod-image:latest")

# apply role and rolebindings
os.system("kubectl apply -f clusterrole.yaml")
os.system("kubectl apply -f clusterrolebinding.yaml")

# delete the pod-asgn
pod_delete_command = "kubectl delete pod " + pod_name
os.system(pod_delete_command)

prune_docker_images()
create_docker_images()

# create pod asgn
pod_create_command = "kubectl apply -f " + pod_name + ".yaml"
os.system(pod_create_command)