#!/usr/bin/env ruby

require 'optparse'

name = nil
time = Time.now.to_i

OptionParser.new do |opts|
  opts.on('-n', '--name NAME', 'Metric name') do |n|
    name = n
  end
end.parse!

name ||= "example"

if ENV['MACKEREL_AGENT_PLUGIN_META'] == '1'
  puts <<-META
# mackerel-agent-plugin
prefix = "foo.bar"
[graphs.#{name}]
label = "My Dice #{name}"
unit = "integer"
[graphs.#{name}.metrics.dice]
label = "The Die"
  META
  exit 0
end

puts "#{name}.dice\t#{rand(6) + 1}\t#{time}"
