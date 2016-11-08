require './helpers'
require 'yaml'

def get_kubeconfig
  kubeconfig = "#{ENV['HOME']}/.kube/config"
  destdir = Dir.pwd
  FileUtils.cp(kubeconfig, destdir)
end

def deploy_kubernetes conf
  system("minikube start --host-only-cidr=\"#{conf['k8s_ip_base']}1/24\"")
  get_kubeconfig
end

def stop_kubernetes
  system("minikube stop")
end

def destroy_kubernetes
  system("minikube delete")
end

def reload_kubernetes conf
  destroy_kubernetes
  deploy_kubernetes conf
end

def handle_kubernetes_action conf
  unless in_path?("minikube")
    puts "please install kubectl before deploying kubernetes"
    exit(1)
  end
  unless in_path?("minikube")
    puts "please install minikube before deploying kubernetes"
    exit(1)
  end

  case ARGV.first
    when "up"
      deploy_kubernetes conf
    when "down"
      stop_kubernetes
    when "reload"
      reload_kubernetes conf
    when "destroy"
      destroy_kubernetes
  end
end
