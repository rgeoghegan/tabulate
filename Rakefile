LIBS = ["tabulate"]
LIB_FILES = FileList[LIBS.map{|n| "src/#{n}/*.go"}]

TOOL_FILES = FileList["src/tools/*.go"]
TOOL_NAME_TO_FILE = {}
TOOL_FILES.each { |e| TOOL_NAME_TO_FILE[e.pathmap("%n")] = e}

DEPS = []
DEV_DEPS = ["github.com/stretchr/testify"]

def go(args)
	ENV['GOPATH'] = Dir.pwd
	sh "go #{args}"
end

desc "Build all the executables"
task :build => :deps

desc "Run specific tool"
task :run, [:toolname] do |t, args|
	tool = TOOL_NAME_TO_FILE[args[:toolname]]
	go("run #{tool}")
end

desc "Start doc server on port 8080"
task :doc do
	ENV['GOPATH'] = Dir.pwd
	sh "godoc -http=:8080"
end

task :default => :build

task :test => :devdeps do
    LIBS.each do |n|
		go "test #{n}"
	end
end

task :test_v => :devdeps do
    LIBS.each do |n|
		go "test #{n}"
	end
end

desc "Try code in play.go"
task :play do
	go "run src/play.go"
end

desc "Clean up any built binaries"
task :clean

TOOL_FILES.each do |n|
    exec = n.pathmap("%n")
    file exec => n do
        go "build #{n}"
	end
    file exec => LIB_FILES

	task :build => exec
	task :clean do
		rm n
	end
end

desc "Install any third-party dependencies"
task :deps do
	DEPS.each do |n|
		go("get #{n}")
	end
end

desc "Install any third-party dev dependencies"
task :devdeps do
	DEV_DEPS.each do |n|
		go("get #{n}")
	end
end
