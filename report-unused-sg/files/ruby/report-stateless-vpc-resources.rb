#!/usr/bin/env ruby

require "pry-byebug"
require "json"
require "bundler/setup"
require "aws-sdk-ec2"
require "aws-sdk-s3"

def nat_gateway_ids_for_vpc(client, vpc_id)
  filter = [ { name: "vpc-id", values: [vpc_id] } ]
  data = client.describe_nat_gateways(filter: filter)
  data.nat_gateways.map { |ng| ng.nat_gateway_id }.sort
end

def subnets_ids_for_vpc(client, vpc_id)
  filters = [ { name: "vpc-id", values: [vpc_id] } ]
  data = client.describe_subnets(filters: filters)
  data.subnets.map { |sn| sn.subnet_id }.sort
end

def internet_gateway_ids_for_vpc(client)
    data = client.describe_internet_gateways()
    data.internet_gateways.map { |ng| ng.internet_gateway_id }.sort
end

def vpc_ids(client)
    data = client.describe_vpcs()
    data.vpcs.map { |ng| ng.vpc_id }.sort
end


def nat_gateway_ids_from_terraform_state(statefile)
  str = File.read("terraform.tfstate")
  data = JSON.parse(str)
  list = data["resources"]
  nat_gateway = list.filter { |m| m["name"] == "private_nat_gateway" }.first
  nat_gateway["instances"].map { |ng| ng.dig("attributes", "nat_gateway_id") }.sort
end


def internet_gateway_ids_from_terraform_state(statefile)
    str = File.read("terraform.tfstate")
    data = JSON.parse(str)
    list = data["resources"]
    nat_gateway = list.filter { |m| m["name"] == "public_internet_gateway" }.first
    nat_gateway["instances"].map { |ng| ng.dig("attributes", "gateway_id") }.sort
end


def route_tables_for_subnet(client, subnet_id)
  filters = [ { name: "association.subnet-id", values: [subnet_id] } ]
  data = client.describe_route_tables(filters: filters)
  data[:route_tables][0][:associations].each do |route|
  #pp route[:route_table_association_id]
    pp route[:route_table_id]
  end
end


def get_state_vpc_names(s3)

  s3_state_bucket_keys = s3.bucket("cloud-platform-terraform-state").objects(prefix:'cloud-platform-network', delimiter: '').collect(&:"key")

  keys = []
  vpc_names = []
  s3_state_bucket_keys.each { |vpc_name| keys << vpc_name.delete(' ') }
  
  
  keys.each do |key|
    key_split = key.split('/')
    vpc_names.push(key_split[1])
  end

  return vpc_names

end 


def get_vpc_ids_from_aws(client)

  vpc_actual_ids = []

  # Get the whole data for all vpcs 
  data = client.describe_vpcs()

  # Create a file containing the vpc ids only
  File.open("output-files/vpc_actual_ids.txt","w") do |f|
    f.puts data.vpcs.map { |vpc| vpc.vpc_id }.sort
  end

  File.open('output-files/vpc_actual_ids.txt').each { |vpc| vpc_actual_ids << vpc }

  #vpc_actual_ids.shift(1)
  return vpc_actual_ids

end 

# The s3 key for the network state contains the vpc name. We therefore need to dynamically fetch each name in the state
# and iterate through them to get the corresponding vpc id. The state file can then be downloaded for each vpc and finally getting vpc id
def get_vpc_ids_from_state(s3)

  vpc_ids_in_state = []
  
  #Iterate the vpc names in the state and get the corresponding vpc ids
  get_state_vpc_names(s3).each do |vpc_name|
    begin
      bucket_name = "cloud-platform-terraform-state"
      key = "cloud-platform-network/"+vpc_name+"/terraform.tfstate"
      statefile_name_output = "state-files/vpc-network-"+vpc_name+".tfstate"
      
      download_state_from_s3(s3, bucket_name, key, statefile_name_output)
    
      str = File.read(statefile_name_output)
      data = JSON.parse(str)
      vpc_id = data['outputs']['vpc_id']['value']
      vpc_ids_in_state.push(vpc_id)
    rescue => e
      #puts vpc_name+' has no vpc id in its state file'
    end
  end

  return vpc_ids_in_state

end


def download_state_from_s3(s3, bucket_name, key, statefile_path)
  # Loop through all the dynamically fetched vpc names and download the network state file from s3
  obj = s3.bucket(bucket_name).object(key)
  obj.get(response_target: statefile_path)
end

s3 = Aws::S3::Resource.new(region:'eu-west-1', profile: ENV["AWS_PROFILE"])

#binding.pry
ec2 = Aws::EC2::Client.new(region:'eu-west-2', profile: ENV["AWS_PROFILE"])



#pp get_vpc_ids_from_state(statefile_name_output)


#Get the actual vpc ids
#vpc_ids_aws = get_vpc_ids_from_aws(ec2)

vpc_ids_from_aws = get_vpc_ids_from_aws(ec2)
vpc_ids_from_state = get_vpc_ids_from_state(s3)

puts 'vpc ids from AWS'
vpc_ids_from_aws.each do |vpc_id|
  puts vpc_id
end

#Get the state vpc ids
puts 'vpc ids from state files'

vpc_ids_from_state.each do |vpc_id|
  puts vpc_id
end





vpc_ids_with_state = []

vpc_ids_from_state.each do |vpc_id_state|
  vpc_ids_from_aws.each do |vpc_id_aws|
      if vpc_id_state.delete(' ') == vpc_id_aws.delete(' ')
          vpc_ids_with_state.push(vpc_id_state)
      end
     # break vpc_id_aws if vpc_id_state == vpc_id_aws
  end
end

puts 'STATELESS'

vpc_ids_with_state.each do |vpc_id_with_state|
  puts vpc_id_with_state
end




#******** Get the natgateway ids **********************

#pp nat_gateway_ids_for_vpc(ec2, "vpc-0c16457fd570a1f0b")
#pp nat_gateway_ids_from_terraform_state("terraform.tfstate")

#********Get the internet gateway*************************************************

#pp internet_gateway_ids_for_vpc(ec2)
#pp internet_gateway_ids_from_terraform_state("terraform.tfstate")

#********Get the vpc ids***********************
#pp vpc_ids(ec2)

#*********Get the subnets*******************

#pp subnets_ids_for_vpc(ec2, "vpc-0c16457fd570a1f0b")


#****Get the route table ids for each subnet*******

#pp route_tables_for_subnet(ec2, "subnet-0386053fb92e69994")
#route_tables_for_subnet(ec2, "subnet-05fe20e709ed69201")


