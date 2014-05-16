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

puts "#{name}.dice\t#{rand(6) + 1}\t#{time}"
