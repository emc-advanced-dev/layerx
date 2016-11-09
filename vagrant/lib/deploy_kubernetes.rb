require_relative './helpers'
require 'yaml'


def get_kubeconfig
  kubeconfig = "#{ENV['HOME']}/.kube/config"
  destdir = Dir.pwd
  FileUtils.cp(kubeconfig, destdir)
end

def deploy_kubernetes conf
  puts "deploying kubernetes"
  system("minikube start --host-only-cidr=\"#{conf['k8s_ip_base']}1/24\"")
  get_kubeconfig
end

def stop_kubernetes
  puts "stopping kubernetes"
  system("minikube stop")
end

def destroy_kubernetes
  unless ARGV.include? "-f"
    puts "Are you sure you want to destroy the 'minikube' VM? [y/N]"
    ans = gets
    unless ans == 'y'
      return
    end
  end
  puts "destroy kubernetes"
  system("minikube delete")
end

def reload_kubernetes conf
  puts "reloading kubernetes"
  destroy_kubernetes
  deploy_kubernetes conf
end

def status_kubernetes
  system("minikube status")
end

def handle_kubernetes_action conf
  unless in_path?("minikube")
    puts "please install minikube before deploying kubernetes (make sure it is in PATH)"
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
    when "status"
      status_kubernetes
  end
end
