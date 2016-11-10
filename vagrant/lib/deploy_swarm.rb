require_relative './helpers'
require 'yaml'
require 'json'
require 'fileutils'

def get_swarm_config conf
  docker_env = `docker-machine env #{conf['docker_machine_name']}`.gsub!("\"", "").gsub!("export ", "").split("\n")
  docker_tls_verify = nil
  docker_host = nil
  docker_cert_path = nil
  docker_env.each do |var|
    if var.include?("DOCKER_TLS_VERIFY")
      docker_tls_verify = var.split("=")[1]
    end
    if var.include?("DOCKER_HOST")
      docker_host = var.split("=")[1]
    end
    if var.include?("DOCKER_CERT_PATH")
      docker_cert_path = var.split("=")[1]
    end
  end

  destdir = "#{Dir.pwd}/swarm/"
  FileUtils.mkdir_p destdir
  ['ca.pem', 'cert.pem', 'key.pem', 'server-key.pem', 'server.pem'].each do |file|
    FileUtils.cp("#{docker_cert_path}/#{file}", "#{destdir}/#{file}")
  end

  docker_machine_ip = JSON.load(`docker-machine inspect #{conf['docker_machine_name']}`)["Driver"]["IPAddress"]

  docker_config = {
      docker_tls_verify: docker_tls_verify,
      docker_host: docker_host,
      docker_machine_name: conf['docker_machine_name'],
      docker_machine_ip: docker_machine_ip
  }
  File.open("#{destdir}/docker_config.yml", 'w') {|f| f.write docker_config.to_yaml }
  docker_config
end

def deploy_docker_swarm conf
  puts "deploying docker_swarm"
  system("docker-machine create -d virtualbox --virtualbox-hostonly-cidr #{conf['docker_swarm_ip_base']}1/24 #{conf['docker_machine_name']}")
  docker_config = get_swarm_config conf
  ENV['DOCKER_TLS_VERIFY'] = docker_config[:docker_tls_verify]
  ENV['DOCKER_HOST'] = docker_config[:docker_host]
  ENV['DOCKER_MACHINE_NAME'] = docker_config[:docker_machine_name]
  ENV['DOCKER_CERT_PATH'] = "#{Dir.pwd}/swarm/"
  system("docker swarm init --advertise-addr #{docker_config[:docker_machine_ip]}")
end

def stop_docker_swarm conf
  puts "stopping docker_swarm"
  system("docker-machine stop #{conf['docker_machine_name']}")
end

def destroy_docker_swarm conf
  puts "destroy docker_swarm"
  system("docker-machine rm --force #{conf['docker_machine_name']}")
end

def reload_docker_swarm conf
  puts "reloading docker_swarm"
  system("docker-machine restart #{conf['docker_machine_name']}")
end

def status_docker_swarm conf
  system("docker-machine status #{conf['docker_machine_name']}")
end

def handle_docker_swarm_action conf
  unless in_path?("docker")
    puts "please install docker before deploying docker_swarm (make sure it is in PATH)"
    exit(1)
  end
  unless in_path?("docker-machine")
    puts "please install docker-machine before deploying docker_swarm (make sure it is in PATH)"
    exit(1)
  end

  case ARGV.first
    when "up"
      deploy_docker_swarm conf
    when "down"
      stop_docker_swarm conf
    when "reload"
      reload_docker_swarm conf
    when "destroy"
      destroy_docker_swarm conf
    when "status"
      status_docker_swarm conf
  end
end
