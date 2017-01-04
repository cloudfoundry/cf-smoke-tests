require 'sinatra'
require 'open3'
require 'json'

STDOUT.sync = true
STDERR.sync = true

get '/put/' do
  host = params[:host]
  path = params[:path]
  port = params[:port]
  data = params[:data]

  stdout, stderr, status = Open3.capture3("curl -XPUT -d #{data} #{host}:#{port}#{path}")

  { stdout: stdout, stderr: stderr, return_code: status.exitstatus }.to_json
end
