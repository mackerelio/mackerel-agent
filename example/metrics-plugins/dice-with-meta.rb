if ENV['MACKEREL_AGENT_PLUGIN_META'] == '1'
  require 'json'

  meta = {
    :graphs => {
      :dice => {
        :label   => 'My Dice',
        :unit    => 'integer',
        :metrics => [
          {
            :name  => 'd6',
            :label => 'Die (d6)'
          }, {
            :name  => 'd20',
            :label => 'Die (d20)'
          }
        ]
      }
    }
  }

  puts '# mackerel-agent-plugin'
  puts meta.to_json
  exit 0
end

puts [ 'dice.d6',  rand( 6)+1, Time.now.to_i ].join("\t")
puts [ 'dice.d20', rand(20)+1, Time.now.to_i ].join("\t")
