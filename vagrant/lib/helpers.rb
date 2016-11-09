
def is_plugin(name)
  if Vagrant.has_plugin?(name)
    puts "using #{name}"
  else
    puts "please run vagrant plugin install #{name}"
    exit(1)
  end
end

def in_path?(cmd)
  exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(';') : ['']
  ENV['PATH'].split(File::PATH_SEPARATOR).each do |path|
    exts.each { |ext|
      exe = File.join(path, "#{cmd}#{ext}")
      return exe if File.executable?(exe) && !File.directory?(exe)
    }
  end
  nil
end