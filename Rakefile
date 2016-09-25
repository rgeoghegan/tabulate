DEV_DEPS = ["github.com/stretchr/testify"]

def go(args)
	ENV['GOPATH'] = Dir.pwd
	sh "go #{args}"
end

desc "Build all the executables"
task :build => :deps

desc "Start doc server on port 8080"
task :doc do
	ENV['GOPATH'] = Dir.pwd
	sh "godoc -http=:8080"
end

task :default => :test

task :test => :devdeps do
    go "test ."
end

task :testv => :devdeps do
    go "test -v ."
end

desc "Try code in play.go"
task :play do
	go "run play.go"
end

desc "Install any third-party dev dependencies"
task :devdeps do
	DEV_DEPS.each do |n|
		go("get #{n}")
	end
end
