require 'json'

file = File.read('actual_network_state.json')


data_hash = JSON.parse(file)

pp data_hash.keys

